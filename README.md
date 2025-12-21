# 🏗️ Gantry

<div align="center">

**Lightweight, self-hosted CI/CD platform for small teams**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18+-61DAFB?style=flat&logo=react)](https://reactjs.org/)
[![Docker](https://img.shields.io/badge/Docker-Required-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

[Features](#-features) • [Quick Start](#-quick-start) • [Documentation](#-documentation) • [Architecture](#-architecture) • [Contributing](#-contributing)

</div>

---

## 🎯 What is Gantry?

Gantry is a lightweight, self-hosted CI/CD platform inspired by GitHub Actions. Built with Go and React, it's designed for small teams who want the power of automated workflows without the complexity of Jenkins or the cost of cloud CI/CD services.

### Why Gantry?

- ✅ **Simple Setup** - Get running in minutes, not hours
- ✅ **Self-Hosted** - Keep your code and secrets on your infrastructure  
- ✅ **YAML Workflows** - Familiar syntax if you've used GitHub Actions
- ✅ **Docker Isolation** - Each job runs in a clean container
- ✅ **Real-Time Monitoring** - Watch your builds live
- ✅ **Lightweight** - Minimal resource footprint
- ✅ **Modular** - Clean architecture, easy to extend

---

## ✨ Features

### Current Features

- 📝 **YAML Workflow Definitions** - Write workflows like GitHub Actions
- 🐳 **Docker Isolation** - Each job runs in a fresh container
- ⚡ **Sequential Execution** - Jobs execute in YAML order
- 📊 **Real-Time Dashboard** - Monitor builds as they happen
- 📜 **Complete Logs** - See every step's output with timestamps
- 🔄 **Auto-Refresh** - UI updates automatically
- 🎯 **Manual Triggers** - Start workflows with one click
- 💾 **Persistent Storage** - MongoDB support for production
- 🧪 **Unit Tests** - Comprehensive test coverage
- 🏗️ **Modular Architecture** - Clean, maintainable codebase

### Coming Soon

- [ ] **Git Webhooks** - Auto-trigger on push/PR
- [ ] **User Authentication** - JWT/OAuth2 support
- [ ] **Secrets Management** - Secure credential storage
- [ ] **Artifacts** - Save build outputs
- [ ] **Matrix Builds** - Test across multiple versions
- [ ] **Notifications** - Email/Slack alerts
- [ ] **Caching** - Speed up repeated builds
- [ ] **Parallel Jobs** - Run independent jobs concurrently
- [ ] **Distributed Runners** - Scale across machines
- [ ] **Plugin System** - Extend with custom actions

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** - [Install Go](https://go.dev/doc/install)
- **Node.js 18+** - [Install Node.js](https://nodejs.org/)
- **Docker** - [Install Docker](https://docs.docker.com/get-docker/)

### Installation

```bash
# Clone the repository
git clone https://github.com/sudzp/gantry.git
cd gantry

# Start backend
cd backend
go mod download
go run ./cmd/server/main.go

# Start frontend (in new terminal)
cd frontend
npm install
npm start
```

Visit http://localhost:3000 and you're ready to go! 🎉

### Using MongoDB (Optional)

```bash
# Start MongoDB with Docker
docker run -d --name gantry-mongodb -p 27017:27017 -v gantry-data:/data/db mongo:7

# Configure backend
cd backend
export STORAGE_TYPE=mongodb
export MONGO_URI=mongodb://localhost:27017
export MONGO_DATABASE=gantry
go run ./cmd/server/main.go
```

---

## 📋 Example Workflow

Create a file `my-workflow.yml`:

```yaml
name: Build and Test

on:
  push:
    branches:
      - main

jobs:
  lint:
    runs-on: alpine
    steps:
      - name: Lint code
        run: |
          echo "Running linter..."
          echo "✓ Lint passed"
          
  test:
    runs-on: ubuntu
    steps:
      - name: Run tests
        run: |
          echo "Running tests..."
          echo "✓ 150 tests passed"
          
  build:
    runs-on: ubuntu
    steps:
      - name: Build application
        run: |
          echo "Building application..."
          echo "✓ Build successful"
          
  deploy:
    runs-on: alpine
    steps:
      - name: Deploy to staging
        run: |
          echo "Deploying to staging..."
          echo "✓ Deployment complete"
```

Upload it through the UI, click the play button, and watch it run!

---

## 🏗️ Architecture

### Backend (Go)

```
backend/
├── cmd/server/              # Application entry point
├── internal/
│   ├── models/             # Data structures
│   ├── parser/             # YAML workflow parsing
│   ├── executor/           # Job execution (Docker)
│   ├── storage/            # Data persistence (Memory/MongoDB)
│   ├── api/                # HTTP handlers & routes
│   └── server/             # Server orchestration
└── go.mod
```

### Frontend (React)

```
frontend/
├── src/
│   ├── components/         # Reusable UI components
│   ├── services/           # API client
│   ├── hooks/              # Custom React hooks
│   ├── utils/              # Helper functions
│   └── App.jsx             # Main application
└── package.json
```

### Flow

```
┌─────────────┐
│  React UI   │ ← User uploads workflows & monitors runs
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Go API    │ ← Parses YAML, manages jobs
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Docker    │ ← Executes jobs in isolated containers
└─────────────┘
```

---

## 🧪 Testing

### Run Backend Tests

```bash
cd backend
chmod +x run-tests.sh
./run-tests.sh
```

Or manually:

```bash
go test -v -cover ./...
```

### Test Coverage

Current coverage: **~85%**

- ✅ Parser tests - YAML parsing & validation
- ✅ Storage tests - Memory & MongoDB operations
- ✅ Models tests - Thread-safe operations
- 🚧 Executor tests - Coming soon
- 🚧 API tests - Coming soon

---

## 📚 Documentation

- **[Setup Guide](docs/SETUP.md)** - Detailed installation instructions
- **[Workflow Syntax](docs/WORKFLOWS.md)** - How to write workflows
- **[API Reference](docs/API.md)** - REST API documentation
- **[MongoDB Setup](docs/MONGODB.md)** - MongoDB configuration
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute
- **[Architecture Guide](docs/ARCHITECTURE.md)** - System design & decisions

---

## 🎨 Screenshots

### Dashboard
![Dashboard](docs/images/dashboard.png)

### Workflow Execution
![Workflow Run](docs/images/workflow-run.png)

### Job Details
![Job Details](docs/images/job-details.png)

---

## 🚀 Deployment

### Docker Compose (Recommended)

```bash
docker-compose up -d
```

### Manual Deployment

**Backend:**
```bash
cd backend
go build -o gantry-server ./cmd/server
./gantry-server
```

**Frontend:**
```bash
cd frontend
npm run build
# Serve build/ directory with nginx or similar
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Backend server port |
| `STORAGE_TYPE` | `memory` | `memory` or `mongodb` |
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection string |
| `MONGO_DATABASE` | `gantry` | MongoDB database name |

---

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. 🍴 Fork the repository
2. 🌿 Create a feature branch: `git checkout -b feature/amazing-feature`
3. ✨ Make your changes
4. ✅ Add tests if applicable
5. 📝 Commit: `git commit -m 'Add amazing feature'`
6. 🚀 Push: `git push origin feature/amazing-feature`
7. 🎉 Open a Pull Request

### Development Setup

```bash
# Backend
cd backend
go test ./...              # Run tests
go run ./cmd/server/main.go  # Start dev server

# Frontend
cd frontend
npm start                  # Start dev server
```

### Code Style

- **Go**: Follow standard Go conventions (`gofmt`, `golint`)
- **JavaScript**: Use Prettier (2 spaces)
- **Commits**: Use [Conventional Commits](https://www.conventionalcommits.org/)

---

## 📊 Project Stats

- **Backend**: ~2,000 lines of Go code
- **Frontend**: ~800 lines of React code
- **Test Coverage**: ~85%
- **Docker Images**: Ubuntu, Alpine
- **Dependencies**: Minimal (see go.mod & package.json)

---

## 🙏 Acknowledgments

Inspired by:
- [GitHub Actions](https://github.com/features/actions) - Workflow syntax
- [Drone CI](https://www.drone.io/) - Architecture concepts
- [Jenkins](https://www.jenkins.io/) - Plugin system ideas

Built with:
- [Go](https://go.dev/) - Backend
- [React](https://reactjs.org/) - Frontend
- [Docker](https://www.docker.com/) - Container runtime
- [MongoDB](https://www.mongodb.com/) - Database
- [Gorilla Mux](https://github.com/gorilla/mux/) - HTTP router
- [Lucide React](https://lucide.dev/) - Icons

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 💬 Support

- 💬 **Issues**: [GitHub Issues](https://github.com/sudzp/gantry/issues)
- 💡 **Discussions**: [GitHub Discussions](https://github.com/sudzp/gantry/discussions)


---

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=sudzp/gantry&type=Date)](https://star-history.com/#sudzp/gantry&Date)

---

<div align="center">

**Built with ❤️ for small teams who need simple CI/CD**

[⬆ Back to Top](#-gantry)

</div>