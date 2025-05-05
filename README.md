[![CI and Test](https://github.com/inovacc/dataprovider/actions/workflows/ci.yml/badge.svg)](https://github.com/inovacc/dataprovider/actions/workflows/ci.yml)
[![Generate SBOM](https://github.com/inovacc/dataprovider/actions/workflows/sbom.yml/badge.svg)](https://github.com/inovacc/dataprovider/actions/workflows/sbom.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/inovacc/dataprovider)](https://goreportcard.com/report/github.com/inovacc/dataprovider)
[![Go Reference](https://pkg.go.dev/badge/github.com/inovacc/dataprovider.svg)](https://pkg.go.dev/github.com/inovacc/dataprovider)

# Dataprovider

dataprovider is a module that provides a uniform interface to interact behind the scene with many databases and
configurations.
On this design is used the [jmoiron/sqlx](https://github.com/jmoiron/sqlx) package that provides a set of extensions on
top of
the excellent built-in [database/sql](https://pkg.go.dev/database/sql) package.

The dataprovider in `memory` and `sqlite` uses `modernc.org/sqlite` package that is a `CGO` free SQLite3 driver for Go.

## Working on

- [ ] Database Migration
- [x] Database Connection
- [x] Database Transaction
- [x] Database Query
- [x] Database QueryRow
- [x] Database QueryRowx

## Supported databases

- [x] Memory
- [x] SQLite3
- [x] MySQL
- [x] PostgreSQL
- [ ] SQL Server
- [x] Oracle
- [ ] CockroachDB
- [ ] ClickHouse
- [ ] Cassandra
- [ ] MongoDB
- [ ] Redis
- [ ] InfluxDB
- [ ] Elasticsearch
- [ ] BigQuery
- [ ] Google Cloud Firestore
- [ ] Google Cloud Spanner
- [ ] Google Cloud Datastore

## How to use

```shell
go get github.com/inovacc/dataprovider
```

Need to build with the tag `mysql`, `postgres`, or `oracle` to use the specific database. Default driver is `sqlite` in
`memory` mode all data is lost when the program ends.

```shell
go build -tags mysql
```

```shell
go build -tags postgres
```

```shell
go build -tags oracle
```

## Example of initialization

```go
package main

import "github.com/inovacc/dataprovider"

func main() {
	// Create a config with driver name to initialize the data provider
	var provider = dataprovider.Must(dataprovider.NewDataProvider(
		dataprovider.NewOptions(dataprovider.WithSqliteDB("test", ".")),
	))

	// Initialize the database
	query := "CREATE TABLE IF NOT EXISTS ...;"
	if err := provider.InitializeDatabase(query); err != nil {
		panic(err)
	}

	// Get the connection and use it as sqlx.DB or sql.DB
	conn := provider.GetConnection()

	query := provider.NewSQLBuilder(provider).
		Table("users").
		Select("id", "name", "email").
		Where("age > 18").
		OrderBy("name ASC").
		Limit(10).
		Offset(5).
		Build()
}

```

## Example of usage

```go
package main

import (
	"encoding/json"
	"github.com/inovacc/dataprovider"
	"os"
)

type User struct {
	Id        int64  `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Gender    string `json:"gender" db:"gender"`
	IpAddress string `json:"ip_address" db:"ip_address"`
	City      string `json:"city" db:"city"`
}

func main() {
	// Create a config with driver name to initialize the data provider
	opts := dataprovider.NewOptions(
		dataprovider.WithDriver(dataprovider.MemoryDataProviderName),
		dataprovider.WithConnectionString("file:test.sqlite3?cache=shared"),
	)

	var provider = dataprovider.Must(dataprovider.NewDataProvider(opts))

	// Initialize the database
	query := "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, first_name TEXT, last_name TEXT, email TEXT, gender TEXT, ip_address TEXT, city TEXT);"
	if err := provider.InitializeDatabase(query); err != nil {
		panic(err)
	}

	// Get the connection
	conn := provider.GetConnection()

	// Begin a transaction
	tx := conn.MustBegin()

	// Insert data
	tx.MustExec("insert into users (first_name, last_name, email, gender, ip_address, city) values ('Marcus', 'Bengefield', 'mbengefield0@vistaprint.com', 'Male', '83.121.11.105', 'Miura');")
	tx.MustExec("insert into users (first_name, last_name, email, gender, ip_address, city) values ('Brandise', 'Mateuszczyk', 'bmateuszczyk1@vistaprint.com', 'Female', '131.187.209.233', 'Dalududalu');")
	tx.MustExec("insert into users (first_name, last_name, email, gender, ip_address, city) values ('Ray', 'Ginnaly', 'rginnaly2@merriam-webster.com', 'Male', '76.71.94.89', 'Al Baqāliţah');")

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		panic(err)
	}

	// Query the data
	var users []User

	query = "SELECT * FROM users"
	rows, err := conn.Queryx(query)
	if err != nil {
		panic(err)
	}

	// Scan the data
	for rows.Next() {
		var user User
		if err = rows.StructScan(&user); err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	// Print the data
	if err = json.NewEncoder(os.Stdout).Encode(users); err != nil {
		panic(err)
	}
}
```
