// Package dataprovider Package query SQLBuilder with adaptations based on SQL and DML/DDL best practices from educational sources and database dialect specifics
package dataprovider

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
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

const (
	joinTemplate              = "JOIN %s ON %s"
	leftJoinTemplate          = "LEFT JOIN %s ON %s"
	rightJoinTemplate         = "RIGHT JOIN %s ON %s"
	whereTemplate             = " WHERE %s"
	groupByTemplate           = "GROUP BY %s"
	orderByTemplate           = "ORDER BY %s"
	limitTemplate             = "LIMIT %d"
	offsetTemplate            = "OFFSET %d"
	havingTemplate            = "HAVING %s"
	selectTemplate            = "SELECT %s FROM %s"
	createTableTemplate       = "CREATE TABLE %s (%s)"
	dropTableTemplate         = "DROP TABLE %s"
	ifExistsDropTableTemplate = "DROP TABLE IF EXISTS %s"
	deleteTemplate            = "DELETE FROM %s"
	insertTemplate            = "INSERT INTO %s (%s) VALUES (%s)"
	updateTemplate            = "UPDATE %s SET %s"
	updateSetTemplate         = "UPDATE %s SET %s"
	andTemplate               = "(%s) AND (%s)"
)

type ColumnBuilder struct {
	tb       *TableBuilder
	colIndex int
}

func (c *ColumnBuilder) setGoType(goType string) *ColumnBuilder {
	col := &c.tb.Columns[c.colIndex]
	col.Type = c.tb.Dialect.TypeMap(goType)
	return c
}

func (c *ColumnBuilder) Int() *ColumnBuilder {
	return c.setGoType("int")
}

func (c *ColumnBuilder) Float64() *ColumnBuilder {
	return c.setGoType("float64")
}

func (c *ColumnBuilder) String() *ColumnBuilder {
	return c.setGoType("string")
}

func (c *ColumnBuilder) Bool() *ColumnBuilder {
	return c.setGoType("bool")
}

func (c *ColumnBuilder) Timestamp() *ColumnBuilder {
	return c.setGoType("time.Time")
}

func (c *ColumnBuilder) NotNull() *ColumnBuilder {
	c.tb.Columns[c.colIndex].NotNull = true
	return c
}

func (c *ColumnBuilder) PrimaryKey() *ColumnBuilder {
	c.tb.Columns[c.colIndex].PrimaryKey = true
	return c
}

func (c *ColumnBuilder) Default(value string) *ColumnBuilder {
	c.tb.Columns[c.colIndex].DefaultValue = value
	return c
}

func (c *ColumnBuilder) Unique() *ColumnBuilder {
	c.tb.Columns[c.colIndex].Unique = true
	return c
}

func (c *ColumnBuilder) ForeignKey(refTable, refColumn, onDelete, onUpdate string) *ColumnBuilder {
	c.tb.Columns[c.colIndex].ForeignKey = &ForeignKey{
		Column:    c.tb.Columns[c.colIndex].Name,
		RefTable:  refTable,
		RefColumn: refColumn,
		OnDelete:  onDelete,
		OnUpdate:  onUpdate,
	}
	return c
}

func (c *ColumnBuilder) Index() *ColumnBuilder {
	col := c.tb.Columns[c.colIndex]
	idxName := fmt.Sprintf("idx_%s_%s", strings.ToLower(c.tb.Name), strings.ToLower(col.Name))
	c.tb.Indexes = append(c.tb.Indexes, Index{
		Name:   idxName,
		Column: col.Name,
	})
	c.tb.Columns[c.colIndex].HasIndex = true
	return c
}

type ForeignKey struct {
	Column    string
	RefTable  string
	RefColumn string
	OnDelete  string
	OnUpdate  string
}

type Index struct {
	Name   string
	Column string
}

type TableBuilder struct {
	Name    string
	Columns []Column
	Indexes []Index
	Dialect
}

