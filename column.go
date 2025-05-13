package dataprovider

import (
	"fmt"
	"strings"
)

type Column struct {
	Name        string
	Type        string
	NotNull     bool
	PrimaryKey  bool
	ForeignKey  *ForeignKey
	HasIndex    bool
	UniqueIndex bool
}

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
