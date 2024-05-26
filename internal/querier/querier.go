package querier

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
	Build() (string, []any)
}

type querier struct {
}

func (q querier) Select(columns ...string) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) From(table string) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) Where(condition string) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) OrderBy(columns ...string) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) Limit(limit int) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) Offset(offset int) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) GroupBy(columns ...string) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) Having(condition string) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) Join(table string, condition string) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) LeftJoin(table string, condition string) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) RightJoin(table string, condition string) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) FullJoin(table string, condition string) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) CrossJoin(table string) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) Union(query Querier) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) UnionAll(query Querier) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) SubQuery(query Querier, alias string) Querier {
	//TODO implement me
	panic("implement me")
}

func (q querier) Build() (string, []any) {
	//TODO implement me
	panic("implement me")
}

// NewQuerier creates a new Querier
func NewQuerier() Querier {
	return &querier{}
}
