# Database Migrations Guide

This project uses **golang-migrate** for database schema versioning and migrations, similar to Alembic (Python) or Prisma (TypeScript).

## 🎯 Why Migrations?

- **Version Control**: Track database schema changes in Git
- **Reproducibility**: Same schema across dev, staging, and production
- **Rollback**: Undo problematic migrations safely
- **Team Collaboration**: Merge schema changes without conflicts
- **Production Ready**: Industry standard for Go applications

## 📦 Installation

The migrate CLI tool is already installed via:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Binary location: `$env:USERPROFILE\go\bin\migrate.exe`

## 📁 Migration Files Structure

```
migrations/
├── 000001_create_users_table.up.sql    # Apply migration
├── 000001_create_users_table.down.sql  # Rollback migration
├── 000002_add_user_indexes.up.sql      # Next migration (example)
└── 000002_add_user_indexes.down.sql    # Rollback for it
```

## 🚀 Quick Start (Windows PowerShell)

### Using Our Helper Script (Recommended)

```powershell
# Show available commands
.\migrate.ps1

# Apply all pending migrations
.\migrate.ps1 up

# Rollback the last migration
.\migrate.ps1 down

# Create a new migration
.\migrate.ps1 create add_user_email_verification

# Check current version
.\migrate.ps1 version
```

### Using migrate CLI Directly

```powershell
$DB_URL = "postgres://devops:devops123@localhost:5432/devops_db?sslmode=disable"

# Apply migrations
& "$env:USERPROFILE\go\bin\migrate.exe" -path ./migrations -database $DB_URL up

# Rollback one migration
& "$env:USERPROFILE\go\bin\migrate.exe" -path ./migrations -database $DB_URL down 1

# Check version
& "$env:USERPROFILE\go\bin\migrate.exe" -path ./migrations -database $DB_URL version

# Force version (if stuck)
& "$env:USERPROFILE\go\bin\migrate.exe" -path ./migrations -database $DB_URL force 1
```

## 🛠️ Common Commands

### Apply Migrations

```powershell
# Apply all pending migrations
.\migrate.ps1 up
```

This runs all `.up.sql` files that haven't been applied yet.

### Rollback Migrations

```powershell
# Rollback the last migration
.\migrate.ps1 down

# Rollback specific number of migrations
& "$env:USERPROFILE\go\bin\migrate.exe" -path ./migrations -database $DB_URL down 2
```

### Create New Migration

```powershell
# Creates two files:
# - migrations/000002_add_profile_fields.up.sql
# - migrations/000002_add_profile_fields.down.sql
.\migrate.ps1 create add_profile_fields
```

Then edit the files:

**000002_add_profile_fields.up.sql:**

```sql
ALTER TABLE accounts_user
ADD COLUMN phone VARCHAR(20),
ADD COLUMN bio TEXT;
```

**000002_add_profile_fields.down.sql:**

```sql
ALTER TABLE accounts_user
DROP COLUMN IF EXISTS phone,
DROP COLUMN IF EXISTS bio;
```

### Check Migration Status

```powershell
# Show current migration version
.\migrate.ps1 version

# Example output: 1
```

### Force Migration Version (Recovery)

If a migration fails halfway and the database is in a "dirty" state:

```powershell
# Force it back to a known good version
.\migrate.ps1 force 1
```

## 🔄 Automatic Migrations on Startup

The application automatically runs pending migrations when it starts:

```go
// In database/db.go
func InitDB() {
    // ... connect to database
    RunMigrations() // Automatically applies pending migrations
}
```

This ensures the database is always up-to-date without manual intervention.

## 📝 Writing Good Migrations

### ✅ DO

- Keep migrations small and focused
- Always provide both `.up.sql` and `.down.sql`
- Test rollbacks before committing
- Add indexes for frequently queried columns
- Use `IF EXISTS` / `IF NOT EXISTS` where appropriate
- Add comments explaining complex changes

### ❌ DON'T

- Modify existing migration files after they're merged
- Delete data without backup strategy
- Make irreversible changes in production
- Skip testing migrations locally first

### Example Migration Templates

**Add Column:**

```sql
-- up
ALTER TABLE accounts_user ADD COLUMN avatar_url VARCHAR(500);

-- down
ALTER TABLE accounts_user DROP COLUMN IF EXISTS avatar_url;
```

**Create Index:**

```sql
-- up
CREATE INDEX idx_accounts_user_created_at ON accounts_user(created_at);

-- down
DROP INDEX IF EXISTS idx_accounts_user_created_at;
```

**Create Table:**

```sql
-- up
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES accounts_user(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- down
DROP TABLE IF EXISTS posts;
```

## 🚦 Migration Workflow

### Development

1. Create migration:

   ```powershell
   .\migrate.ps1 create add_new_feature
   ```

2. Edit the generated files in `migrations/`

3. Apply migration:

   ```powershell
   .\migrate.ps1 up
   ```

4. Test the application

5. Test rollback:

   ```powershell
   .\migrate.ps1 down
   .\migrate.ps1 up
   ```

6. Commit migration files to Git

### Production Deployment

1. Pull latest code with new migrations
2. Backup database
3. Apply migrations:
   ```bash
   ./migrate.ps1 up
   # or let the app auto-apply on startup
   ```
4. Verify application works
5. If issues: rollback migrations
   ```bash
   ./migrate.ps1 down
   ```

## 🐛 Troubleshooting

### "Dirty database version"

This happens when a migration failed halfway. Fix it:

```powershell
# Check current state
.\migrate.ps1 version

# Force to last known good version
.\migrate.ps1 force 1

# Try migration again
.\migrate.ps1 up
```

### "No change" error when running up

This is normal - it means all migrations are already applied.

### Can't connect to database

Check:

1. PostgreSQL container is running: `docker ps`
2. Connection string is correct in `conf/app.conf`
3. Network connectivity to database

### Migration not found

Ensure:

1. You're in the correct directory (`devops-api/`)
2. `migrations/` folder exists with `.sql` files
3. File naming follows pattern: `000001_name.up.sql`

## 📊 Migration State

Migrations are tracked in the `schema_migrations` table:

```sql
-- Check migration status manually
SELECT * FROM schema_migrations;

-- Output:
-- version | dirty
-- --------|-------
--      1  | false
```

- `version`: Current migration number
- `dirty`: `true` if last migration failed, `false` if clean

## 🔗 Resources

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [Migration Best Practices](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md)
- [PostgreSQL ALTER TABLE](https://www.postgresql.org/docs/current/sql-altertable.html)

## 🆚 Comparison with Other Tools

| Feature          | golang-migrate | Alembic (Python) | Prisma (TypeScript) |
| ---------------- | -------------- | ---------------- | ------------------- |
| Language         | Go             | Python           | TypeScript/Node     |
| Auto-generate    | No             | Yes              | Yes                 |
| Rollback         | ✅ Yes         | ✅ Yes           | ⚠️ Limited          |
| CLI Tool         | ✅ Yes         | ✅ Yes           | ✅ Yes              |
| Production Ready | ✅ Yes         | ✅ Yes           | ✅ Yes              |
| Multiple DBs     | ✅ Yes         | ✅ Yes           | ✅ Yes              |

golang-migrate doesn't auto-generate from models, but it gives you full SQL control and is the Go standard.

---

**Your migrations are ready!** The application will automatically apply them on startup. 🚀
