// Package query SQLBuilder with adaptations based on SQL and DML/DDL best practices from educational sources and database dialect specifics
package query

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
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
	updateSetTemplate   = "UPDATE %s SET %s"
)

type stringKinds string

const (
	stringKindSelect stringKinds = "select"
	stringKindInsert stringKinds = "insert"
	stringKindUpdate stringKinds = "update"
	stringKindDelete stringKinds = "delete"
	stringKindCreate stringKinds = "create"
	stringKindDrop   stringKinds = "drop"
)

// SQLBuilder interface models typical SQL DDL and DML operations for various dialects
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
	InsertInto(table string, columns ...string) SQLBuilder
	Values(args ...any) SQLBuilder
	Update(table string) SQLBuilder
	Set(column string, value any) SQLBuilder
	Build() (string, []any)
	Clear() SQLBuilder
	As(alias string) SQLBuilder
	Raw(clause string, args ...any) SQLBuilder
	MergeInto(table string) SQLBuilder
	On(condition string) SQLBuilder
	WhenMatched(updateSet map[string]any) SQLBuilder
	WhenNotMatchedInsert(columns []string, values []any) SQLBuilder
	Union(other SQLBuilder) SQLBuilder
	ExportAsJSON() (string, error)
	ExportAsXML() (string, error)
	StructToSQL(data any, table string, isInsert bool) (string, []any, error)
}

type queryBuilder struct {
	opts            provider.Options
	kind            stringKinds
	table           string
	columns         []string
	joins           []string
	where           []string
	groupBy         []string
	having          []string
	orderBy         []string
	insertCols      []string
	insertVals      []string
	updateSet       []string
	alias           string
	rawClauses      []string
	mergeTable      string
	mergeOn         string
	mergeMatchedSet []string
	mergeInsertCols []string
	mergeInsertVals []string
	args            []any
	limit           *int
	offset          *int
	special         string
	formatter       PlaceholderFormatter
}

func NewQueryBuilder(opts provider.Options) SQLBuilder {
	return &queryBuilder{
		opts:      opts,
		formatter: NewFormatter(opts.Driver),
	}
}

func (b *queryBuilder) As(alias string) SQLBuilder {
	b.alias = alias
	return b
}

func (b *queryBuilder) Raw(clause string, args ...any) SQLBuilder {
	b.rawClauses = append(b.rawClauses, clause)
	b.args = append(b.args, args...)
	return b
}

func (b *queryBuilder) MergeInto(table string) SQLBuilder {
	b.mergeTable = table
	return b
}

func (b *queryBuilder) On(condition string) SQLBuilder {
	b.mergeOn = condition
	return b
}

func (b *queryBuilder) WhenMatched(updateSet map[string]any) SQLBuilder {
	for col, val := range updateSet {
		b.mergeMatchedSet = append(b.mergeMatchedSet, fmt.Sprintf("%s = ?", col))
		b.args = append(b.args, val)
	}
	return b
}

func (b *queryBuilder) WhenNotMatchedInsert(columns []string, values []any) SQLBuilder {
	b.mergeInsertCols = columns
	b.mergeInsertVals = make([]string, len(values))
	for i := range values {
		b.mergeInsertVals[i] = "?"
	}
	b.args = append(b.args, values...)
	return b
}

func (b *queryBuilder) InsertInto(table string, columns ...string) SQLBuilder {
	b.kind = stringKindInsert
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
	b.kind = stringKindUpdate
	b.table = table
	return b
}

func (b *queryBuilder) Set(column string, value any) SQLBuilder {
	b.updateSet = append(b.updateSet, fmt.Sprintf("%s = ?", column))
	b.args = append(b.args, value)
	return b
}

func (b *queryBuilder) CreateTable(table string, definition string) SQLBuilder {
	b.kind = stringKindCreate
	b.special = fmt.Sprintf(createTableTemplate, table, definition)
	return b
}

func (b *queryBuilder) DropTable(table string) SQLBuilder {
	b.kind = stringKindDrop
	b.special = fmt.Sprintf(dropTableTemplate, table)
	return b
}

func (b *queryBuilder) DeleteFrom(table string) SQLBuilder {
	b.kind = stringKindDelete
	b.special = fmt.Sprintf(deleteTemplate, table)
	return b
}

