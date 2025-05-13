package dataprovider

import (
	"fmt"
	"strings"
	"testing"
)

func TestQuery(t *testing.T) {
	query := NewQueryBuilder(&Options{Driver: SQLiteDataProviderName}).
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
	query := NewQueryBuilder(&Options{Driver: SQLiteDataProviderName}).
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
	query, _ := NewQueryBuilder(&Options{Driver: SQLiteDataProviderName}).
		CreateTable("users").
		IfNotExists().
		Columns(
			Column("id").Type("SERIAL").PrimaryKey(),
			Column("first_name").Type("TEXT"),
			Column("last_name").Type("TEXT"),
			Column("email").Type("TEXT"),
			Column("ip_address").Type("TEXT"),
			Column("city").Type("TEXT"),
		).
		Build()

	t.Log("CREATE TABLE Query:")
	t.Log(query)

	if query != "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, first_name TEXT, last_name TEXT, email TEXT, ip_address TEXT, city TEXT)" {
		t.Error("CREATE TABLE Query is not correct")
	}
}

func TestQueryUpdate(t *testing.T) {
	query := NewQueryBuilder(&Options{Driver: SQLiteDataProviderName}).
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

func TestQueryBuilderVariants(t *testing.T) {
	tests := []struct {
		driver      databaseKind
		expectedSQL string
	}{
		{
			driver:      OracleDatabaseProviderName,
			expectedSQL: "SELECT id, name FROM users WHERE status = :p1 ORDER BY name LIMIT 10",
		},
		{
			driver:      PostgresSQLDatabaseProviderName,
			expectedSQL: "SELECT id, name FROM users WHERE status = $1 ORDER BY name LIMIT 10",
		},
		{
			driver:      MySQLDatabaseProviderName,
			expectedSQL: "SELECT id, name FROM users WHERE status = ? ORDER BY name LIMIT 10",
		},
		{
			driver:      SQLiteDataProviderName,
			expectedSQL: "SELECT id, name FROM users WHERE status = ? ORDER BY name LIMIT 10",
		},
		{
			driver:      "mariadb", // assuming MariaDB shares behavior with MySQL
			expectedSQL: "SELECT id, name FROM users WHERE status = ? ORDER BY name LIMIT 10",
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.driver), func(t *testing.T) {
			opt := &Options{Driver: tt.driver}
			qb := NewQueryBuilder(opt)
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
		driver       databaseKind
		operation    string
		builderFunc  func(SQLBuilder) (string, []any)
		expectedSQLs map[string]string
	}{
		{
			driver:    OracleDatabaseProviderName,
			operation: "create",
			builderFunc: func(b SQLBuilder) (string, []any) {
				return b.CreateTable("users", "id INT, name VARCHAR2(50)").Build()
			},
			expectedSQLs: map[string]string{
				"oracle": "CREATE TABLE users (id INT, name VARCHAR2(50))",
			},
		},
		{
			driver:    PostgresSQLDatabaseProviderName,
			operation: "drop",
			builderFunc: func(b SQLBuilder) (string, []any) {
				return b.DropTable("users").Build()
			},
			expectedSQLs: map[string]string{
				"postgres": "DROP TABLE users",
			},
		},
		{
			driver:    MySQLDatabaseProviderName,
			operation: "delete",
			builderFunc: func(b SQLBuilder) (string, []any) {
				return b.DeleteFrom("users").Where("name = ?", "john").Build()
			},
			expectedSQLs: map[string]string{
				"mysql": "DELETE FROM users WHERE name = ?",
			},
		},
		{
			driver:    SQLiteDataProviderName,
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
			opt := &Options{Driver: tt.driver}
			builder := NewQueryBuilder(opt)
			sql, _ := tt.builderFunc(builder)

			key := tt.driver
			if val, ok := tt.expectedSQLs[string(key)]; ok && sql != val {
				t.Errorf("expected SQL %q, got %q", val, sql)
			}
		})
	}
}

func TestInsertAndUpdateSQL(t *testing.T) {
	tests := []struct {
		driver       databaseKind
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
			opt := &Options{Driver: tt.driver}
			builder := NewQueryBuilder(opt)
			sql, _ := tt.builderFunc(builder)

			expected := tt.expectedSQLs[string(tt.driver)]
			if sql != expected {
				t.Errorf("expected SQL %q, got %q", expected, sql)
			}
		})
	}
}

