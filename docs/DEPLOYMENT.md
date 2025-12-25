# Deployment Guide

## Quick Start (5 minutes)

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker (for running workflows)

### Step 1: Start Backend
```bash
cd backend
go mod download
go run ./cmd/server/main.go
```

Server starts on `http://localhost:8080`

### Step 2: Start Frontend (new terminal)
```bash
cd frontend
npm install
npm start
```

Browser opens at `http://localhost:3000`

### Step 3: Upload Your First Workflow

Create a file `my-workflow.yml`:
```yaml
name: Hello World

on:
  push:
    branches: [main]

jobs:
  greet:
    runs-on: ubuntu
    steps:
      - name: Say Hello
        run: echo "Hello from Gantry!"
```

Upload through the UI and click the play button!

## Docker Compose (Recommended)

```bash
docker-compose up -d
```

This starts:
- MongoDB: http://localhost:27017
- Backend: http://localhost:8080
- Frontend: http://localhost:3000

## Pre-Deployment Checklist

### Backend Verification
```bash
cd backend

# Build check
go build -o /tmp/test-build ./cmd/server/main.go

# Format check
go fmt ./...

# Lint check
go vet ./...

# Test check - all tests must pass
go test -v -race ./...
go test -v -race -tags=integration ./internal/server/...

# Coverage check - must be ≥60%
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
```

### Frontend Verification
```bash
cd frontend

# Install dependencies
npm ci

# Lint check
npx eslint src/
npx prettier --check "src/**/*.{js,jsx,json,css}"

# Test check - all tests must pass
npm test -- --coverage --watchAll=false
npm test -- App.test.integration.js --coverage --watchAll=false

# Build check
npm run build
```

### Security Scanning
```bash
# Go security
gosec ./backend/...

# JavaScript security
cd frontend && npm audit --audit-level=moderate
```

## Manual Deployment

### Build Backend
```bash
cd backend
go build -o gantry-server ./cmd/server/main.go
```

### Build Frontend
```bash
cd frontend
npm run build
```

### Environment Variables

```bash
# Backend
export PORT=8080
export STORAGE_TYPE=mongodb  # or 'memory'
export MONGO_URI=mongodb://localhost:27017
export MONGO_DATABASE=gantry

# Start server
./gantry-server
```

### Deploy Frontend
```bash
# Serve the build/ directory with nginx or similar
# Or use Node.js server:
cd frontend
npm install -g serve
serve -s build -l 3000
```

## Production Checklist

- [ ] Enable HTTPS/TLS
- [ ] Use production database (MongoDB Atlas recommended)
- [ ] Enable authentication (JWT or OAuth2)
- [ ] Set up secrets management
- [ ] Enable logging and monitoring
- [ ] Configure backups
- [ ] Set resource limits
- [ ] Enable rate limiting
- [ ] Run security scanning (Gosec, npm audit)
- [ ] All tests passing (60%+ coverage)
- [ ] All CI/CD jobs passing

## Troubleshooting

### Docker Connection Issues
```bash
# Verify Docker daemon
docker ps

# On Linux, add user to docker group
sudo usermod -aG docker $USER
```

### Port Already in Use
```bash
# Change backend port
export PORT=8081
./gantry-server

# Change frontend port
cd frontend
PORT=3001 npm start
```

### MongoDB Connection
```bash
# Test MongoDB connection
mongo mongodb://localhost:27017

# Check MongoDB logs
docker logs gantry-mongodb
```

### CORS Issues
- Ensure frontend URL matches API origin
- Check CORS middleware configuration in backend

## CI/CD Pipeline

The GitHub Actions pipeline runs 10 jobs automatically:

1. **Backend Lint** - gofmt, go vet, golangci-lint
2. **Backend Tests** - unit + integration tests (60%+ required)
3. **Backend Security** - Gosec scanning
4. **Backend Build** - go build verification
5. **Frontend Lint** - ESLint, Prettier
6. **Frontend Tests** - unit + integration tests
7. **Frontend Security** - npm audit
8. **Frontend Build** - npm run build verification
9. **All Checks Passed** - Gate requiring all jobs to pass

View the pipeline at `.github/workflows/ci.yml`

## Rollback

If issues occur after deployment:

```bash
# Revert to previous version
git revert HEAD

# Or checkout previous commit
git checkout <commit-hash>

# Rebuild and redeploy
go build ./cmd/server/main.go
npm run build
```

## Monitoring

Monitor these metrics:
- Build success rate
- Test coverage trends
- Security scan results
- Response times
- Error rates
- Workflow execution times

## Support

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common issues or open an issue on GitHub.
