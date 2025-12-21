#!/bin/bash

# Backend Linting Script

echo "🔍 Running Backend Linters"
echo "=========================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if we're in backend directory
if [ ! -f "go.mod" ]; then
    echo "${RED}Error: Please run this script from the backend directory${NC}"
    exit 1
fi

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "${YELLOW}golangci-lint not found. Installing...${NC}"
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
fi

echo "📝 Running gofmt..."
GOFMT_FILES=$(gofmt -l .)
if [ -z "$GOFMT_FILES" ]; then
    echo "${GREEN}✓ gofmt passed${NC}"
    GOFMT_PASS=true
else
    echo "${RED}✗ gofmt failed${NC}"
    echo "Files need formatting:"
    echo "$GOFMT_FILES"
    echo ""
    echo "Run ${YELLOW}gofmt -w .${NC} to fix"
    GOFMT_PASS=false
fi

echo ""
echo "📝 Running go vet..."
go vet ./...

if [ $? -eq 0 ]; then
    echo "${GREEN}✓ go vet passed${NC}"
    GOVET_PASS=true
else
    echo "${RED}✗ go vet failed${NC}"
    GOVET_PASS=false
fi

echo ""
echo "📝 Running golangci-lint..."
golangci-lint run ./...

if [ $? -eq 0 ]; then
    echo "${GREEN}✓ golangci-lint passed${NC}"
    GOLANGCI_PASS=true
else
    echo "${RED}✗ golangci-lint failed${NC}"
    GOLANGCI_PASS=false
fi

echo ""
echo "=========================="

if [ "$GOFMT_PASS" = true ] && [ "$GOVET_PASS" = true ] && [ "$GOLANGCI_PASS" = true ]; then
    echo "${GREEN}✅ All linting checks passed!${NC}"
    exit 0
else
    echo "${RED}❌ Some linting checks failed${NC}"
    exit 1
fi