# Database Migrations Guide

This guide explains how to use the database migration system for the 4Planet backend. The system is built around **GORM models** and ensures your database schema always perfectly matches your Go code.

## 🎯 **Overview**

### **Key Benefits**
- ✅ **Zero drift** between GORM models and database schema
- ✅ **Automatic validation** through migration testing
- ✅ **Consistent structure** across all environments
- ✅ **Simple commands** with `make` targets
- ✅ **Clear documentation** for all operations

## 📝 **Summary**

### **Quick Reference**

```bash
# Development setup
make deps
make docker-run
make migrate
make seed

# Production deployment
migrate -path migrations -database "$PROD_DSN" up

# Rollback if needed
migrate -path migrations -database "$PROD_DSN" down 1
```

### **Key Commands**

- `make migrate` - Run all pending migrations
- `make seed` - Seed database with initial data
- `migrate version` - Check current migration version
- `migrate up` - Apply pending migrations
- `migrate down` - Rollback migrations

### **File Locations**

- **Migrations**: `./migrations/`
- **Seed Data**: `./cmd/seed/main.go` (Go-based)
- **Environment**: `./env.example`

### **Go-Based Seeding**

The seeding system uses Go code (`cmd/seed/main.go`) with GORM models for type-safe data insertion:

```bash
# Check if Go seeder exists
ls -la cmd/seed/main.go

# Create seeder if needed
mkdir -p cmd/seed
touch cmd/seed/main.go

# Add your seeding logic using GORM models
# Example:
# func seedTreePrices(ctx context.Context, db *gorm.DB) error {
#     prices := []models.TreePrice{
#         {Currency: "USD", PriceMinor: 2500},
#     }
#     return db.Create(&prices).Error
# }
```

**Benefits of Go-based seeding**:
- **Type Safety**: Compile-time validation of data structure
- **Model Consistency**: Always matches your GORM models exactly
- **Easier Maintenance**: Update models, seeding automatically stays in sync
- **Better Testing**: Can test seeding logic with Go tests
- **Environment Flexibility**: Easy to seed different environments differently


## 📋 **Prerequisites**

### 1. Install golang-migrate

```bash
# macOS (using Homebrew)
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

```

### 2. Verify Installation

```bash
migrate --version
# Should output: migrate version v4.18.3
```

## 🚀 **Quick Start**

### **Option 1: Using Makefile (Recommended)**

```bash
# Install dependencies first
make deps

# Validate migrations
make migrate-validate

# Start development environment
make docker-run

# Run migrations
make migrate

# Seed database
make seed

# View migration status
make migrate-status
```

### **Option 2: Manual Migration Commands**

```bash
# Set your database connection string (optional)
export DB_DSN="postgres://postgres:postgres@localhost:5432/planet?sslmode=disable"

# Run all pending migrations
migrate -path migrations -database "$DB_DSN" up

# Seed the database (using Go seeder)
go run ./cmd/seed
```

## 📚 **Migration Commands**

### **Available Make Targets**

```bash
make migrate          # Run all pending migrations
make migrate-status   # Check current migration version
make migrate-down     # Rollback last migration
make migrate-reset    # Rollback all migrations
make migrate-create   # Create new migration files
make migrate-validate # Validate migration files
```

### **Basic Commands**

```bash
# Run all pending migrations
migrate -path migrations -database "$DB_DSN" up

# Rollback last migration
migrate -path migrations -database "$DB_DSN" down 1

# Rollback all migrations
migrate -path migrations -database "$DB_DSN" down

# Check migration status
migrate -path migrations -database "$DB_DSN" version

# Force migration version (use with caution)
migrate -path migrations -database "$DB_DSN" force VERSION
```

### **Advanced Commands**

```bash
# Run migrations with custom table name
migrate -path migrations -database "$DB_DSN" -x migrations-table=4planet_migrations up

# Run migrations with timeout
migrate -path migrations -database "$DB_DSN" -x timeout=30s up

# Run migrations with custom source
migrate -path migrations -database "$DB_DSN" -source file://migrations up

# Dry run (validate migrations without applying)
migrate -path migrations -database "$DB_DSN" -dry-run up
```

## 🔧 **Environment Configuration**

### **Database Connection String Format**

```bash
# Local development
DB_DSN="postgres://postgres:postgres@localhost:5432/planet?sslmode=disable"

# Docker Compose
DB_DSN="postgres://postgres:postgres@postgres:5432/planet?sslmode=disable"

# Production
DB_DSN="postgres://user:password@host:5432/planet?sslmode=require"
```

### **Environment Variables**

```bash
# Copy example environment file
cp env.example .env

# Edit .env file with your database settings
DB_DSN=postgres://postgres:postgres@localhost:5432/planet?sslmode=disable
```

### **Makefile Environment Support**

The Makefile automatically uses the `DB_DSN` environment variable if set:

```bash
# Use default connection (local development)
make migrate

# Use custom connection
export DB_DSN="postgres://user:pass@host:5432/prod_db?sslmode=require"
make migrate

# Use different connection for specific command
DB_DSN="postgres://test:test@localhost:5432/test_db" make migrate
```

**Default connection**: `postgres://postgres:postgres@localhost:5432/planet?sslmode=disable`

## 📁 **Migration File Structure**

### **File Naming Convention**

```
migrations/
├── 000001_init.up.sql          # Up migration
├── 000001_init.down.sql        # Down migration
├── 000002_add_feature.up.sql   # Future migration
├── 000002_add_feature.down.sql # Future rollback
```

### **Migration File Format**

