# Project Analysis and Possible Improvements

This document summarizes observations about the current codebase and lists actionable improvements. It also points out potential pitfalls discovered during a quick review.

## Overview

This repository provides a unified data provider abstraction over multiple relational databases using sqlx (on top of database/sql). It includes:
- Provider interface (connectivity, initialization, migration hooks, status, SQL builder)
- Implementations for SQLite (including in-memory), MySQL, PostgreSQL, Oracle (behind build tags)
- Options struct with functional options
- A simple SQL builder

## Strengths
- Clear Provider interface abstraction
- Uses sqlx for ergonomics
- Functional options for configuration
- Build tags to keep drivers optional
- Basic pooling knobs (SetMaxOpenConns, SetMaxIdleConns) already present for MySQL/Postgres

## Notable Issues / Bugs
- postgres.go: PGSQLProvider struct lacks an `options *Options` field but methods reference `p.options`. The constructor also tries to set `options: options` in the returned struct which is not defined. This will not compile under postgres build tag.
- mysql.go: `context.WithTimeout(m.Context, 5)` uses a `time.Duration` of 5 nanoseconds, likely intended as seconds. Should be `5 * time.Second`.
- postgres.go: same timeout issue as MySQL implementation.
- Migration-related methods are `panic("implement me")` and therefore not safe to call.
- Global mutable `opts` in options.go and `driverName` in dataprovider.go introduce hidden state and potential race conditions in concurrent or multi-provider setups.
- DataSourceName assembly is simplistic, missing various connection options (TLS, timeouts, parsing time for MySQL, etc.).

## Design/Architecture Improvements
1. Eliminate global mutable state: 
   - Remove global `opts` singleton; NewOptions should return a fresh struct initialized with defaults. ✓ (proposal)
   - Avoid package-global `driverName`; keep driver in the Provider instance (Status can read from provider’s options).
2. Provider lifecycle & context:
   - Store `options *Options` consistently in all provider structs.
   - Make `Context` in Options optional and default to context.Background(). Ensure all timeouts use explicit durations.
   - Use context and timeouts with clear units.
3. Connection configuration:
   - Expose pooling settings beyond size: MaxIdleConns, ConnMaxLifetime, ConnMaxIdleTime.
   - Support DSN flags (e.g., `parseTime=true` for MySQL) and SSL/TLS options.
   - Allow setting connection timeout and read/write timeouts via Options per driver.
4. Error handling & health checks:
   - Replace `panic("implement me")` with proper error returns.
   - Add healthcheck helpers (Ping with timeout) and maybe a background checker to update status.
5. Migrations:
   - Provide a simple migration runner or integrate with a known tool (golang-migrate) behind an interface.
6. SQL Builder:
   - Consider separating SQL builder into a dedicated package; add tests for edge cases (quoting, injection safety, dialect quoting differences).
7. Testing & CI:
   - Add more unit tests for providers and SQL builder.
   - Add integration tests using Docker for MySQL/Postgres (via GitHub Actions services) gated by build tags.
8. Documentation:
   - Expand README with clear examples per driver, options listing, and pitfalls.
   - Document build tags and how to enable drivers.
9. Observability:
   - Add optional logging hooks for queries (with sampling) and metrics (Prometheus) to measure pool saturation, latency, errors.
10. Extensibility:
   - Consider an interface to plug in read/write splitting at the Provider level (master/replica aware), even though the current issue requests GORM-based guidance (see GORM_READ_WRITE.md).

## Quick Fixes (Low Effort, High Value)
- Fix postgres.go struct to include `options *Options` and timeout units. 
- Fix mysql.go timeout units. 
- Add ConnMaxLifetime/IdleTime configuration to all providers. 
- Make NewOptions return a new Options value instead of mutating a package-global.

These fixes would improve reliability and make the package safer in production contexts.
