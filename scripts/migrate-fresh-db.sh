#!/bin/bash

# Config
DB_NAME="livestudio"
DB_USER="postgres"
DB_PASS="postgres"  
DB_HOST="127.0.0.1"
DB_PORT="5432"
MIGRATION_PATH="D:\MKI\live-studio-api\cmd\migrate\main.go"
SEEDER_ROLE_PATH="D:\MKI\live-studio-api\cmd\seed\main.go"

export PGPASSWORD=$DB_PASS

echo "🔄 Dropping database $DB_NAME ..."
psql -U $DB_USER -h $DB_HOST -p $DB_PORT -c "DROP DATABASE IF EXISTS $DB_NAME;"

echo "🆕 Creating database $DB_NAME ..."
psql -U $DB_USER -h $DB_HOST -p $DB_PORT -c "CREATE DATABASE $DB_NAME;"

echo "🎉 Database $DB_NAME created successfully!"

echo "⚡ Running migrations ..."
go run $MIGRATION_PATH

echo " ⚡ Running seed roles ..."
go run $SEEDER_ROLE_PATH role-permission

echo "✅ Done!"
air
