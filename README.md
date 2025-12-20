# 🏗️ Gantry

**Lightweight, self-hosted CI/CD platform for small teams**

[![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18+-61DAFB?style=flat&logo=react)](https://reactjs.org/)
[![Docker](https://img.shields.io/badge/Docker-Required-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 🎯 What is Gantry?

Gantry is a lightweight CI/CD platform that brings the power of GitHub Actions to your own infrastructure. Built with Go and React, it's designed for small teams who want:

- ✅ **Simple setup** - Get running in minutes, not hours
- ✅ **Self-hosted** - Keep your code and secrets on your infrastructure
- ✅ **YAML workflows** - Familiar syntax if you've used GitHub Actions
- ✅ **Docker isolation** - Each job runs in a clean container
- ✅ **Real-time monitoring** - Watch your builds live
- ✅ **Lightweight** - Minimal resource footprint

## 🚀 Quick Start

### Prerequisites

- Go 1.19+
- Node.js 16+
- Docker
- 5 minutes of your time

### Installation

```bash
# Clone or create project directory
mkdir gantry && cd gantry

# Backend setup
mkdir backend && cd backend
go mod init gantry
go get github.com/gorilla/mux github.com/docker/docker/client gopkg.in/yaml.v3

# Copy the main.go from artifacts
# Start the server
go run main.go

# Frontend setup (in new terminal)
cd ../
npx create-react-app frontend
cd frontend
npm install lucide-react

# Copy the App.jsx from artifacts to src/App.js
# Start the UI
npm start
```

Visit `http://localhost:3000` and you're ready to go! 🎉

## 📋 Example Workflow

Create a file `build.yml`:

```yaml
name: Build and Test

on:
  push:
    branches:
      - main

jobs:
  test:
    runs-on: ubuntu
    steps:
      - name: Run tests
        run: |
          echo "Running tests..."
          npm test
          
      - name: Build
        run: |
          echo "Building application..."
          npm run build
```

Upload it through the UI, click the play button, and watch it run!

## 🎨 Features

### Current Features

- 📝 **YAML Workflow Definitions** - Write workflows like GitHub Actions
- 🐳 **Docker Isolation** - Each job runs in a fresh container
- ⚡ **Concurrent Execution** - Run multiple workflows simultaneously
- 📊 **Real-time Dashboard** - Monitor builds as they happen
- 📜 **Complete Logs** - See every step's output
- 🔄 **Auto-refresh** - UI updates every 3 seconds
- 🎯 **Manual Triggers** - Start workflows with one click

### Roadmap

- [ ] **Git Webhooks** - Auto-trigger on push/PR
- [ ] **User Authentication** - JWT/OAuth2 support
- [ ] **Secrets Management** - Secure credential storage
- [ ] **Artifacts** - Save build outputs
- [ ] **Matrix Builds** - Test across multiple versions
- [ ] **Notifications** - Email/Slack alerts
- [ ] **Caching** - Speed up repeated builds
- [ ] **Database Backend** - PostgreSQL/MySQL support
- [ ] **Distributed Runners** - Scale across machines
- [ ] **Plugin System** - Extend with custom actions

## 🏗️ Architecture

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

## 📚 Documentation

- [Setup Guide](SETUP.md) - Detailed installation instructions
- [Workflow Syntax](WORKFLOWS.md) - How to write workflows
- [API Reference](API.md) - REST API documentation
- [Contributing](CONTRIBUTING.md) - How to contribute

## 🤝 Contributing

We welcome contributions! Whether it's:

- 🐛 Bug reports
- 💡 Feature requests
- 📖 Documentation improvements
- 🔧 Code contributions

Check out our [Contributing Guide](CONTRIBUTING.md) to get started.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

## 🙏 Acknowledgments

Inspired by:
- [GitHub Actions](https://github.com/features/actions)
- [Drone CI](https://www.drone.io/)
- [Jenkins](https://www.jenkins.io/)
- [CircleCI](https://circleci.com/)

Built with:
- [Go](https://go.dev/) - Backend language
- [React](https://reactjs.org/) - Frontend framework
- [Docker](https://www.docker.com/) - Container runtime
- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router
- [Lucide React](https://lucide.dev/) - Icons

## 📬 Support

- 💬 [Discussions](https://github.com/yourorg/gantry/discussions)
- 🐛 [Issues](https://github.com/yourorg/gantry/issues)
- 📧 Email: support@gantry.dev (if you set up email)

---

<div align="center">
  <strong>Built with ❤️ for small teams who need simple CI/CD</strong>
  <br>
  <sub>Gantry lifts your code from development to deployment</sub>
</div>