#!/bin/bash

# Gantry Test Runner

echo "🧪 Running Gantry Backend Tests"
echo "================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if we're in the backend directory
if [ ! -f "go.mod" ]; then
    echo "${RED}Error: Please run this script from the backend directory${NC}"
    exit 1
fi

# Run tests with coverage
echo "📊 Running tests with coverage..."
go test -v -cover -coverprofile=coverage.out ./... 2>&1 | tee test-output.log

# Check if tests passed
if [ $? -eq 0 ]; then
    echo ""
    echo "${GREEN}✅ All tests passed!${NC}"
    echo ""
    
    # Generate coverage report
    echo "📈 Coverage Report:"
    go tool cover -func=coverage.out | tail -1
    echo ""
    
    # Generate HTML coverage report
    echo "🌐 Generating HTML coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    echo "${GREEN}✓${NC} Coverage report saved to coverage.html"
    echo ""
    
    # Show summary
    echo "📝 Test Summary:"
    grep -E "^(PASS|FAIL)" test-output.log | sort | uniq -c
    echo ""
    
    exit 0
else
    echo ""
    echo "${RED}❌ Tests failed!${NC}"
    echo ""
    echo "Failed tests:"
    grep "FAIL" test-output.log
    echo ""
    exit 1
fi