func (b *TableBuilder) CreateTable(name string) *TableBuilder {
	b.Name = name
	return b
}

func (b *TableBuilder) Column(name string) *ColumnBuilder {
	b.Columns = append(b.Columns, Column{Name: name})
	return &ColumnBuilder{tb: b, colIndex: len(b.Columns) - 1}
}

func (b *TableBuilder) NotNull() *TableBuilder {
	b.Columns[len(b.Columns)-1].NotNull = true
	return b
}

func (b *TableBuilder) PrimaryKey() *TableBuilder {
	b.Columns[len(b.Columns)-1].PrimaryKey = true
	return b
}

func (b *TableBuilder) ForeignKey(refTable, refColumn string) *TableBuilder {
	col := &b.Columns[len(b.Columns)-1]
	col.ForeignKey = &ForeignKey{
		Column:    col.Name,
		RefTable:  refTable,
		RefColumn: refColumn,
	}
	return b
}

func (b *TableBuilder) Index(unique bool) *TableBuilder {
	col := b.Columns[len(b.Columns)-1]
	idxName := fmt.Sprintf("idx_%s_%s", strings.ToLower(b.Name), strings.ToLower(col.Name))
	b.Indexes = append(b.Indexes, Index{
		Name:   idxName,
		Column: col.Name,
	})
	b.Columns[len(b.Columns)-1].HasIndex = true
	b.Columns[len(b.Columns)-1].UniqueIndex = unique
	return b
}

func (b *TableBuilder) Build() string {
	var lines []string
	for _, c := range b.Columns {
		col := fmt.Sprintf("%s %s", b.Dialect.QuoteIdentifier(c.Name), b.Dialect.TypeMap(c.Type))
		if c.PrimaryKey {
			col += " " + b.Dialect.PrimaryKeySyntax()
		}
		if c.NotNull {
			col += " NOT NULL"
		}
		lines = append(lines, col)
	}

	for _, c := range b.Columns {
		if c.ForeignKey != nil {
			lines = append(lines, fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)",
				b.Dialect.QuoteIdentifier(c.Name),
				b.Dialect.QuoteIdentifier(c.ForeignKey.RefTable),
				b.Dialect.QuoteIdentifier(c.ForeignKey.RefColumn)))
		}
	}

	createStmt := fmt.Sprintf("CREATE TABLE %s (\n  %s\n);",
		b.Dialect.QuoteIdentifier(b.Name),
		strings.Join(lines, ",\n  "))

	var indexStmts []string
	for _, idx := range b.Indexes {
		indexStmts = append(indexStmts, fmt.Sprintf("CREATE INDEX %s ON %s(%s);",
			b.Dialect.QuoteIdentifier(idx.Name),
			b.Dialect.QuoteIdentifier(b.Name),
			b.Dialect.QuoteIdentifier(idx.Column)))
	}

	return createStmt + "\n" + strings.Join(indexStmts, "\n")
}

func (b *TableBuilder) ToFile(path string) error {
	sql := b.Build()
	return os.WriteFile(path, []byte(sql), 0644)
}

// SQLBuilder interface models typical SQL DDL and DML operations for various dialects
type SQLBuilder interface {
	Select(table string, columns ...string) SQLBuilder
	Where(condition string, args ...any) SQLBuilder
	And(condition string, args ...any) SQLBuilder
	Join(table, onCondition string) SQLBuilder
	LeftJoin(table, onCondition string) SQLBuilder
	RightJoin(table, onCondition string) SQLBuilder
	GroupBy(columns ...string) SQLBuilder
	Having(condition string, args ...any) SQLBuilder
	OrderBy(columns ...string) SQLBuilder
	Limit(n int) SQLBuilder
	Offset(n int) SQLBuilder
	CreateTable(table string) SQLBuilder
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
	ExportAsJSON(pretty bool) (string, error)
	ExportAsXML() (string, error)
	ExportAsYAML() (string, error)
	StructToSQL(data any, table string, isInsert bool) (string, []any, error)
	IfNotExists() SQLBuilder
	IfExists() SQLBuilder
	ExportStructuredQuery() StructuredQuery
	ImportStructuredQuery(s StructuredQuery) SQLBuilder
	Columns(columns ...Column) SQLBuilder
	Column(name string, colType string, opts ...func(*Column)) SQLBuilder
}