func TestMergeBuilder(t *testing.T) {
	tests := []struct {
		driver   databaseKind
		expected string
	}{
		{
			driver:   PostgresSQLDatabaseProviderName,
			expected: "MERGE INTO users ON id = $1 WHEN MATCHED THEN UPDATE SET email = $2 WHEN NOT MATCHED THEN INSERT (id, email) VALUES ($3, $4)",
		},
		{
			driver:   OracleDatabaseProviderName,
			expected: "MERGE INTO users ON id = :p1 WHEN MATCHED THEN UPDATE SET email = :p2 WHEN NOT MATCHED THEN INSERT (id, email) VALUES (:p3, :p4)",
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.driver), func(t *testing.T) {
			opt := &Options{Driver: tt.driver}
			q := NewQueryBuilder(opt).
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
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opt).
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
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opt).
		Select("users", "id", "email").
		As("u")

	sql, _ := q.Build()
	expectedSQL := "SELECT id, email FROM users AS u"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
}

func TestGroupByHaving(t *testing.T) {
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opt).
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
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opt).
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
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opt).
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
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	inner := NewQueryBuilder(opt).
		Select("payments", "user_id").
		Where("status = ?", "completed")

	subSQL, subArgs := inner.Build()

	q := NewQueryBuilder(opt).
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
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	q1 := NewQueryBuilder(opt).
		Select("users", "id", "email").
		Where("role = ?", "admin")

	q2 := NewQueryBuilder(opt).
		Select("users", "id", "email").
		Where("role = ?", "manager")

	query := q1.Union(q2)
	sql, args := query.Build()

	expectedSQL := "SELECT id, email FROM users WHERE role = $1 UNION SELECT id, email FROM users WHERE role = $2"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
	if len(args) != 2 || args[0] != "admin" || args[1] != "manager" {
		t.Errorf("Expected args to be [\"admin\", \"manager\"], got %v", args)
	}
}

func TestCaseWhenClause(t *testing.T) {
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opt).
		Select("orders", "id", "amount", "CASE WHEN amount > 100 THEN 'high' ELSE 'low' END AS category")

	sql, _ := q.Build()
	expectedSQL := "SELECT id, amount, CASE WHEN amount > 100 THEN 'high' ELSE 'low' END AS category FROM orders"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
}

func TestWithCTE(t *testing.T) {
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	cte := NewQueryBuilder(opt).
		Select("payments", "user_id", "SUM(amount) AS total").
		GroupBy("user_id")

	cteSQL, cteArgs := cte.Build()
	main := NewQueryBuilder(opt).
		Select("summary", "user_id", "total").
		Raw(fmt.Sprintf("WITH summary AS (%s) SELECT user_id, total FROM summary WHERE total > ?", cteSQL), append(cteArgs, 1000)...) // injects entire CTE with final condition

	sql, args := main.Build()
	expectedSQL := "WITH summary AS (SELECT user_id, SUM(amount) AS total FROM payments GROUP BY user_id) SELECT user_id, total FROM summary WHERE total > $1"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
	if len(args) != 1 || args[0] != 1000 {
		t.Errorf("Expected args to be [1000], got %v", args)
	}
}

func TestExistsClause(t *testing.T) {
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	sub := NewQueryBuilder(opt).
		Select("orders", "1").Where("orders.user_id = users.id")
	subSQL, subArgs := sub.Build()

	q := NewQueryBuilder(opt).
		Select("users", "id", "email").
		Where(fmt.Sprintf("EXISTS (%s)", subSQL), subArgs...)

	sql, args := q.Build()
	expectedSQL := "SELECT id, email FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
	if len(args) != 0 {
		t.Errorf("Expected no args, got %v", args)
	}
}

