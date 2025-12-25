# Troubleshooting Guide

## Backend Issues

### Server Won't Start

**Error: "Address already in use"**
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use a different port
export PORT=8081
go run ./cmd/server/main.go
```

**Error: "Docker daemon not accessible"**
```bash
# Verify Docker is running
docker ps

# On Linux, add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# On macOS, restart Docker Desktop
open -a Docker
```

**Error: "MongoDB connection refused"**
```bash
# Start MongoDB
docker run -d -p 27017:27017 --name gantry-db mongo:7

# Or with docker-compose
docker-compose up -d

# Verify connection
mongo mongodb://localhost:27017
```

### Tests Failing

**"go test: no Go files in /path/to/directory"**
```bash
# Make sure you're in backend directory
cd backend
go test ./...
```

**"Race condition detected"**
```bash
# Run tests with race detector
go test -race ./...

# Check for goroutine leaks
go test -race -count=5 ./...
```

**"Coverage below 60%"**
```bash
# Check coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# View HTML coverage report
go tool cover -html=coverage.out
```

### Build Failures

**"go: package not found"**
```bash
# Download dependencies
go mod download

# Tidy go.mod
go mod tidy

# Rebuild
go build ./cmd/server/main.go
```

**"lint: unexpected error"**
```bash
# Check code formatting
go fmt ./...

# Run vet
go vet ./...

# Install golangci-lint
brew install golangci-lint
golangci-lint run ./...
```

## Frontend Issues

### npm Install Fails

**"npm ERR! peer dep missing"**
```bash
# Install missing peer dependencies
npm install --legacy-peer-deps

# Or update npm
npm install -g npm@latest
```

**"npm ERR! ERESOLVE could not resolve dependency"**
```bash
# Use legacy resolver
npm install --legacy-peer-deps

# Or delete node_modules and lock file
rm -rf node_modules package-lock.json
npm install
```

### Tests Failing

**"Cannot find module 'jest'"**
```bash
# Install test dependencies
npm install --save-dev jest @testing-library/react

# Or reinstall all dependencies
npm ci
```

**"FAIL: React component not rendering"**
```bash
# Clear Jest cache
npm test -- --clearCache

# Run with verbose output
npm test -- --verbose

# Run specific test
npm test -- App.test.integration.js
```

### Build Issues

**"npm run build fails"**
```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
npm run build
```

**"Build succeeds but app doesn't load"**
```bash
# Check for console errors
# Open http://localhost:3000 in browser
# Check Network tab and Console for errors

# Verify API is running
curl http://localhost:8080/api/workflows

# Check CORS headers
curl -v http://localhost:8080/api/workflows
```

## Docker Issues

### Cannot Connect to Docker Daemon

**"Cannot connect to the Docker daemon"**
```bash
# macOS
open -a Docker

# Linux
sudo systemctl start docker

# Verify connection
docker ps
```

### Container Won't Start

**"Error response from daemon"**
```bash
# Check logs
docker logs <container-name>

# Inspect container
docker inspect <container-name>

# Stop and remove
docker stop <container-name>
docker rm <container-name>

# Retry
docker-compose up -d
```

### Volume Mount Issues

**"Cannot mount volume"**
```bash
# Check volume exists
docker volume ls

# Create volume
docker volume create gantry-data

# Remove and recreate
docker volume rm gantry-data
docker volume create gantry-data
```

## API Issues

### Workflow Upload Fails

**"400 Bad Request: Failed to parse workflow"**
```bash
# Check YAML syntax
yamllint my-workflow.yml

# Or validate manually
python3 -c "import yaml; yaml.safe_load(open('my-workflow.yml'))"

# Check required fields
# - name
# - on.push.branches
# - jobs
# - jobs[name].runs-on
# - jobs[name].steps
```

### Workflow Execution Fails

**"500 Internal Server Error"**
```bash
# Check backend logs
# Verify Docker connection
docker ps

# Check workflow in database
# Run specific job in Docker manually
docker run -it ubuntu:latest echo "test"
```

**"Workflow runs but no output"**
```bash
# Check step output in UI
# View logs in "Recent Runs" section

# Or get via API
curl http://localhost:8080/api/runs/<run-id>
```

## Database Issues

### MongoDB Connection Failed

**"Error: connect ECONNREFUSED 127.0.0.1:27017"**
```bash
# Start MongoDB
docker run -d -p 27017:27017 mongo:7

