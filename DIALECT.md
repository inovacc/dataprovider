# Diferencias de Sintaxis SQL entre Motores Populares

Este documento resume las principales diferencias de sintaxis SQL entre los motores de base de datos más usados: PostgreSQL, MySQL/MariaDB, SQLite y SQL Server.

---

## ⚙️ 1. Tipos de Datos

| Tipo              | PostgreSQL      | MySQL             | SQLite           | SQL Server       |
|------------------|------------------|------------------|------------------|------------------|
| Entero pequeño   | `SMALLINT`       | `SMALLINT`       | `INTEGER`        | `SMALLINT`       |
| Entero normal    | `INTEGER`        | `INT`            | `INTEGER`        | `INT`            |
| Cadena corta     | `VARCHAR(n)`     | `VARCHAR(n)`     | `TEXT`/`VARCHAR` | `VARCHAR(n)`     |
| Texto largo      | `TEXT`           | `TEXT`           | `TEXT`           | `VARCHAR(MAX)`   |
| Booleano         | `BOOLEAN`        | `TINYINT(1)`     | `INTEGER`        | `BIT`            |
| Auto-incremental | `SERIAL`         | `AUTO_INCREMENT` | `AUTOINCREMENT`  | `IDENTITY`       |
| JSON             | `JSON` / `JSONB` | `JSON` (5.7+)    | No nativo        | `NVARCHAR(MAX)`  |

---

## 🔑 2. Claves Primarias y Auto Incremento

| Sintaxis                     | PostgreSQL                 | MySQL                  | SQLite                  | SQL Server              |
|-----------------------------|----------------------------|------------------------|-------------------------|--------------------------|
| Auto-increment              | `id SERIAL PRIMARY KEY`    | `id INT AUTO_INCREMENT PRIMARY KEY` | `id INTEGER PRIMARY KEY AUTOINCREMENT` | `id INT IDENTITY(1,1) PRIMARY KEY` |
| UUID como PK                | `uuid UUID DEFAULT gen_random_uuid()` | `CHAR(36)` y lo llenás tú | `TEXT` manual           | `UNIQUEIDENTIFIER` con `NEWID()` |

---

## 🔍 3. Búsqueda y Operadores

| Operación       | PostgreSQL            | MySQL                   | SQLite             | SQL Server             |
|----------------|------------------------|-------------------------|--------------------|-------------------------|
| ILIKE (case-insensitive) | `ILIKE`        | `LIKE` con `COLLATE`    | `LIKE` es case-insensitive por defecto | `LIKE` con `COLLATE` |
| Concatenar     | `'||'`                | `'CONCAT(a, b)'`        | `'||'`             | `'+'` o `CONCAT()`     |
| Regex          | `~` o `SIMILAR TO`    | `REGEXP`                | No estándar        | `LIKE` limitado o CLR  |

---

## 📅 4. Fechas y Horas

| Función                  | PostgreSQL       | MySQL              | SQLite           | SQL Server        |
|--------------------------|------------------|--------------------|------------------|-------------------|
| Fecha actual             | `CURRENT_DATE`   | `CURDATE()`        | `date('now')`    | `GETDATE()`       |
| Fecha y hora actual      | `NOW()`          | `NOW()`            | `datetime('now')`| `GETDATE()`       |
| Extraer año              | `EXTRACT(YEAR FROM date)` | `YEAR(date)` | `strftime('%Y', date)` | `YEAR(date)`    |

---

## 🧱 5. DDL – Crear Tablas, Índices, Constraints

| Característica          | PostgreSQL                 | MySQL                       | SQLite             | SQL Server              |
|------------------------|-----------------------------|-----------------------------|--------------------|--------------------------|
| CHECK constraint       | ✅                          | ✅ (a veces ignorado)       | ✅                 | ✅                       |
| Foreign key            | ✅                          | ✅                          | ✅ (limitado)      | ✅                       |
| Enum                   | `CREATE TYPE ... AS ENUM`   | `ENUM(...)`                 | No nativo          | No nativo (usa `CHECK`) |
| Índice parcial         | ✅ `WHERE condition`        | No                          | No                 | Parcialmente con `FILTER` |

---

## 📦 6. JSON

| Operación           | PostgreSQL             | MySQL                 | SQLite               | SQL Server             |
|---------------------|------------------------|------------------------|-----------------------|-------------------------|
| Acceso por clave    | `json_col->>'key'`     | `JSON_EXTRACT(...)`   | No nativo             | `JSON_VALUE(...)`       |
| Indexación JSON     | ✅ (`jsonb`)           | ✅ (`GENERATED`)      | ❌                    | ✅ (`computed column`)  |
| Validación JSON     | ✅                     | ✅                    | ❌                    | ✅                      |

---

## 🧠 7. CTEs, Funciones de Ventana

| Característica     | PostgreSQL | MySQL       | SQLite       | SQL Server |
|--------------------|------------|-------------|--------------|------------|
| CTE (`WITH`)       | ✅         | ✅ (8.0+)    | ✅           | ✅         |
| Window Functions   | ✅         | ✅ (8.0+)    | ✅ (3.25+)    | ✅         |
| Lateral Joins      | ✅         | ❌           | ❌           | ❌         |

---

## 🧪 8. Funcionalidades Especiales

| Función / característica | PostgreSQL         | MySQL             | SQLite           | SQL Server        |
|--------------------------|--------------------|-------------------|------------------|-------------------|
| UPSERT (`INSERT ... ON CONFLICT`) | `ON CONFLICT DO UPDATE` | `ON DUPLICATE KEY UPDATE` | `INSERT OR REPLACE` | `MERGE` |
| Array                     | ✅ (`int[]`)       | ❌ (usa JSON)     | ❌               | ❌                |
| Materialized Views       | ✅                 | ❌                | ❌               | ✅                |

---

## 🎯 Conclusión

| Si necesitas...                                | Usa...               |
|------------------------------------------------|----------------------|
| Sintaxis rica, potente y estricta              | PostgreSQL           |
| Compatibilidad con muchos hosts y ORMs         | MySQL/MariaDB        |
| Ligereza y portabilidad (no necesita servidor) | SQLite               |
| Integración con herramientas Microsoft         | SQL Server           |