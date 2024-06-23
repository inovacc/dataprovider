package queries

import (
	"fmt"
	"github.com/spf13/afero"
	"path/filepath"
	"strings"
)

// this package is a query builder for SQL queries like SELECT, UPDATE, INSERT, DELETE, PL/SQL, etc.

type Queries interface {
	// Select generates a SELECT query
	Select(columns ...string) Queries

	// From generates a FROM clause
	From(table string) Queries

	// Where generates a WHERE clause
	Where(condition string) Queries

	// OrderBy generates an ORDER BY clause
	OrderBy(columns ...string) Queries

	// Limit generates a LIMIT clause
	Limit(limit int) Queries

	// Offset generates an OFFSET clause
	Offset(offset int) Queries

	// GroupBy generates a GROUP BY clause
	GroupBy(columns ...string) Queries

	// Having generates a HAVING clause
	Having(condition string) Queries

	// Join generates a JOIN clause
	Join(table string, condition string) Queries

	// LeftJoin generates a LEFT JOIN clause
	LeftJoin(table string, condition string) Queries

	// RightJoin generates a RIGHT JOIN clause
	RightJoin(table string, condition string) Queries

	// FullJoin generates a FULL JOIN clause
	FullJoin(table string, condition string) Queries

	// CrossJoin generates a CROSS JOIN clause
	CrossJoin(table string) Queries

	// Union generates a UNION clause
	Union(query Queries) Queries

	// UnionAll generates a UNION ALL clause
	UnionAll(query Queries) Queries

	// SubQuery generates a subquery
	SubQuery(query Queries, alias string) Queries

	// Build generates the final SQL query
	Build() string
}

type queries struct {
	selectCols   []string
	fromTable    string
	whereClause  string
	orderByCols  []string
	limitValue   int
	offsetValue  int
	groupByCols  []string
	havingClause string
	joins        []string
	unions       []string
	subQueries   []string
}

func NewQueries() Queries {
	return &queries{}
}

func (q *queries) Select(columns ...string) Queries {
	q.selectCols = columns
	return q
}

func (q *queries) From(table string) Queries {
	q.fromTable = table
	return q
}

func (q *queries) Where(condition string) Queries {
	q.whereClause = condition
	return q
}

func (q *queries) OrderBy(columns ...string) Queries {
	q.orderByCols = columns
	return q
}

func (q *queries) Limit(limit int) Queries {
	q.limitValue = limit
	return q
}

func (q *queries) Offset(offset int) Queries {
	q.offsetValue = offset
	return q
}

func (q *queries) GroupBy(columns ...string) Queries {
	q.groupByCols = columns
	return q
}

func (q *queries) Having(condition string) Queries {
	q.havingClause = condition
	return q
}

func (q *queries) Join(table string, condition string) Queries {
	joinClause := fmt.Sprintf("JOIN %s ON %s", table, condition)
	q.joins = append(q.joins, joinClause)
	return q
}

func (q *queries) LeftJoin(table string, condition string) Queries {
	joinClause := fmt.Sprintf("LEFT JOIN %s ON %s", table, condition)
	q.joins = append(q.joins, joinClause)
	return q
}

func (q *queries) RightJoin(table string, condition string) Queries {
	joinClause := fmt.Sprintf("RIGHT JOIN %s ON %s", table, condition)
	q.joins = append(q.joins, joinClause)
	return q
}

func (q *queries) FullJoin(table string, condition string) Queries {
	joinClause := fmt.Sprintf("FULL JOIN %s ON %s", table, condition)
	q.joins = append(q.joins, joinClause)
	return q
}

func (q *queries) CrossJoin(table string) Queries {
	joinClause := fmt.Sprintf("CROSS JOIN %s", table)
	q.joins = append(q.joins, joinClause)
	return q
}

func (q *queries) Union(query Queries) Queries {
	unionClause := fmt.Sprintf("UNION (%s)", query.Build())
	q.unions = append(q.unions, unionClause)
	return q
}

func (q *queries) UnionAll(query Queries) Queries {
	unionClause := fmt.Sprintf("UNION ALL (%s)", query.Build())
	q.unions = append(q.unions, unionClause)
	return q
}

func (q *queries) SubQuery(query Queries, alias string) Queries {
	subQueryClause := fmt.Sprintf("(%s) AS %s", query.Build(), alias)
	q.subQueries = append(q.subQueries, subQueryClause)
	return q
}

func (q *queries) Build() string {
	var query []string

	if len(q.selectCols) > 0 {
		query = append(query, fmt.Sprintf("SELECT %s", strings.Join(q.selectCols, ", ")))
	}

	if q.fromTable != "" {
		query = append(query, fmt.Sprintf("FROM %s", q.fromTable))
	}

	if len(q.joins) > 0 {
		query = append(query, strings.Join(q.joins, " "))
	}

	if q.whereClause != "" {
		query = append(query, fmt.Sprintf("WHERE %s", q.whereClause))
	}

	if len(q.groupByCols) > 0 {
		query = append(query, fmt.Sprintf("GROUP BY %s", strings.Join(q.groupByCols, ", ")))
	}

	if q.havingClause != "" {
		query = append(query, fmt.Sprintf("HAVING %s", q.havingClause))
	}

	if len(q.orderByCols) > 0 {
		query = append(query, fmt.Sprintf("ORDER BY %s", strings.Join(q.orderByCols, ", ")))
	}

	if q.limitValue > 0 {
		query = append(query, fmt.Sprintf("LIMIT %d", q.limitValue))
	}

	if q.offsetValue > 0 {
		query = append(query, fmt.Sprintf("OFFSET %d", q.offsetValue))
	}

	if len(q.unions) > 0 {
		query = append(query, strings.Join(q.unions, " "))
	}

	return strings.Join(query, " ")
}

func GetQueryFromFile(filename string) (string, error) {
	fs := afero.NewOsFs()

	ok, err := afero.DirExists(fs, filepath.Dir(filename))
	if err != nil {
		return "", err
	}

	if !ok {
		return "", fmt.Errorf("directory %s does not exist", filepath.Dir(filename))
	}

	ok, err = afero.Exists(fs, filename)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", fmt.Errorf("file %s does not exist", filename)
	}

	content, err := afero.ReadFile(fs, filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
