package dataprovider

import (
	"fmt"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	b := NewBuilder(PostgresDialect{})
	sql := b.CreateTable("users").
		Column("id", b.Dialect.AutoIncrement()).PrimaryKey().
		Column("email", "varchar(255)").NotNull().
		Column("created_at", "timestamp").
		Build()

	fmt.Println(sql)

}
