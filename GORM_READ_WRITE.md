# Production-Ready Read/Write Splitting in GORM (1 master, N replicas)

This guide shows how to implement read/write splitting with GORM using `gorm.io/plugin/dbresolver`. It covers:
- Topology: 1 master (write) + N replicas (read)
- Connection pooling
- Basic failover handling strategies
- Health checks and operational tips

Works with MySQL and PostgreSQL (dialector examples included).

## Dependencies

```bash
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get gorm.io/driver/postgres
go get gorm.io/plugin/dbresolver
```

## Minimal Example (MySQL)

```go
package main

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func main() {
	masterDSN := "user:pass@tcp(master.db:3306)/app?parseTime=true&charset=utf8mb4,utf8&loc=UTC&timeout=5s&readTimeout=5s&writeTimeout=5s"
	replica1DSN := "user:pass@tcp(replica-1.db:3306)/app?parseTime=true&charset=utf8mb4,utf8&loc=UTC&timeout=5s&readTimeout=5s&writeTimeout=5s"
	replica2DSN := "user:pass@tcp(replica-2.db:3306)/app?parseTime=true&charset=utf8mb4,utf8&loc=UTC&timeout=5s&readTimeout=5s&writeTimeout=5s"

	db, err := gorm.Open(mysql.Open(masterDSN), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	resolver := dbresolver.Register(dbresolver.Config{
		// Source = master(s), Replicas = replica(s)
		Sources:  []gorm.Dialector{mysql.Open(masterDSN)},
		Replicas: []gorm.Dialector{mysql.Open(replica1DSN), mysql.Open(replica2DSN)},
		Policy:   dbresolver.RandomPolicy{}, // random among replicas for reads
	})

	if err := db.Use(resolver);
	err != nil {
		log.Fatal(err)
	}

	// Pooling (defaults can be too low/high; tune per your workload)
	// Apply to all connections registered in resolver (master+replicas)
	resolver.SetMaxOpenConns(50)
	resolver.SetMaxIdleConns(25)
	resolver.SetConnMaxLifetime(30 * time.Minute)
	resolver.SetConnMaxIdleTime(5 * time.Minute)

	// Default routing: SELECT -> replicas, INSERT/UPDATE/DELETE -> master.
	// In transactions, all statements go to the same DB where the tx began.

	// Example read (goes to a replica):
	var count int64
	if err := db.Table("users").Count(&count).Error; err != nil {
		log.Fatal(err)
	}

	// Example write (goes to master):
	if err := db.Exec("INSERT INTO users(name) VALUES (?)", "alice").Error; err != nil {
		log.Fatal(err)
	}
}
```

## PostgreSQL Dialectors

```go
master := postgres.Open("host=master.db user=app password=secret dbname=app port=5432 sslmode=disable connect_timeout=5")
rep1  := postgres.Open("host=replica-1.db user=app password=secret dbname=app port=5432 sslmode=disable connect_timeout=5")
rep2  := postgres.Open("host=replica-2.db user=app password=secret dbname=app port=5432 sslmode=disable connect_timeout=5")

resolver := dbresolver.Register(dbresolver.Config{
    Sources:  []gorm.Dialector{master},
    Replicas: []gorm.Dialector{rep1, rep2},
    Policy:   dbresolver.RandomPolicy{},
})
```

## Routing Controls

- Reads by default: SELECT, Scan, Find, First, etc.
- Writes by default: Create, Save, Update, Delete, Exec, Raw that changes data.
- Transaction scope: Once a transaction begins on a connection, all queries within it will stick to the same DB. GORM typically uses the write source for transactions, which is desirable to avoid replication lag issues.
- Manual override: you can force read or write with clauses if needed:

```go
import "gorm.io/plugin/dbresolver"

// Force read on master (e.g., require freshest data immediately after write)
db.Clauses(dbresolver.Write).First(&user, id)

// Force read on replicas explicitly
db.Clauses(dbresolver.Read).Find(&users)
```

## Connection Pooling Best Practices

- MaxOpenConns: Set based on DB capacity and app concurrency. Monitor for saturation.
- MaxIdleConns: Keep a warm pool to reduce connection churn.
- ConnMaxLifetime: Rotate connections under load to avoid server-side limits and mitigate memory leaks.
- ConnMaxIdleTime: Trim idle connections over time.
- Apply settings through the resolver so they apply to every registered pool (sources and replicas).

