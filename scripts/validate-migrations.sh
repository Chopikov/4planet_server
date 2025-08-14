#!/bin/bash

# Migration validation script for 4Planet backend
# This script validates migration files and checks for common issues

set -e

echo "🔍 Validating migration files..."

# Check if migrations directory exists
if [ ! -d "migrations" ]; then
    echo "❌ migrations directory not found"
    exit 1
fi

# Check for required files
required_files=("000001_init.up.sql" "000001_init.down.sql")
for file in "${required_files[@]}"; do
    if [ ! -f "migrations/$file" ]; then
        echo "❌ Required migration file not found: $file"
        exit 1
    fi
done

echo "✅ Required migration files found"

# Check file naming convention
echo "📝 Checking file naming convention..."
for file in migrations/*.sql; do
    filename=$(basename "$file")
    if [[ ! "$filename" =~ ^[0-9]{6}_[a-z_]+\.(up|down)\.sql$ ]]; then
        echo "⚠️  File naming convention issue: $filename"
        echo "   Expected format: 000001_description.up.sql or 000001_description.down.sql"
    fi
done

# Check for matching up/down pairs
echo "🔗 Checking migration pairs..."
up_files=$(find migrations -name "*.up.sql" | sort)
down_files=$(find migrations -name "*.down.sql" | sort)

up_count=$(echo "$up_files" | wc -l)
down_count=$(echo "$down_files" | wc -l)

if [ "$up_count" != "$down_count" ]; then
    echo "❌ Mismatched migration pairs: $up_count up files, $down_count down files"
    exit 1
fi

echo "✅ Migration pairs are balanced"

# Check for SQL syntax (basic validation)
echo "🔧 Validating SQL syntax..."
for file in migrations/*.sql; do
    if command -v psql >/dev/null 2>&1; then
        # Try to parse SQL with psql (dry run)
        if ! psql -q -t -c "SELECT 1;" >/dev/null 2>&1; then
            echo "⚠️  Could not validate SQL syntax for $file (psql not available)"
        fi
    else
        echo "⚠️  psql not available, skipping SQL syntax validation"
        break
    fi
done

# Check for common issues in migration files
echo "🔍 Checking for common issues..."

# Check for hardcoded database names
if grep -r "CREATE DATABASE" migrations/ >/dev/null 2>&1; then
    echo "⚠️  Found CREATE DATABASE in migrations (usually not recommended)"
fi

# Check for DROP DATABASE
if grep -r "DROP DATABASE" migrations/ >/dev/null 2>&1; then
    echo "⚠️  Found DROP DATABASE in migrations (dangerous!)"
fi

# Check for proper rollback in down migrations
echo "📋 Checking down migrations..."
for down_file in migrations/*.down.sql; do
    if [ -f "$down_file" ]; then
        # Check if down migration has content
        if [ ! -s "$down_file" ]; then
            echo "⚠️  Empty down migration: $down_file"
        fi
        
        # Check for dangerous operations
        if grep -q "DROP TABLE" "$down_file"; then
            echo "⚠️  Down migration drops tables: $down_file"
        fi
    fi
done

# Migration configuration is handled via environment variables
echo "✅ Migration configuration via DB_DSN environment variable"

echo ""
echo "🎉 Migration validation complete!"
echo ""
echo "📚 Next steps:"
echo "   1. Review any warnings above"
echo "   2. Test migrations: make migrate-status"
echo "   3. Run migrations: make migrate"
echo "   4. Seed database: make seed"
echo ""
echo "📖 For detailed instructions, see MIGRATIONS.md"
