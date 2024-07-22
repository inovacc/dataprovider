package dataprovider

import (
	"testing"
)

func TestQuery(t *testing.T) {
	provider := Must(NewDataProvider(NewOptions(WithMemoryDB())))

	query := provider.SqlBuilder().
		Table("users").
		Select("id", "name", "email").
		Where("age > 18").
		OrderBy("name ASC").
		Limit(10).
		Offset(5).
		Build()

	t.Log("SELECT Query:")
	t.Log(query)

	if query != "SELECT id, name, email FROM users WHERE age > 18 ORDER BY name ASC LIMIT 10 OFFSET 5" {
		t.Error("SELECT Query is not correct")
	}
}

func TestQuerySchema(t *testing.T) {
	provider := Must(NewDataProvider(NewOptions(WithMemoryDB())))

	query := provider.SqlBuilder().
		Schema("public").
		Table("users").
		Select("id", "name", "email").
		Where("age > 18").
		OrderBy("name ASC").
		Limit(10).
		Offset(5).
		Build()

	t.Log("SELECT Query:")
	t.Log(query)

	if query != "SELECT id, name, email FROM public.users WHERE age > 18 ORDER BY name ASC LIMIT 10 OFFSET 5" {
		t.Error("SELECT Query is not correct")
	}
}
