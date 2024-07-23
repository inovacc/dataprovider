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

func TestQueryCreate(t *testing.T) {
	provider := Must(NewDataProvider(NewOptions(WithMemoryDB())))

	query := provider.SqlBuilder().
		CreateTable("users").
		IfNotExists().
		Columns(
			provider.SqlBuilder().Column("id").Type("SERIAL").PrimaryKey(),
			provider.SqlBuilder().Column("first_name").Type("TEXT"),
			provider.SqlBuilder().Column("last_name").Type("TEXT"),
			provider.SqlBuilder().Column("email").Type("TEXT"),
			provider.SqlBuilder().Column("ip_address").Type("TEXT"),
			provider.SqlBuilder().Column("city").Type("TEXT"),
		).
		Build()

	t.Log("CREATE TABLE Query:")
	t.Log(query)

	if query != "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, first_name TEXT, last_name TEXT, email TEXT, ip_address TEXT, city TEXT)" {
		t.Error("CREATE TABLE Query is not correct")
	}
}

func TestQueryUpdate(t *testing.T) {
	provider := Must(NewDataProvider(NewOptions(WithMemoryDB())))

	query := provider.SqlBuilder().
		Table("users").
		Update().
		Set(map[string]any{
			"name":  "John Doe",
			"email": "test@admin.com",
		}).
		Where("id = 1").
		Build()

	t.Log("UPDATE Query:")
	t.Log(query)

	if query != "UPDATE users SET name = 'John Doe', email = 'test@admin.com' WHERE id = 1" {
		t.Error("UPDATE Query is not correct")
	}
}