func (b *queryBuilder) Select(table string, columns ...string) SQLBuilder {
	b.kind = stringKindSelect
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

func (b *queryBuilder) Union(other SQLBuilder) SQLBuilder {
	s1, a1 := b.Build()
	s2, a2 := other.Build()

	if strings.Contains(s1, "$1") && strings.Contains(s2, "$1") {
		s2 = strings.Replace(s2, "$1", "$2", 1)
	}

	combined := fmt.Sprintf("%s UNION %s", s1, s2)
	b.rawClauses = []string{combined}
	b.args = append(a1, a2...)
	return b
}

func (b *queryBuilder) Clear() SQLBuilder {
	q := &queryBuilder{
		opts:      b.opts,
		formatter: NewFormatter(b.opts.Driver),
	}
	*b = *q
	return b
}

type StructuredQuery struct {
	Kind            stringKinds       `json:"kind"`
	Columns         []string          `json:"columns,omitempty"`
	From            string            `json:"from,omitempty"`
	Where           string            `json:"where,omitempty"`
	GroupBy         []string          `json:"group_by,omitempty"`
	Having          []string          `json:"having,omitempty"`
	OrderBy         []string          `json:"order_by,omitempty"`
	Limit           *int              `json:"limit,omitempty"`
	Offset          *int              `json:"offset,omitempty"`
	Args            []any             `json:"args,omitempty"`
	Alias           string            `json:"alias,omitempty"`
	Joins           []string          `json:"joins,omitempty"`
	SQL             string            `json:"sql,omitempty"`
	Special         string            `json:"special,omitempty"`
	MergeTable      string            `json:"merge_table,omitempty"`
	MergeOn         string            `json:"merge_on,omitempty"`
	MergeMatchedSet []string          `json:"merge_matched,omitempty"`
	MergeInsertCols []string          `json:"merge_insert_cols,omitempty"`
	MergeInsertVals []string          `json:"merge_insert_vals,omitempty"`
	RawClauses      []string          `json:"raw,omitempty"`
	Queries         []StructuredQuery `json:"queries,omitempty"`
}

func (b *queryBuilder) ExportStructuredQuery() StructuredQuery {
	query, args := b.Build()
	return StructuredQuery{
		Kind:            b.kind,
		Columns:         b.columns,
		From:            b.table,
		Where:           strings.ReplaceAll(strings.Join(b.where, " AND "), "?", fmt.Sprintf("$%%d", 1)),
		GroupBy:         b.groupBy,
		Having:          b.having,
		OrderBy:         b.orderBy,
		Limit:           b.limit,
		Offset:          b.offset,
		Args:            args,
		Alias:           b.alias,
		Joins:           b.joins,
		SQL:             query,
		Special:         b.special,
		MergeTable:      b.mergeTable,
		MergeOn:         b.mergeOn,
		MergeMatchedSet: b.mergeMatchedSet,
		MergeInsertCols: b.mergeInsertCols,
		MergeInsertVals: b.mergeInsertVals,
		RawClauses:      b.rawClauses,
	}
}

// ExportAsJSON converts the builder's current query state to a JSON object
func (b *queryBuilder) ExportAsJSON() (string, error) {
	sq := b.ExportStructuredQuery()
	out, err := json.MarshalIndent(sq, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ExportAsXML converts the builder's current query state to an XML structure
func (b *queryBuilder) ExportAsXML() (string, error) {
	query, args := b.Build()
	data := struct {
		SQL  string `xml:"sql"`
		Args []any  `xml:"args>arg"`
	}{
		SQL:  query,
		Args: args,
	}
	out, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// StructToSQL converts a struct into SQL INSERT or UPDATE syntax using reflection (like GORM)
func (b *queryBuilder) StructToSQL(data any, table string, isInsert bool) (string, []any, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("StructToSQL expects a struct, got %s", v.Kind())
	}

	typeOf := v.Type()
	var columns []string
	var values []any
	var placeholders []string

	for i := 0; i < v.NumField(); i++ {
		field := typeOf.Field(i)
		tag := field.Tag.Get("db")
		if tag == "-" || tag == "" {
			continue
		}
		columns = append(columns, tag)
		values = append(values, v.Field(i).Interface())
		placeholders = append(placeholders, "?")
	}

	if isInsert {
		query := fmt.Sprintf(insertTemplate, table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
		return query, values, nil
	}

	set := make([]string, len(columns))
	for i, col := range columns {
		set[i] = fmt.Sprintf("%s = ?", col)
	}
	query := fmt.Sprintf(updateSetTemplate, table, strings.Join(set, ", "))
	return query, values, nil
}

func (b *queryBuilder) Build() (string, []any) {
	if b.mergeTable != "" {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("MERGE INTO %s", b.mergeTable))
		if b.mergeOn != "" {
			sb.WriteString(fmt.Sprintf(" ON %s", b.mergeOn))
		}
		if len(b.mergeMatchedSet) > 0 {
			sb.WriteString(" WHEN MATCHED THEN UPDATE SET ")
			sb.WriteString(strings.Join(b.mergeMatchedSet, ", "))
		}
		if len(b.mergeInsertCols) > 0 && len(b.mergeInsertVals) > 0 {
			sb.WriteString(" WHEN NOT MATCHED THEN INSERT (")
			sb.WriteString(strings.Join(b.mergeInsertCols, ", "))
			sb.WriteString(") VALUES (")
			sb.WriteString(strings.Join(b.mergeInsertVals, ", "))
			sb.WriteString(")")
		}
		return b.formatter.ReplacePlaceholders(sb.String()), b.args
	}

	if len(b.rawClauses) > 0 {
		return b.formatter.ReplacePlaceholders(strings.Join(b.rawClauses, " ")), b.args
	}

	if b.special != "" && strings.HasPrefix(b.special, "DELETE FROM") {
		query := b.special
		if len(b.where) > 0 {
			query += fmt.Sprintf(whereTemplate, strings.Join(b.where, " AND "))
		}
		return b.formatter.ReplacePlaceholders(query), b.args
	}

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

	if b.alias != "" {
		sb.WriteString(fmt.Sprintf(" AS %s", b.alias))
	}

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
