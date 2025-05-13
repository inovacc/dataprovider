package dataprovider

import (
	"fmt"
	"testing"
)

func TestNewDialects(t *testing.T) {
	opt := &Options{
		Driver: PostgresSQLDatabaseProviderName,
	}

	dialect := NewDialect(opt)
	if dialect == nil {
		t.Error("Dialect not found")
	}

	fmt.Println(dialect.TypeMap("string"))
}
