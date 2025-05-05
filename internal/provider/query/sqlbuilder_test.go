package query

import (
	"fmt"
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
		t.Run(fmt.Sprintf("%s_%s", tt.operation, tt.driver), func(t *testing.T) {
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
		t.Run(fmt.Sprintf("%s_%s", tt.opType, tt.driver), func(t *testing.T) {
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

func TestMergeBuilder(t *testing.T) {
	tests := []struct {
		driver   string
		expected string
	}{
		{
			driver:   provider.PostgresSQLDatabaseProviderName,
			expected: "MERGE INTO users ON id = $1 WHEN MATCHED THEN UPDATE SET email = $2 WHEN NOT MATCHED THEN INSERT (id, email) VALUES ($3, $4)",
		},
		{
			driver:   provider.OracleDatabaseProviderName,
			expected: "MERGE INTO users ON id = :p1 WHEN MATCHED THEN UPDATE SET email = :p2 WHEN NOT MATCHED THEN INSERT (id, email) VALUES (:p3, :p4)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.driver, func(t *testing.T) {
			opts := provider.Options{Driver: tt.driver}
			q := NewQueryBuilder(opts).
				MergeInto("users").
				On("id = ?").
				WhenMatched(map[string]any{"email": "updated@example.com"}).
				WhenNotMatchedInsert([]string{"id", "email"}, []any{123, "new@example.com"})

			sql, _ := q.Build()
			if sql != tt.expected {
				t.Errorf("Expected: %q\nGot:      %q", tt.expected, sql)
			}
		})
	}
}

func TestRawSQLInjection(t *testing.T) {
	opts := provider.Options{Driver: provider.PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opts).
		Raw("SELECT * FROM users WHERE email = ?", "john@example.com")

	sql, args := q.Build()
	expectedSQL := "SELECT * FROM users WHERE email = $1"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
	if len(args) != 1 || args[0] != "john@example.com" {
		t.Errorf("Expected args to be [\"john@example.com\"], got %v", args)
	}
}

func TestAliasSupport(t *testing.T) {
	opts := provider.Options{Driver: provider.PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opts).
		Select("users", "id", "email").
		As("u")

	sql, _ := q.Build()
	expectedSQL := "SELECT id, email FROM users AS u"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
}

func TestGroupByHaving(t *testing.T) {
	opts := provider.Options{Driver: provider.PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opts).
		Select("orders", "customer_id", "COUNT(*)").
		GroupBy("customer_id").
		Having("COUNT(*) > ?", 5).
		OrderBy("customer_id")

	sql, args := q.Build()
	expectedSQL := "SELECT customer_id, COUNT(*) FROM orders GROUP BY customer_id HAVING COUNT(*) > $1 ORDER BY customer_id"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
	if len(args) != 1 || args[0] != 5 {
		t.Errorf("Expected args to be [5], got %v", args)
	}
}

func TestJoinClauses(t *testing.T) {
	opts := provider.Options{Driver: provider.PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opts).
		Select("orders", "orders.id", "users.email").
		Join("users", "orders.user_id = users.id").
		LeftJoin("payments", "orders.id = payments.order_id").
		RightJoin("shipments", "orders.id = shipments.order_id")

	sql, _ := q.Build()
	expectedSQL := "SELECT orders.id, users.email FROM orders JOIN users ON orders.user_id = users.id LEFT JOIN payments ON orders.id = payments.order_id RIGHT JOIN shipments ON orders.id = shipments.order_id"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
}

func TestJoinGroupByHaving(t *testing.T) {
	opts := provider.Options{Driver: provider.PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opts).
		Select("orders", "users.name", "COUNT(orders.id)").
		Join("users", "orders.user_id = users.id").
		GroupBy("users.name").
		Having("COUNT(orders.id) > ?", 10).
		OrderBy("users.name")

	sql, args := q.Build()
	expectedSQL := "SELECT users.name, COUNT(orders.id) FROM orders JOIN users ON orders.user_id = users.id GROUP BY users.name HAVING COUNT(orders.id) > $1 ORDER BY users.name"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
	if len(args) != 1 || args[0] != 10 {
		t.Errorf("Expected args to be [10], got %v", args)
	}
}

func TestNestedSelect(t *testing.T) {
	opts := provider.Options{Driver: provider.PostgresSQLDatabaseProviderName}
	inner := NewQueryBuilder(opts).
		Select("payments", "user_id").
		Where("status = ?", "completed")

	subSQL, subArgs := inner.Build()

	q := NewQueryBuilder(opts).
		Select("users", "id", "email").
		Where(fmt.Sprintf("id IN (%s)", subSQL), subArgs...)

	sql, args := q.Build()
	expectedSQL := "SELECT id, email FROM users WHERE id IN (SELECT user_id FROM payments WHERE status = $1)"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
	if len(args) != 1 || args[0] != "completed" {
		t.Errorf("Expected args to be [\"completed\"], got %v", args)
	}
}

func TestUnionQueries(t *testing.T) {
	opts := provider.Options{Driver: provider.PostgresSQLDatabaseProviderName}
	q1 := NewQueryBuilder(opts).
		Select("users", "id", "email").
		Where("role = ?", "admin")

	s1, a1 := q1.Build()

	q2 := NewQueryBuilder(opts).
		Select("users", "id", "email").
		Where("role = ?", "manager")

	s2, a2 := q2.Build()

	sql := fmt.Sprintf("%s UNION %s", s1, s2)
	args := append(a1, a2...)

	expectedSQL := "SELECT id, email FROM users WHERE role = $1 UNION SELECT id, email FROM users WHERE role = $2"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
	if len(args) != 2 || args[0] != "admin" || args[1] != "manager" {
		t.Errorf("Expected args to be [\"admin\", \"manager\"], got %v", args)
	}
}

func TestCaseWhenClause(t *testing.T) {
	opts := provider.Options{Driver: provider.PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opts).
		Select("orders", "id", "amount", "CASE WHEN amount > 100 THEN 'high' ELSE 'low' END AS category")

	sql, _ := q.Build()
	expectedSQL := "SELECT id, amount, CASE WHEN amount > 100 THEN 'high' ELSE 'low' END AS category FROM orders"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
}

func TestWindowFunction(t *testing.T) {
	opts := provider.Options{Driver: provider.PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opts).
		Select("orders", "id", "amount", "RANK() OVER (PARTITION BY customer_id ORDER BY amount DESC) AS rank")

	sql, _ := q.Build()
	expectedSQL := "SELECT id, amount, RANK() OVER (PARTITION BY customer_id ORDER BY amount DESC) AS rank FROM orders"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
}
