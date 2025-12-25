# Documentation Index

Welcome to Gantry! Here's what you need to know:

## 🚀 Getting Started

**New to Gantry?** Start here:
1. [Quick Start - 5 minutes](docs/DEPLOYMENT.md#quick-start-5-minutes) - Get running immediately
2. [Setup Guide](docs/SETUP.md) - Detailed installation for different setups
3. [Example Workflows](workflows/) - Real workflow examples to try

## 📖 Core Documentation

| Guide | Purpose |
|-------|---------|
| **[API Reference](docs/API.md)** | REST API endpoints and usage |
| **[Workflow Syntax](docs/WORKFLOWS.md)** | How to write YAML workflows |
| **[Testing Guide](docs/TESTING.md)** | Running tests locally & CI/CD pipeline |
| **[Deployment Guide](docs/DEPLOYMENT.md)** | Production deployment & checklists |
| **[Troubleshooting](docs/TROUBLESHOOTING.md)** | Common issues & solutions |

## 💡 Common Tasks

### Run Gantry Locally
```bash
cd backend && go run ./cmd/server/main.go
# In new terminal:
cd frontend && npm start
```
→ See [Setup Guide](docs/SETUP.md)

### Write Your First Workflow
```yaml
name: My Workflow
on:
  push:
    branches: [main]
jobs:
  test:
    runs-on: ubuntu
    steps:
      - name: Say Hello
        run: echo "Hello Gantry!"
```
→ See [Workflow Syntax](docs/WORKFLOWS.md)

### Run Tests
```bash
# Backend
cd backend && go test -v ./...

# Frontend  
cd frontend && npm test
```
→ See [Testing Guide](docs/TESTING.md)

### Deploy to Production
```bash
# Build backend
cd backend && go build ./cmd/server/main.go

# Build frontend
cd frontend && npm run build
```
→ See [Deployment Guide](docs/DEPLOYMENT.md)

## 🔍 Find Answers

**Something not working?** Check [Troubleshooting](docs/TROUBLESHOOTING.md) for:
- Backend issues (Docker, MongoDB, build errors)
- Frontend issues (npm, tests, build)
- API issues (workflows, runs, execution)
- CI/CD pipeline problems

## 🤝 Contributing

Want to contribute? See [CONTRIBUTING.md](CONTRIBUTING.md) for:
- How to fork and set up development
- Code style guidelines
- Testing requirements
- Pull request process

## 📚 Complete File Structure

```
docs/
├── API.md                 # REST API reference
├── DEPLOYMENT.md          # Production deployment & quick start
├── SETUP.md              # Installation & configuration
├── TESTING.md            # Testing & CI/CD
├── TROUBLESHOOTING.md    # Common issues & solutions
└── WORKFLOWS.md          # Workflow syntax guide

root/
├── README.md             # Project overview
├── CONTRIBUTING.md       # Contributing guidelines
├── .github/
│   ├── copilot-instructions.md
│   └── workflows/ci.yml  # CI/CD pipeline
└── workflows/            # Example workflows
```

## ⚡ Quick Links

- **📝 Code**: [backend/](backend/) • [frontend/](frontend/)
- **🧪 Tests**: `go test ./...` • `npm test`
- **🐳 Docker**: [docker-compose.yml](docker-compose.yml)
- **📋 Examples**: [workflows/](workflows/)
- **🔧 Config**: [.github/copilot-instructions.md](.github/copilot-instructions.md)

## 🆘 Need Help?

1. **Check the docs** - Most answers are in [Troubleshooting](docs/TROUBLESHOOTING.md)
2. **Run diagnostic** - Use the quick diagnostic in [Troubleshooting](docs/TROUBLESHOOTING.md#quick-diagnostic)
3. **Open an issue** - Include error message, steps, and system info
4. **Start a discussion** - For questions or ideas

---

**Ready to build?** Start with [Quick Start](docs/DEPLOYMENT.md#quick-start-5-minutes) 🚀