func TestMultiRowInsert(t *testing.T) {
	opt := &Options{Driver: MySQLDatabaseProviderName}
	builder := NewQueryBuilder(opt)
	values := []any{"john@example.com", "doe@example.com"}
	placeholders := []string{"(?)", "(?)"} // simulated

	builder.Raw(fmt.Sprintf("INSERT INTO users (email) VALUES %s", strings.Join(placeholders, ", ")), values...)
	sql, args := builder.Build()
	expectedSQL := "INSERT INTO users (email) VALUES (?), (?)"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
	if len(args) != 2 || args[0] != "john@example.com" || args[1] != "doe@example.com" {
		t.Errorf("Expected args to be correct, got %v", args)
	}
}

func TestTransactionalQuery(t *testing.T) {
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	start := NewQueryBuilder(opt).Raw("BEGIN")
	commit := NewQueryBuilder(opt).Raw("COMMIT")

	s1, _ := start.Build()
	s2, _ := commit.Build()

	if s1 != "BEGIN" {
		t.Errorf("Expected BEGIN, got %q", s1)
	}
	if s2 != "COMMIT" {
		t.Errorf("Expected COMMIT, got %q", s2)
	}
}

func TestWindowFunction(t *testing.T) {
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	q := NewQueryBuilder(opt).
		Select("orders", "id", "amount", "RANK() OVER (PARTITION BY customer_id ORDER BY amount DESC) AS rank")

	sql, _ := q.Build()
	expectedSQL := "SELECT id, amount, RANK() OVER (PARTITION BY customer_id ORDER BY amount DESC) AS rank FROM orders"

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %q\nGot: %q", expectedSQL, sql)
	}
}

func TestExportAsJSON(t *testing.T) {
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	builder := NewQueryBuilder(opt).
		Select("users", "id", "email").
		Where("email = ?", "test@example.com")

	jsonOut, err := builder.ExportAsJSON()
	if err != nil {
		t.Fatalf("ExportAsJSON failed: %v", err)
	}

	expedted := "{\"kind\":\"select\",\"columns\":[\"id\",\"email\"],\"from\":\"users\",\"where\":\"email = $1\",\"sql\":\"SELECT id, email FROM users WHERE email = $1\",\"args\":[\"test@example.com\"]}"

	if jsonOut != expedted {
		t.Errorf("Expected JSON output: %q\nGot: %q", expedted, jsonOut)
	}
}

func TestExportAsXML(t *testing.T) {
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	builder := NewQueryBuilder(opt).
		Select("users", "id", "email").
		Where("email = ?", "john@example.com")

	xmlOut, err := builder.ExportAsXML()
	if err != nil {
		t.Fatalf("ExportAsXML failed: %v", err)
	}

	expected := "<query><kind>select</kind><columns>id</columns><columns>email</columns><from>users</from><where>email = $1</where><alias></alias><SQL>SELECT id, email FROM users WHERE email = $1</SQL><special></special><mergeTable></mergeTable><mergeOn></mergeOn><args>john@example.com</args></query>"

	if xmlOut != expected {
		t.Errorf("Expected XML output: %q\nGot: %q", expected, xmlOut)
	}
}

func TestExportAsYaml(t *testing.T) {
	opt := &Options{Driver: PostgresSQLDatabaseProviderName}
	builder := NewQueryBuilder(opt).
		Select("users", "id", "email").
		Where("email = ?", "john@example.com")

	yamlOut, err := builder.ExportAsYAML()
	if err != nil {
		t.Fatalf("ExportAsYAML failed: %v", err)
	}

	expected := `xmlname:
    space: ""
    local: query
kind: select
columns:
    - id
    - email
from: users
where: email = $1
groupby: []
having: []
orderby: []
limit: null
offset: null
alias: ""
joins: []
sql: SELECT id, email FROM users WHERE email = $1
special: ""
mergetable: ""
mergeon: ""
mergematchedset: []
mergeinsertcols: []
mergeinsertvals: []
rawclauses: []
queries: []
args:
    - john@example.com
data: null
`

	if yamlOut != expected {
		t.Errorf("Expected YAML output: %q\nGot: %q", expected, yamlOut)
	}
}