```sql
-- 000001_init.up.sql
-- Up migration: creates tables, indexes, etc.

-- 000001_init.down.sql  
-- Down migration: drops tables, indexes, etc.
```

## 🐳 **Docker Environment**

### **Start Services**

```bash
# Start PostgreSQL and other services
make docker-run

# Wait for services to be ready
sleep 10

# Run migrations
make migrate

# Seed database
make seed
```

### **Migration in Docker**

```bash
# Run migrations from host (recommended)
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/planet?sslmode=disable" up

# Or use migrate from inside container
docker-compose exec postgres psql -U postgres -d planet -f /docker-entrypoint-initdb.d/01-schema.sql
```

## 🧪 **Testing Migrations**

### **Test Migration Rollback**

```bash
# Apply migration
migrate -path migrations -database "$DB_DSN" up

# Verify tables exist
psql "$DB_DSN" -c "\dt"

# Rollback migration
migrate -path migrations -database "$DB_DSN" down 1

# Verify tables are dropped
psql "$DB_DSN" -c "\dt"
```

### **Validate Migration Files**

```bash
# Check for syntax errors
migrate -path migrations -database "$DB_DSN" -dry-run up

# Validate file naming
ls migrations/*.sql | grep -E "^migrations/[0-9]{6}_[a-z_]+\.(up|down)\.sql$"

# Run automated validation
make migrate-validate
```

## 🔄 **Creating New Migrations**

### **Generate Migration Files**

```bash
# Create new migration
make migrate-create
# Enter: add_user_phone

# This creates:
# migrations/000002_add_user_phone.up.sql
# migrations/000002_add_user_phone.down.sql
```

### **Migration Template**

```sql
-- 000002_add_user_phone.up.sql
-- Add new feature

-- Example: Add new column
ALTER TABLE users ADD COLUMN phone_number text;

-- Example: Add new index
CREATE INDEX idx_users_phone ON users(phone_number);
```

```sql
-- 000002_add_user_phone.down.sql
-- Rollback new feature

-- Example: Remove new column
ALTER TABLE users DROP COLUMN phone_number;

-- Example: Remove new index
DROP INDEX IF EXISTS idx_users_phone;
```

## 🎯 **GORM Integration**

### **How Migrations Work with GORM**

1. **GORM Models Define Schema**: Your Go structs with GORM tags define the database structure
2. **Migrations Apply Schema**: SQL migrations create the actual database tables and constraints
3. **Perfect Alignment**: The migration system ensures your database matches your GORM models exactly

### **GORM Auto-Migration**

```go
// In development, you can use GORM's auto-migration
func autoMigrate() error {
    return DB.AutoMigrate(
        &models.User{},
        &models.Session{},
        &models.Project{},
        // ... all your models
    )
}
```

### **Migration vs Auto-Migration**

- **Migrations**: For production, version control, rollbacks
- **Auto-Migration**: For development, prototyping, testing
- **Best Practice**: Use migrations for production, auto-migration for development

## 🚨 **Troubleshooting**

### **Common Issues**

#### **1. Connection Refused**

```bash
# Error: dial error (dial tcp [::1]:5432: connect: connection refused)

# Solution: Check if PostgreSQL is running
brew services list | grep postgresql
# or
docker-compose ps
```

#### **2. Database Does Not Exist**

```bash
# Error: database "planet" does not exist

# Solution: Create database
createdb -U postgres planet
# or
psql -U postgres -c "CREATE DATABASE planet;"
```

#### **3. Permission Denied**

```bash
# Error: permission denied for table

# Solution: Check user permissions
psql -U postgres -c "\du"
# Grant necessary permissions
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE planet TO postgres;"
```

#### **4. Migration Already Applied**

```bash
# Error: migration already applied

# Solution: Check current version
migrate -path migrations -database "$DB_DSN" version

# Force version if needed (use with caution)
migrate -path migrations -database "$DB_DSN" force VERSION
```

### **Reset Database**

```bash
# Drop and recreate database (DANGER: loses all data)
dropdb -U postgres planet
createdb -U postgres planet

# Run migrations from scratch
migrate -path migrations -database "$DB_DSN" up
```

## 📊 **Migration Status**

### **Check Current Status**

```bash
# View applied migrations
migrate -path migrations -database "$DB_DSN" version

# View migration history
psql "$DB_DSN" -c "SELECT * FROM schema_migrations ORDER BY version;"
```

### **Migration Logs**

```bash
# View PostgreSQL logs
docker-compose logs postgres

# View application logs
docker-compose logs api
```

## 🎯 **Best Practices**

### **1. Always Test Rollbacks**

```bash
# Test migration
make migrate

# Test rollback
make migrate-down

# Re-apply
make migrate
```

### **2. Use Transactions**

```sql
-- Wrap complex migrations
BEGIN;

-- Your migration SQL here
ALTER TABLE users ADD COLUMN phone_number text;

COMMIT;
-- or ROLLBACK; if something goes wrong
```

### **3. Validate Before Committing**

```bash
# Run validation
make migrate-validate

# Fix any warnings
# Commit only after validation passes
```

### **4. Test with GORM Models**

```bash
# After running migrations, verify GORM can connect
go run ./cmd/api

# Check that models can be created/queried
# This ensures perfect alignment between code and database
```

### **5. Backup Before Production**

```bash
# Create backup before running migrations
pg_dump -U postgres planet > backup_before_migration.sql

# Run migrations
migrate -path migrations -database "$DB_DSN" up

# Verify everything works
# If not, restore from backup
psql -U postgres planet < backup_before_migration.sql
```


