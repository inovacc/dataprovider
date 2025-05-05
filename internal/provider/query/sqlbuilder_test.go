package query

import (
	"testing"

	"github.com/inovacc/dataprovider/internal/provider"
)

func TestQueryBuilderVariants(t *testing.T) {
	tests := []struct {
		driver      string
		expectedSQL string
	}{
		{
			driver:      provider.OracleDatabaseProviderName,
			expectedSQL: "SELECT id, name FROM users WHERE status = :p1 ORDER BY name LIMIT 10",
		},
		{
			driver:      provider.PostgresSQLDatabaseProviderName,
			expectedSQL: "SELECT id, name FROM users WHERE status = $1 ORDER BY name LIMIT 10",
		},
		{
			driver:      provider.MySQLDatabaseProviderName,
			expectedSQL: "SELECT id, name FROM users WHERE status = ? ORDER BY name LIMIT 10",
		},
		{
			driver:      provider.SQLiteDataProviderName,
			expectedSQL: "SELECT id, name FROM users WHERE status = ? ORDER BY name LIMIT 10",
		},
		{
			driver:      "mariadb", // assuming MariaDB shares behavior with MySQL
			expectedSQL: "SELECT id, name FROM users WHERE status = ? ORDER BY name LIMIT 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.driver, func(t *testing.T) {
			opts := provider.Options{Driver: tt.driver}
			qb := NewQueryBuilder(opts)
			sql, _ := qb.Select("users", "id", "name").
				Where("status = ?", "active").
				OrderBy("name").
				Limit(10).
				Build()

			if sql != tt.expectedSQL {
				t.Errorf("driver %s: expected %q, got %q", tt.driver, tt.expectedSQL, sql)
			}
		})
	}
}

func TestCreateDropDeleteSQL(t *testing.T) {
	tests := []struct {
		driver       string
		operation    string
		builderFunc  func(SQLBuilder) (string, []any)
		expectedSQLs map[string]string
	}{
		{
			driver:    provider.OracleDatabaseProviderName,
			operation: "create",
			builderFunc: func(b SQLBuilder) (string, []any) {
				return b.CreateTable("users", "id INT, name VARCHAR2(50)").Build()
			},
			expectedSQLs: map[string]string{
				"oracle": "CREATE TABLE users (id INT, name VARCHAR2(50))",
			},
		},
		{
			driver:    provider.PostgresSQLDatabaseProviderName,
			operation: "drop",
			builderFunc: func(b SQLBuilder) (string, []any) {
				return b.DropTable("users").Build()
			},
			expectedSQLs: map[string]string{
				"postgres": "DROP TABLE users",
			},
		},
		{
			driver:    provider.MySQLDatabaseProviderName,
			operation: "delete",
			builderFunc: func(b SQLBuilder) (string, []any) {
				return b.DeleteFrom("users").Where("name = ?", "john").Build()
			},
			expectedSQLs: map[string]string{
				"mysql": "DELETE FROM users WHERE name = ?",
			},
		},
		{
			driver:    provider.SQLiteDataProviderName,
			operation: "delete",
			builderFunc: func(b SQLBuilder) (string, []any) {
				return b.DeleteFrom("users").Where("name = ?", "john").Build()
			},
			expectedSQLs: map[string]string{
				"sqlite": "DELETE FROM users WHERE name = ?",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.operation+"_"+tt.driver, func(t *testing.T) {
			opts := provider.Options{Driver: tt.driver}
			builder := NewQueryBuilder(opts)
			sql, _ := tt.builderFunc(builder)

			key := tt.driver
			if val, ok := tt.expectedSQLs[key]; ok && sql != val {
				t.Errorf("expected SQL %q, got %q", val, sql)
			}
		})
	}
}

func TestInsertAndUpdateSQL(t *testing.T) {
	tests := []struct {
		driver       string
		opType       string
		builderFunc  func(SQLBuilder) (string, []any)
		expectedSQLs map[string]string
	}{
		{
			driver: "postgres",
			opType: "insert",
			builderFunc: func(b SQLBuilder) (string, []any) {
				return b.InsertInto("users", "name", "email").
					Values("john", "john@example.com").Build()
			},
			expectedSQLs: map[string]string{
				"postgres": "INSERT INTO users (name, email) VALUES ($1, $2)",
			},
		},
		{
			driver: "oracle",
			opType: "insert",
			builderFunc: func(b SQLBuilder) (string, []any) {
				return b.InsertInto("accounts", "username", "balance").
					Values("alice", 1000).Build()
			},
			expectedSQLs: map[string]string{
				"oracle": "INSERT INTO accounts (username, balance) VALUES (:p1, :p2)",
			},
		},
		{
			driver: "postgres",
			opType: "update",
			builderFunc: func(b SQLBuilder) (string, []any) {
				return b.Update("products").Set("price", 9.99).Where("id = ?", 42).Build()
			},
			expectedSQLs: map[string]string{
				"postgres": "UPDATE products SET price = $1 WHERE id = $2",
			},
		},
		{
			driver: "oracle",
			opType: "update",
			builderFunc: func(b SQLBuilder) (string, []any) {
				return b.Update("inventory").Set("stock", 30).Where("item_id = ?", 7).Build()
			},
			expectedSQLs: map[string]string{
				"oracle": "UPDATE inventory SET stock = :p1 WHERE item_id = :p2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.opType+"_"+tt.driver, func(t *testing.T) {
			opts := provider.Options{Driver: tt.driver}
			builder := NewQueryBuilder(opts)
			sql, _ := tt.builderFunc(builder)

			expected := tt.expectedSQLs[tt.driver]
			if sql != expected {
				t.Errorf("expected SQL %q, got %q", expected, sql)
			}
		})
	}
}