# Or via docker-compose
docker-compose up -d

# Verify running
docker ps | grep mongo
```

### Data Not Persisting

**"Data lost after restart"**
```bash
# Check storage type
echo $STORAGE_TYPE  # Should be 'mongodb' not 'memory'

# Restart with MongoDB
export STORAGE_TYPE=mongodb
export MONGO_URI=mongodb://localhost:27017
go run ./cmd/server/main.go
```

### MongoDB Performance

**"Queries are slow"**
```bash
# Check indexes
mongo mongodb://localhost:27017
db.workflows.getIndexes()
db.runs.getIndexes()

# Monitor connections
db.serverStatus().connections
```

## CI/CD Issues

### GitHub Actions Failing

**"Workflow job fails at lint stage"**
```bash
# Run locally
cd backend
go fmt ./...
go vet ./...
golangci-lint run ./...
```

**"Coverage check fails"**
```bash
# Check local coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total

# Must be ≥60%
```

**"Tests pass locally but fail in CI"**
```bash
# Check for race conditions
go test -race ./...

# Run with same Go version as CI
go version  # Check CI uses same version

# Check for hardcoded paths
grep -r "/Users/" backend/
```

## Performance Issues

### Slow Workflow Execution

**"Workflow takes too long"**
```bash
# Check Docker performance
docker stats

# Check resource limits
docker inspect <container-name> | grep -A 20 'HostConfig'

# Increase resources
docker run -m 2g --cpus 2 ...
```

### High Memory Usage

**"Process using too much memory"**
```bash
# Monitor process
top
ps aux | grep gantry

# Check for goroutine leaks
curl http://localhost:8080/debug/pprof/goroutine

# Or check in logs
go run ./cmd/server/main.go 2>&1 | grep -i "goroutine"
```

## Logging & Debugging

### Enable Debug Mode

```bash
# Backend
export LOG_LEVEL=debug
go run ./cmd/server/main.go

# Check logs
tail -f gantry.log
```

### Get API Response Details

```bash
# Verbose curl output
curl -v http://localhost:8080/api/workflows

# Pretty print JSON
curl http://localhost:8080/api/runs | jq .

# Check headers
curl -i http://localhost:8080/api/workflows
```

### Browser DevTools

1. Open http://localhost:3000
2. Right-click → Inspect → Console
3. Check for JavaScript errors
4. Check Network tab for failed requests
5. Check Application → Local Storage for saved state

## Common Error Messages

| Error | Cause | Solution |
|-------|-------|----------|
| "dial tcp: lookup docker.sock" | Docker not accessible | Start Docker, add user to docker group |
| "connection refused" | Service not running | Start backend/MongoDB |
| "404 Not Found" | Route doesn't exist | Check API endpoint URL |
| "CORS error" | Origin not allowed | Check CORS middleware config |
| "Workflow syntax error" | Invalid YAML | Validate YAML structure |
| "No space left on device" | Disk full | Clean up Docker volumes |
| "Permission denied" | File permissions | Check file ownership and chmod |

## Getting Help

1. Check [TESTING.md](TESTING.md) for test-related issues
2. Check [DEPLOYMENT.md](DEPLOYMENT.md) for deployment questions
3. Check [API.md](API.md) for API-related problems
4. Open an issue on GitHub with:
   - Exact error message
   - Steps to reproduce
   - System information (OS, Go version, Node version, Docker version)
   - Relevant logs

## Quick Diagnostic

Run this to check system status:

```bash
#!/bin/bash
echo "=== System Status ==="
echo "Go version: $(go version)"
echo "Node version: $(node -v)"
echo "npm version: $(npm -v)"
echo "Docker: $(docker --version)"
echo "Docker daemon: $(docker ps > /dev/null && echo 'Running' || echo 'Not running')"
echo ""
echo "=== Backend Status ==="
curl -s http://localhost:8080/api/workflows > /dev/null && echo "✅ Backend running" || echo "❌ Backend not running"
echo ""
echo "=== Frontend Status ==="
curl -s http://localhost:3000 > /dev/null && echo "✅ Frontend running" || echo "❌ Frontend not running"
echo ""
echo "=== Database Status ==="
mongo mongodb://localhost:27017 --eval "db.adminCommand('ping')" > /dev/null 2>&1 && echo "✅ MongoDB running" || echo "❌ MongoDB not running"
```

Save as `diagnostic.sh` and run: `bash diagnostic.sh`
