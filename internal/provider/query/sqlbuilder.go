package query

import (
	"fmt"
	"strings"

	"github.com/inovacc/dataprovider/internal/provider"
)

const (
	joinTemplate        = "JOIN %s ON %s"
	leftJoinTemplate    = "LEFT JOIN %s ON %s"
	rightJoinTemplate   = "RIGHT JOIN %s ON %s"
	whereTemplate       = " WHERE %s"
	groupByTemplate     = "GROUP BY %s"
	orderByTemplate     = "ORDER BY %s"
	limitTemplate       = "LIMIT %d"
	offsetTemplate      = "OFFSET %d"
	havingTemplate      = "HAVING %s"
	selectTemplate      = "SELECT %s FROM %s"
	createTableTemplate = "CREATE TABLE %s (%s)"
	dropTableTemplate   = "DROP TABLE %s"
	deleteTemplate      = "DELETE FROM %s"
	insertTemplate      = "INSERT INTO %s (%s) VALUES (%s)"
	updateTemplate      = "UPDATE %s SET %s"
)

type PlaceholderFormatter interface {
	ReplacePlaceholders(query string) string
}

type postgresFormatter struct{}

func (f *postgresFormatter) ReplacePlaceholders(query string) string {
	for i := 1; strings.Contains(query, "?"); i++ {
		query = strings.Replace(query, "?", fmt.Sprintf("$%d", i), 1)
	}
	return query
}

type oracleFormatter struct{}

func (f *oracleFormatter) ReplacePlaceholders(query string) string {
	for i := 1; strings.Contains(query, "?"); i++ {
		query = strings.Replace(query, "?", fmt.Sprintf(":p%d", i), 1)
	}
	return query
}

type defaultFormatter struct{}

func (f *defaultFormatter) ReplacePlaceholders(query string) string {
	return query
}

func NewFormatter(driver string) PlaceholderFormatter {
	switch driver {
	case provider.PostgresSQLDatabaseProviderName:
		return &postgresFormatter{}
	case provider.OracleDatabaseProviderName:
		return &oracleFormatter{}
	default:
		return &defaultFormatter{}
	}
}

type SQLBuilder interface {
	Select(table string, columns ...string) SQLBuilder
	Where(condition string, args ...any) SQLBuilder
	Join(table, onCondition string) SQLBuilder
	LeftJoin(table, onCondition string) SQLBuilder
	RightJoin(table, onCondition string) SQLBuilder
	GroupBy(columns ...string) SQLBuilder
	Having(condition string, args ...any) SQLBuilder
	OrderBy(columns ...string) SQLBuilder
	Limit(n int) SQLBuilder
	Offset(n int) SQLBuilder
	CreateTable(table string, definition string) SQLBuilder
	DropTable(table string) SQLBuilder
	DeleteFrom(table string) SQLBuilder
	Build() (string, []any)
	Clear() SQLBuilder
	InsertInto(table string, columns ...string) SQLBuilder
	Values(args ...any) SQLBuilder
	Update(table string) SQLBuilder
	Set(column string, value any) SQLBuilder
}

type queryBuilder struct {
	opts       provider.Options
	table      string
	columns    []string
	joins      []string
	where      []string
	groupBy    []string
	having     []string
	orderBy    []string
	insertCols []string
	insertVals []string
	updateSet  []string
	args       []any
	limit      *int
	offset     *int
	special    string
	formatter  PlaceholderFormatter
}

func NewQueryBuilder(opts provider.Options) SQLBuilder {
	return &queryBuilder{
		opts:      opts,
		formatter: NewFormatter(opts.Driver),
	}
}

func (b *queryBuilder) InsertInto(table string, columns ...string) SQLBuilder {
	b.table = table
	b.insertCols = columns
	return b
}

func (b *queryBuilder) Values(args ...any) SQLBuilder {
	placeholders := make([]string, len(args))
	for i := range args {
		placeholders[i] = "?"
	}
	b.insertVals = placeholders
	b.args = append(b.args, args...)
	return b
}

func (b *queryBuilder) Update(table string) SQLBuilder {
	b.table = table
	return b
}

func (b *queryBuilder) Set(column string, value any) SQLBuilder {
	b.updateSet = append(b.updateSet, fmt.Sprintf("%s = ?", column))
	b.args = append(b.args, value)
	return b
}

func (b *queryBuilder) CreateTable(table string, definition string) SQLBuilder {
	b.special = fmt.Sprintf(createTableTemplate, table, definition)
	return b
}

