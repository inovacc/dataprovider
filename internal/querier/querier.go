package querier

import (
	"fmt"
	"strings"
)

// this package is a query builder for SQL queries like SELECT, UPDATE, INSERT, DELETE, PL/SQL, etc.

type Querier interface {
	// Select generates a SELECT query
	Select(columns ...string) Querier

	// From generates a FROM clause
	From(table string) Querier

	// Where generates a WHERE clause
	Where(condition string) Querier

	// OrderBy generates an ORDER BY clause
	OrderBy(columns ...string) Querier

	// Limit generates a LIMIT clause
	Limit(limit int) Querier

	// Offset generates an OFFSET clause
	Offset(offset int) Querier

	// GroupBy generates a GROUP BY clause
	GroupBy(columns ...string) Querier

	// Having generates a HAVING clause
	Having(condition string) Querier

	// Join generates a JOIN clause
	Join(table string, condition string) Querier

	// LeftJoin generates a LEFT JOIN clause
	LeftJoin(table string, condition string) Querier

	// RightJoin generates a RIGHT JOIN clause
	RightJoin(table string, condition string) Querier

	// FullJoin generates a FULL JOIN clause
	FullJoin(table string, condition string) Querier

	// CrossJoin generates a CROSS JOIN clause
	CrossJoin(table string) Querier

	// Union generates a UNION clause
	Union(query Querier) Querier

	// UnionAll generates a UNION ALL clause
	UnionAll(query Querier) Querier

	// SubQuery generates a subquery
	SubQuery(query Querier, alias string) Querier

	// Build generates the final SQL query
	Build() string
}

type querier struct {
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

func NewQuerier() Querier {
	return &querier{}
}

func (q *querier) Select(columns ...string) Querier {
	q.selectCols = columns
	return q
}

func (q *querier) From(table string) Querier {
	q.fromTable = table
	return q
}

func (q *querier) Where(condition string) Querier {
	q.whereClause = condition
	return q
}

func (q *querier) OrderBy(columns ...string) Querier {
	q.orderByCols = columns
	return q
}

func (q *querier) Limit(limit int) Querier {
	q.limitValue = limit
	return q
}

func (q *querier) Offset(offset int) Querier {
	q.offsetValue = offset
	return q
}

func (q *querier) GroupBy(columns ...string) Querier {
	q.groupByCols = columns
	return q
}

func (q *querier) Having(condition string) Querier {
	q.havingClause = condition
	return q
}

func (q *querier) Join(table string, condition string) Querier {
	joinClause := fmt.Sprintf("JOIN %s ON %s", table, condition)
	q.joins = append(q.joins, joinClause)
	return q
}

func (q *querier) LeftJoin(table string, condition string) Querier {
	joinClause := fmt.Sprintf("LEFT JOIN %s ON %s", table, condition)
	q.joins = append(q.joins, joinClause)
	return q
}

func (q *querier) RightJoin(table string, condition string) Querier {
	joinClause := fmt.Sprintf("RIGHT JOIN %s ON %s", table, condition)
	q.joins = append(q.joins, joinClause)
	return q
}

func (q *querier) FullJoin(table string, condition string) Querier {
	joinClause := fmt.Sprintf("FULL JOIN %s ON %s", table, condition)
	q.joins = append(q.joins, joinClause)
	return q
}

func (q *querier) CrossJoin(table string) Querier {
	joinClause := fmt.Sprintf("CROSS JOIN %s", table)
	q.joins = append(q.joins, joinClause)
	return q
}

func (q *querier) Union(query Querier) Querier {
	unionClause := fmt.Sprintf("UNION (%s)", query.Build())
	q.unions = append(q.unions, unionClause)
	return q
}

func (q *querier) UnionAll(query Querier) Querier {
	unionClause := fmt.Sprintf("UNION ALL (%s)", query.Build())
	q.unions = append(q.unions, unionClause)
	return q
}

func (q *querier) SubQuery(query Querier, alias string) Querier {
	subQueryClause := fmt.Sprintf("(%s) AS %s", query.Build(), alias)
	q.subQueries = append(q.subQueries, subQueryClause)
	return q
}

func (q *querier) Build() string {
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
