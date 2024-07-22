package provider

import (
	"fmt"
	"strings"
)

// SQLBuilder estructura para construir consultas SQL
type SQLBuilder struct {
	queryType string
	schema    string
	table     string
	driver    string
	columns   []string
	values    []string
	where     []string
	joins     []string
	groupBy   []string
	orderBy   []string
	limit     int
	offset    int
}

// NewSQLBuilder crea una nueva instancia de SQLBuilder
func NewSQLBuilder(driver string) *SQLBuilder {
	return &SQLBuilder{driver: driver}
}

// Schema establece el esquema de la tabla
func (b *SQLBuilder) Schema(schema string) *SQLBuilder {
	b.schema = schema
	return b
}

func (b *SQLBuilder) Table(table string) *SQLBuilder {
	b.table = table
	return b
}

// Select establece las columnas a seleccionar
func (b *SQLBuilder) Select(columns ...string) *SQLBuilder {
	b.queryType = "SELECT"
	b.columns = columns
	return b
}

// Insert establece la tabla y columnas para la consulta INSERT
func (b *SQLBuilder) Insert(columns ...string) *SQLBuilder {
	b.queryType = "INSERT"
	b.columns = columns
	return b
}

// Values establece los valores para la consulta INSERT
func (b *SQLBuilder) Values(values ...string) *SQLBuilder {
	b.values = values
	return b
}

// Where agrega una condición WHERE
func (b *SQLBuilder) Where(condition string) *SQLBuilder {
	b.where = append(b.where, condition)
	return b
}

// Join agrega una cláusula JOIN
func (b *SQLBuilder) Join(join string) *SQLBuilder {
	b.joins = append(b.joins, join)
	return b
}

// GroupBy establece la cláusula GROUP BY
func (b *SQLBuilder) GroupBy(columns ...string) *SQLBuilder {
	b.groupBy = columns
	return b
}

// OrderBy establece la cláusula ORDER BY
func (b *SQLBuilder) OrderBy(columns ...string) *SQLBuilder {
	b.orderBy = columns
	return b
}

// Limit establece la cláusula LIMIT
func (b *SQLBuilder) Limit(limit int) *SQLBuilder {
	b.limit = limit
	return b
}

// Offset establece la cláusula OFFSET
func (b *SQLBuilder) Offset(offset int) *SQLBuilder {
	b.offset = offset
	return b
}

// Build construye la consulta SQL
func (b *SQLBuilder) Build() string {
	var sb strings.Builder

	switch b.queryType {
	case "SELECT":
		sb.WriteString("SELECT ")
		if len(b.columns) > 0 {
			sb.WriteString(strings.Join(b.columns, ", "))
		} else {
			sb.WriteString("*")
		}
		sb.WriteString(" FROM ")
		if b.schema != "" {
			sb.WriteString(b.schema)
			sb.WriteString(".")
		}
		sb.WriteString(b.table)
	case "INSERT":
		sb.WriteString("INSERT INTO ")
		if b.schema != "" {
			sb.WriteString(b.schema)
			sb.WriteString(".")
		}
		sb.WriteString(b.table)
		sb.WriteString(" (")
		sb.WriteString(strings.Join(b.columns, ", "))
		sb.WriteString(") VALUES (")
		sb.WriteString(strings.Join(b.values, ", "))
		sb.WriteString(")")
	default:
		return ""
	}

	if len(b.joins) > 0 {
		for _, join := range b.joins {
			sb.WriteString(" ")
			sb.WriteString(join)
		}
	}

	if len(b.where) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(b.where, " AND "))
	}

	if len(b.groupBy) > 0 {
		sb.WriteString(" GROUP BY ")
		sb.WriteString(strings.Join(b.groupBy, ", "))
	}

	if len(b.orderBy) > 0 {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(strings.Join(b.orderBy, ", "))
	}

	switch b.driver {
	case OracleDatabaseProviderName:
		if b.offset > 0 || b.limit > 0 {
			sb.WriteString(fmt.Sprintf(" OFFSET %d ROWS", b.offset))
			if b.limit > 0 {
				sb.WriteString(fmt.Sprintf(" FETCH NEXT %d ROWS ONLY", b.limit))
			}
		}
	default:
		if b.limit > 0 {
			sb.WriteString(fmt.Sprintf(" LIMIT %d", b.limit))
		}
		if b.offset > 0 {
			sb.WriteString(fmt.Sprintf(" OFFSET %d", b.offset))
		}
	}

	return sb.String()
}

// Function define una función
func (b *SQLBuilder) Function(name string, params ...string) string {
	var sb strings.Builder
	sb.WriteString("CREATE FUNCTION ")
	sb.WriteString(name)
	sb.WriteString("(")
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") RETURNS void AS $$ BEGIN RETURN; END; $$ LANGUAGE plpgsql;")
	return sb.String()
}

// Procedure define un procedimiento almacenado
func (b *SQLBuilder) Procedure(name string, params ...string) string {
	var sb strings.Builder
	sb.WriteString("CREATE PROCEDURE ")
	sb.WriteString(name)
	sb.WriteString("(")
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") AS $$ BEGIN RETURN; END; $$ LANGUAGE plpgsql;")
	return sb.String()
}
