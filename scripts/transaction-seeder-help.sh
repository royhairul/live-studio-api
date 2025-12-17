#!/bin/bash
# Transaction Seeder Quick Reference Script
# This script provides examples and shortcuts for running transaction seeders

echo "==================================="
echo "Transaction Seeder Quick Reference"
echo "==================================="
echo ""

# Color codes for better readability
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}Available Commands:${NC}"
echo ""

echo -e "${GREEN}1. Seed single month:${NC}"
echo "   go run cmd/seed/main.go transaction-month <tenant_id> <year> <month>"
echo ""
echo "   Examples:"
echo "   go run cmd/seed/main.go transaction-month tenant-123 2024 11"
echo "   go run cmd/seed/main.go transaction-month studio-abc 2024 1"
echo ""

echo -e "${GREEN}2. Seed date range:${NC}"
echo "   go run cmd/seed/main.go transaction-range <tenant_id> <start_year> <start_month> <end_year> <end_month>"
echo ""
echo "   Examples:"
echo "   go run cmd/seed/main.go transaction-range tenant-123 2024 1 2024 12   # Full year"
echo "   go run cmd/seed/main.go transaction-range tenant-123 2024 10 2024 12  # Q4 2024"
echo ""

echo -e "${YELLOW}Common Use Cases:${NC}"
echo ""

echo "Seed last 3 months for tenant-123:"
echo "go run cmd/seed/main.go transaction-range tenant-123 2024 10 2024 12"
echo ""

echo "Seed November 2024 for studio-abc:"
echo "go run cmd/seed/main.go transaction-month studio-abc 2024 11"
echo ""

echo "Seed entire year 2024 for tenant-456:"
echo "go run cmd/seed/main.go transaction-range tenant-456 2024 1 2024 12"
echo ""

echo -e "${BLUE}Need Help?${NC}"
echo "See: database/seeders/README_TRANSACTION_SEEDER.md"
echo ""
