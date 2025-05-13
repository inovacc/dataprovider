package dataprovider

import (
	"fmt"
	"os"
	"strings"
)

type ForeignKey struct {
	Column    string
	RefTable  string
	RefColumn string
}

type Index struct {
	Name   string
	Column string
}

type TableBuilder struct {
	Name    string
	Columns []Column
	Indexes []Index
	Dialect Dialect
}

func NewBuilder(d Dialect) *TableBuilder {
	return &TableBuilder{Dialect: d}
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

func main() {
	builder := NewBuilder(PostgresDialect{})
	builder.CreateTable("enrollment").
		Column("student_id").Int().NotNull().PrimaryKey().
		Column("grade").Float64().
		Column("comment").String().Default("'none'").NotNull().
		Column("verified").Bool().Default("false")

	sql := builder.Build()
	fmt.Println(sql)

	// Exportar a archivo
	if err := builder.ToFile("orders_postgres.sql"); err != nil {
		fmt.Println("❌ Error al exportar archivo:", err)
	}
}
