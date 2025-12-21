#!/bin/bash

# Frontend Linting Script

echo "🔍 Running Frontend Linters"
echo "============================"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if we're in frontend directory
if [ ! -f "package.json" ]; then
    echo "${RED}Error: Please run this script from the frontend directory${NC}"
    exit 1
fi

# Install eslint if not present
if ! npm list eslint > /dev/null 2>&1; then
    echo "${YELLOW}Installing ESLint...${NC}"
    npm install --save-dev eslint
fi

# Install prettier if not present
if ! npm list prettier > /dev/null 2>&1; then
    echo "${YELLOW}Installing Prettier...${NC}"
    npm install --save-dev prettier
fi

echo "📝 Running ESLint..."
npx eslint src/ --max-warnings 0

if [ $? -eq 0 ]; then
    echo "${GREEN}✓ ESLint passed${NC}"
    ESLINT_PASS=true
else
    echo "${RED}✗ ESLint failed${NC}"
    ESLINT_PASS=false
fi

echo ""
echo "💅 Running Prettier..."
npx prettier --check "src/**/*.{js,jsx,json,css}"

if [ $? -eq 0 ]; then
    echo "${GREEN}✓ Prettier passed${NC}"
    PRETTIER_PASS=true
else
    echo "${RED}✗ Prettier failed${NC}"
    echo ""
    echo "Run ${YELLOW}npx prettier --write \"src/**/*.{js,jsx,json,css}\"${NC} to fix"
    PRETTIER_PASS=false
fi

echo ""
echo "============================"

if [ "$ESLINT_PASS" = true ] && [ "$PRETTIER_PASS" = true ]; then
    echo "${GREEN}✅ All linting checks passed!${NC}"
    exit 0
else
    echo "${RED}❌ Some linting checks failed${NC}"
    exit 1
fi