type Column struct {
	Name         string
	Type         string
	NotNull      bool
	PrimaryKey   bool
	ForeignKey   *ForeignKey
	HasIndex     bool
	UniqueIndex  bool
	DefaultValue any
	Unique       bool
}

type queryBuilder struct {
	opts            *Options
	kind            stringKinds
	table           string
	selectCols      []string
	columns         []Column
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
	ifNotExists     bool
	ifExists        bool
}

func NewQueryBuilder(opts *Options) SQLBuilder {
	return &queryBuilder{
		opts:      opts,
		formatter: NewFormatter(opts.Driver),
	}
}

func (b *queryBuilder) As(alias string) SQLBuilder {
	b.alias = alias
	return b
}

func (b *queryBuilder) And(condition string, args ...any) SQLBuilder {
	b.where = append(b.where, fmt.Sprintf(andTemplate, b.where[len(b.where)-1], condition))
	b.args = append(b.args, args...)
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

func (b *queryBuilder) CreateTable(table string) SQLBuilder {
	b.kind = stringKindCreate
	b.table = table
	return b
}

func (b *queryBuilder) DropTable(table string) SQLBuilder {
	b.kind = stringKindDrop
	if b.ifExists {
		b.special = fmt.Sprintf(ifExistsDropTableTemplate, table)
		return b
	}
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
	b.selectCols = columns
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

// Offset sets the offset for the query.
func (b *queryBuilder) Offset(n int) SQLBuilder {
	b.offset = &n
	return b
}

// Union combines two queries into a single UNION query
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

// Clear resets the builder to its initial state
func (b *queryBuilder) Clear() SQLBuilder {
	*b = queryBuilder{opts: b.opts, formatter: NewFormatter(b.opts.Driver)}
	return b
}

type StructuredQuery struct {
	XMLName         xml.Name          `json:"-" xml:"query"`
	Kind            stringKinds       `json:"kind" xml:"kind"`
	Columns         []Column          `json:"columns,omitempty" xml:"columns"`
	From            string            `json:"from,omitempty" xml:"from"`
	Where           string            `json:"where,omitempty" xml:"where"`
	GroupBy         []string          `json:"groupBy,omitempty" xml:"groupBy"`
	Having          []string          `json:"having,omitempty" xml:"having"`
	OrderBy         []string          `json:"orderBy,omitempty" xml:"orderBy"`
	Limit           *int              `json:"limit,omitempty" xml:"limit"`
	Offset          *int              `json:"offset,omitempty" xml:"offset"`
	Alias           string            `json:"alias,omitempty" xml:"alias"`
	Joins           []string          `json:"joins,omitempty" xml:"joins"`
	SQL             string            `json:"sql,omitempty" xml:"SQL"`
	Special         string            `json:"special,omitempty" xml:"special"`
	MergeTable      string            `json:"mergeTable,omitempty" xml:"mergeTable"`
	MergeOn         string            `json:"mergeOn,omitempty" xml:"mergeOn"`
	MergeMatchedSet []string          `json:"mergeMatchedSet,omitempty" xml:"mergeMatchedSet"`
	MergeInsertCols []string          `json:"mergeInsertCols,omitempty" xml:"mergeInsertCols"`
	MergeInsertVals []string          `json:"mergeInsertVals,omitempty" xml:"mergeInsertVals"`
	RawClauses      []string          `json:"raw,omitempty" xml:"rawClauses"`
	Queries         []StructuredQuery `json:"queries,omitempty" xml:"queries"`
	Args            []any             `json:"args,omitempty" xml:"args"`
	Data            any               `json:"data,omitempty" xml:"data"`
}

// ExportStructuredQuery decomposes the builder into a StructuredQuery
func (b *queryBuilder) ExportStructuredQuery() StructuredQuery {
	formatStrings := func(clauses []string) []string {
		formatted := make([]string, len(clauses))
		for i, clause := range clauses {
			formatted[i] = NewFormatter(b.opts.Driver).ReplacePlaceholders(clause)
		}
		return formatted
	}

	query, args := b.Build()
	formatter := NewFormatter(b.opts.Driver)
	return StructuredQuery{
		Kind:            b.kind,
		Columns:         b.columns,
		From:            b.table,
		Where:           formatter.ReplacePlaceholders(strings.Join(b.where, " AND ")),
		GroupBy:         b.groupBy,
		Having:          formatStrings(b.having),
		OrderBy:         b.orderBy,
		Limit:           b.limit,
		Offset:          b.offset,
		Args:            args,
		Alias:           b.alias,
		Joins:           b.joins,
		SQL:             query,
		Special:         b.special,
		MergeTable:      b.mergeTable,
		MergeOn:         formatter.ReplacePlaceholders(b.mergeOn),
		MergeMatchedSet: formatStrings(b.mergeMatchedSet),
		MergeInsertCols: b.mergeInsertCols,
		MergeInsertVals: b.mergeInsertVals,
		RawClauses:      formatStrings(b.rawClauses),
		XMLName:         xml.Name{Local: "query"},
	}
}

// ImportStructuredQuery applies a StructuredQuery to a queryBuilder
func (b *queryBuilder) ImportStructuredQuery(s StructuredQuery) SQLBuilder {
	b.columns = s.Columns
	b.table = s.From
	if s.Where != "" {
		b.Where(s.Where)
	}
	b.groupBy = s.GroupBy
	b.having = s.Having
	b.orderBy = s.OrderBy
	b.limit = s.Limit
	b.offset = s.Offset
	b.args = s.Args
	b.alias = s.Alias
	b.joins = s.Joins
	b.special = s.Special
	b.mergeTable = s.MergeTable
	b.mergeOn = s.MergeOn
	b.mergeMatchedSet = s.MergeMatchedSet
	b.mergeInsertCols = s.MergeInsertCols
	b.mergeInsertVals = s.MergeInsertVals
	b.rawClauses = s.RawClauses
	return b
}

// ExportAsJSON converts the builder's current query state to a JSON object
func (b *queryBuilder) ExportAsJSON(pretty bool) (string, error) {
	sq := b.ExportStructuredQuery()
	var out []byte
	var err error
	if pretty {
		out, err = json.MarshalIndent(sq, "", "  ")
	} else {
		out, err = json.Marshal(sq)
	}
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ExportAsXML converts the builder's current query state to an XML structure
func (b *queryBuilder) ExportAsXML() (string, error) {
	sq := b.ExportStructuredQuery()
	out, err := xml.Marshal(sq)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ExportAsYAML converts the builder's current query state to an YAML structure
func (b *queryBuilder) ExportAsYAML() (string, error) {
	sq := b.ExportStructuredQuery()
	out, err := yaml.Marshal(sq)
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

// StructToSQLWithPK converts a struct into SQL INSERT or UPDATE syntax using reflection (like GORM)
// and sets the primary key column to the value of the primary key field
func (b *queryBuilder) StructToSQLWithPK(data any, table string, isInsert bool) (string, []any, error) {
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

	pk := typeOf.Field(0).Tag.Get("db")
	if pk == "" {
		return "", nil, fmt.Errorf("StructToSQLWithPK expects a struct with a primary key field, got %s", typeOf.Field(0).Name)
	}

	return query, append(values, data.(map[string]any)[pk]), nil
}

func (b *queryBuilder) Column(name string, colType string, opts ...func(*Column)) SQLBuilder {
	col := Column{Name: name, Type: colType}
	for _, opt := range opts {
		opt(&col)
	}
	b.columns = append(b.columns, col)
	return b
}

func (b *queryBuilder) Columns(cols ...Column) SQLBuilder {
	b.columns = append(b.columns, cols...)
	return b
}

func (b *queryBuilder) IfNotExists() SQLBuilder {
	b.ifNotExists = true
	return b
}

func (b *queryBuilder) IfExists() SQLBuilder {
	b.ifExists = true
	return b
}

// Build builds the query and returns the query string and a slice of arguments
func (b *queryBuilder) Build() (string, []any) {
	var sb strings.Builder

	switch b.kind {
	case stringKindSelect:
		sb.WriteString("SELECT ")
		if len(b.selectCols) > 0 {
			sb.WriteString(strings.Join(b.selectCols, ", "))
		} else {
			sb.WriteString("*")
		}
		sb.WriteString(" FROM " + b.table)

		if len(b.joins) > 0 {
			sb.WriteString(" " + strings.Join(b.joins, " "))
		}
		if len(b.where) > 0 {
			sb.WriteString(" WHERE " + strings.Join(b.where, " AND "))
		}
		if len(b.groupBy) > 0 {
			sb.WriteString(" GROUP BY " + strings.Join(b.groupBy, ", "))
		}
		if len(b.having) > 0 {
			sb.WriteString(" HAVING " + strings.Join(b.having, " AND "))
		}
		if len(b.orderBy) > 0 {
			sb.WriteString(" ORDER BY " + strings.Join(b.orderBy, ", "))
		}
		if b.limit != nil {
			sb.WriteString(fmt.Sprintf(" LIMIT %d", *b.limit))
		}
		if b.offset != nil {
			sb.WriteString(fmt.Sprintf(" OFFSET %d", *b.offset))
		}
		return sb.String(), b.args

	case stringKindCreate:
		var parts []string
		for _, col := range b.columns {
			def := fmt.Sprintf("%s %s", b.formatter.QuoteIdentifier(col.Name), col.Type)
			if col.PrimaryKey {
				def += " PRIMARY KEY"
			}
			if col.NotNull {
				def += " NOT NULL"
			}
			parts = append(parts, def)
		}
		stmt := "CREATE TABLE"
		if b.ifNotExists {
			stmt += " IF NOT EXISTS"
		}
		stmt += fmt.Sprintf(" %s (%s)", b.table, strings.Join(parts, ", "))
		return stmt, nil

	case stringKindInsert:
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			b.table,
			strings.Join(b.insertCols, ", "),
			strings.Join(b.insertVals, ", "))
		return query, b.args

	case stringKindUpdate:
		query := fmt.Sprintf("UPDATE %s SET %s", b.table, strings.Join(b.updateSet, ", "))
		if len(b.where) > 0 {
			query += fmt.Sprintf(" WHERE %s", strings.Join(b.where, " AND "))
		}
		return query, b.args

	case stringKindDelete:
		query := fmt.Sprintf("DELETE FROM %s", b.table)
		if len(b.where) > 0 {
			query += fmt.Sprintf(" WHERE %s", strings.Join(b.where, " AND "))
		}
		return query, b.args

	case stringKindDrop:
		query := "DROP TABLE"
		if b.ifExists {
			query += " IF EXISTS"
		}
		query += " " + b.table
		return query, nil

	default:
		return "", nil
	}
}

func NotNull() func(*Column) {
	return func(c *Column) { c.NotNull = true }
}

func PrimaryKey() func(*Column) {
	return func(c *Column) { c.PrimaryKey = true }
}

func Unique() func(*Column) {
	return func(c *Column) { c.Unique = true }
}

func Default(val any) func(*Column) {
	return func(c *Column) { c.DefaultValue = val }
}