func (b *queryBuilder) DropTable(table string) SQLBuilder {
	b.special = fmt.Sprintf(dropTableTemplate, table)
	return b
}

func (b *queryBuilder) DeleteFrom(table string) SQLBuilder {
	b.special = fmt.Sprintf(deleteTemplate, table)
	return b
}

func (b *queryBuilder) Select(table string, columns ...string) SQLBuilder {
	b.table = table
	b.columns = columns
	return b
}

func (b *queryBuilder) Where(condition string, args ...any) SQLBuilder {
	b.where = append(b.where, condition)
	b.args = append(b.args, args...)
	return b
}

func (b *queryBuilder) Join(table, onCondition string) SQLBuilder {
	b.joins = append(b.joins, fmt.Sprintf(joinTemplate, table, onCondition))
	return b
}

func (b *queryBuilder) LeftJoin(table, onCondition string) SQLBuilder {
	b.joins = append(b.joins, fmt.Sprintf(leftJoinTemplate, table, onCondition))
	return b
}

func (b *queryBuilder) RightJoin(table, onCondition string) SQLBuilder {
	b.joins = append(b.joins, fmt.Sprintf(rightJoinTemplate, table, onCondition))
	return b
}

func (b *queryBuilder) GroupBy(columns ...string) SQLBuilder {
	b.groupBy = append(b.groupBy, columns...)
	return b
}

func (b *queryBuilder) Having(condition string, args ...any) SQLBuilder {
	b.having = append(b.having, condition)
	b.args = append(b.args, args...)
	return b
}

func (b *queryBuilder) OrderBy(columns ...string) SQLBuilder {
	b.orderBy = append(b.orderBy, columns...)
	return b
}

func (b *queryBuilder) Limit(n int) SQLBuilder {
	b.limit = &n
	return b
}

func (b *queryBuilder) Offset(n int) SQLBuilder {
	b.offset = &n
	return b
}

func (b *queryBuilder) Clear() SQLBuilder {
	b.table = ""
	b.columns = nil
	b.joins = nil
	b.where = nil
	b.groupBy = nil
	b.having = nil
	b.orderBy = nil
	b.args = nil
	b.limit = nil
	b.offset = nil
	b.insertCols = nil
	b.insertVals = nil
	b.updateSet = nil
	b.special = ""
	return b
}

func (b *queryBuilder) Build() (string, []any) {
	if b.special != "" {
		return b.special, b.args
	}

	if len(b.insertCols) > 0 && len(b.insertVals) > 0 {
		query := fmt.Sprintf(insertTemplate,
			b.table,
			strings.Join(b.insertCols, ", "),
			strings.Join(b.insertVals, ", "))
		return b.formatter.ReplacePlaceholders(query), b.args
	}

	if len(b.updateSet) > 0 {
		query := fmt.Sprintf(updateTemplate, b.table, strings.Join(b.updateSet, ", "))
		if len(b.where) > 0 {
			query += fmt.Sprintf(whereTemplate, strings.Join(b.where, " AND "))
		}
		return b.formatter.ReplacePlaceholders(query), b.args
	}

	var sb strings.Builder

	columns := "*"
	if len(b.columns) > 0 {
		columns = strings.Join(b.columns, ", ")
	}

	sb.WriteString(fmt.Sprintf(selectTemplate, columns, b.table))

	if len(b.joins) > 0 {
		sb.WriteString(" ")
		sb.WriteString(strings.Join(b.joins, " "))
	}

	if len(b.where) > 0 {
		sb.WriteString(fmt.Sprintf(whereTemplate, strings.Join(b.where, " AND ")))
	}

	if len(b.groupBy) > 0 {
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf(groupByTemplate, strings.Join(b.groupBy, ", ")))
	}

	if len(b.having) > 0 {
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf(havingTemplate, strings.Join(b.having, " AND ")))
	}

	if len(b.orderBy) > 0 {
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf(orderByTemplate, strings.Join(b.orderBy, ", ")))
	}

	if b.limit != nil {
		sb.WriteString(fmt.Sprintf(" %s", fmt.Sprintf(limitTemplate, *b.limit)))
	}

	if b.offset != nil {
		sb.WriteString(fmt.Sprintf(" %s", fmt.Sprintf(offsetTemplate, *b.offset)))
	}

	return b.formatter.ReplacePlaceholders(sb.String()), b.args
}