## Failover Handling

GORM dbresolver provides load balancing across replicas and basic error handling. However, write failover (master outage) requires decisions:

1) Infrastructure-level failover (recommended):
   - Use a managed DB (RDS/Aurora, Cloud SQL) or orchestrated cluster (Patroni, Orchestrator, MHA) with a single writable endpoint (VIP/DNS, ProxySQL, PgBouncer/HAProxy).
   - Your app continues to use the same DSN; the endpoint is moved to the new primary automatically.

2) Application-level awareness:
   - Register multiple sources (if your topology allows multiple potential primaries, but normally only one is writable at a time).
   - Implement a health checker that pings master and replicas periodically; on prolonged failures, either:
     - Rebuild the resolver with a promoted replica DSN, or
     - Switch routing temporarily (risky) until promotion completes.

3) Custom selection policy:
   - Implement `dbresolver.Policy` to choose among replicas based on health/latency. Example skeleton:

```go
type HealthAwarePolicy struct { /* your fields */ }

func (p HealthAwarePolicy) Resolve(connPools []gorm.ConnPool) gorm.ConnPool {
    // pick only healthy pools; fallback to any available
    // you can maintain a health table updated by a separate goroutine
    return connPools[rand.Intn(len(connPools))]
}

// Usage:
resolver := dbresolver.Register(dbresolver.Config{
    Sources:  []gorm.Dialector{master},
    Replicas: []gorm.Dialector{rep1, rep2},
    Policy:   HealthAwarePolicy{},
})
```

4) Graceful replica removal:
   - On health failure, rebuild the resolver or maintain DSN lists and re-register plugin to drop unhealthy nodes.

## Health Checks

- Periodically ping each pool using `sqlDB := db.ConnPool.(interface{ DB() *sql.DB })` pattern, or maintain separate sql.DB ping logic outside GORM using the same DSNs.
- Record latency and error rates; if a node is degraded, remove it from replicas and re-register the resolver.
- For MySQL, guard against replication lag: for read-after-write, either route to master using `db.Clauses(dbresolver.Write)` or read from a replica with GTID/semisync awareness via an external proxy.

## Observability

- Log slow queries (GORM logger with thresholds).
- Export DB pool metrics (pprof, Prometheus) including open/idle connections, wait counts, wait durations, max lifetime closes.
- Correlate application errors with database node health.

## Putting It Together: Configurable Bootstrap

```go
type DBConfig struct {
    MasterDSN  string
    ReplicaDSN []string

    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}

func NewGormRW(cfg DBConfig) (*gorm.DB, error) {
    db, err := gorm.Open(mysql.Open(cfg.MasterDSN), &gorm.Config{})
    if err != nil { return nil, err }

    replicas := make([]gorm.Dialector, 0, len(cfg.ReplicaDSN))
    for _, d := range cfg.ReplicaDSN { replicas = append(replicas, mysql.Open(d)) }

    resolver := dbresolver.Register(dbresolver.Config{
        Sources:  []gorm.Dialector{mysql.Open(cfg.MasterDSN)},
        Replicas: replicas,
        Policy:   dbresolver.RandomPolicy{},
    })

    if err := db.Use(resolver); err != nil { return nil, err }

    if cfg.MaxOpenConns > 0 { resolver.SetMaxOpenConns(cfg.MaxOpenConns) }
    if cfg.MaxIdleConns > 0 { resolver.SetMaxIdleConns(cfg.MaxIdleConns) }
    if cfg.ConnMaxLifetime > 0 { resolver.SetConnMaxLifetime(cfg.ConnMaxLifetime) }
    if cfg.ConnMaxIdleTime > 0 { resolver.SetConnMaxIdleTime(cfg.ConnMaxIdleTime) }

    return db, nil
}
```

## Integration Notes for This Repository

This project currently uses sqlx, not GORM. If you decide to migrate or provide an optional GORM-based Provider:
- Create a new provider variant that internally exposes a `*gorm.DB` instead of `*sqlx.DB`.
- Use the approach in this guide to set up read/write splitting via dbresolver.
- Keep the existing sqlx Provider for backwards compatibility; consider feature flagging which backend (sqlx vs gorm) to use.

For more details on project-wide improvements, see `IMPROVEMENTS.md`.
