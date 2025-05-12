# SQLBuilder Testing Suite

## Overview

This testing suite validates a SQL query builder with dialect-specific support and composable query generation. The builder targets cloud-native applications requiring dynamic, multi-database SQL composition with high reliability and test coverage.

---

## Supported SQL Features

### 🗄️ DDL (Data Definition Language)

* `CREATE TABLE`
* `DROP TABLE`

### 🔄 DML (Data Manipulation Language)

* `SELECT`
* `INSERT INTO` (including multi-row)
* `UPDATE`
* `DELETE`
* `MERGE` / `UPSERT`

### 🔎 Query Enhancements

* `JOIN`, `LEFT JOIN`, `RIGHT JOIN`
* `GROUP BY`, `HAVING`, `ORDER BY`, `LIMIT`, `OFFSET`
* `ALIAS` with `AS`
* `CASE WHEN`, `RANK()`, `OVER()`
* `WITH` (CTE), `EXISTS`, nested queries
* `UNION`
* Raw SQL injection (`Raw()`)

### ⚙️ Control and Extensibility

* Placeholder substitution by dialect (e.g., `$1` for PostgreSQL, `:p1` for Oracle)
* Transactional queries (`BEGIN`, `COMMIT`, `ROLLBACK`)
* Dynamic argument binding

---

## Database Dialect Support

* PostgreSQL ✅
* Oracle ✅
* MySQL ✅
* SQLite ✅
* MariaDB (inherits MySQL behavior) ✅

---

## Test Coverage Summary

| Test                       | Purpose                                  |
| -------------------------- | ---------------------------------------- |
| `TestQueryBuilderVariants` | Validates dialect placeholder formatting |
| `TestCreateDropDeleteSQL`  | DDL and DELETE queries with WHERE clause |
| `TestInsertAndUpdateSQL`   | Parameterized inserts and updates        |
| `TestMergeBuilder`         | `MERGE INTO`/`UPSERT` logic              |
| `TestRawSQLInjection`      | Arbitrary SQL with dynamic args          |
| `TestAliasSupport`         | `AS` alias handling                      |
| `TestGroupByHaving`        | Aggregation filtering                    |
| `TestJoinClauses`          | Join clause variations                   |
| `TestJoinGroupByHaving`    | Combined grouping and join               |
| `TestNestedSelect`         | Subqueries in `WHERE IN`                 |
| `TestUnionQueries`         | `UNION` support across queries           |
| `TestCaseWhenClause`       | Conditional select logic                 |
| `TestWithCTE`              | CTE-based query composition              |
| `TestExistsClause`         | `EXISTS` clause validation               |
| `TestMultiRowInsert`       | Multi-row insert structure               |
| `TestTransactionalQuery`   | Transaction begin/commit handling        |
| `TestWindowFunction`       | Ranking via window functions             |

---

## Run Instructions

To run tests:

```bash
go test ./...
```

---

## Authors & License

© 2025 - Built and maintained by distributed systems engineers. Released under MIT License.

---

Need help extending to NoSQL or GraphQL DSLs? Reach out anytime.
