# Gantry - Lightweight CI/CD Platform

![Gantry Logo](https://img.shields.io/badge/Gantry-CI%2FCD-blue?style=for-the-badge)

A lightweight, self-hosted CI/CD platform built with Go and React, inspired by GitHub Actions.

**Gantry** lifts your code from development to deployment with ease.

## Prerequisites

- **Go** 1.21+ installed
- **Node.js** 18+ and npm
- **Docker** installed and running
- Git (optional)

## Project Structure

```
gantry/
├── backend/
│   ├── main.go
│   └── go.mod
├── frontend/
│   ├── src/
│   │   └── App.jsx
│   ├── package.json
│   └── index.html
└── workflows/
    └── example-workflow.yml
```

## Backend Setup (Go)

### 1. Create the project directory

```bash
mkdir -p gantry/backend
cd gantry/backend
```

### 2. Initialize Go module

```bash
go mod init gantry
```

### 3. Install dependencies

```bash
go get github.com/gorilla/mux
go get github.com/docker/docker/client
go get gopkg.in/yaml.v3
```

### 4. Create main.go

Copy the backend code from the artifact into `main.go`

### 5. Run the backend

```bash
# Make sure Docker is running
go run main.go
```

The API server will start on `http://localhost:8080`

## Frontend Setup (React)

### 1. Create React app

```bash
cd ../
npx create-react-app frontend
cd frontend
```

### 2. Install dependencies

```bash
npm install lucide-react
```

### 3. Replace src/App.js

Copy the React component from the artifact and replace the content of `src/App.js`

### 4. Start development server

```bash
npm start
```

The UI will open at `http://localhost:3000`

## Testing

### Backend Tests

```bash
cd backend
go test -v -race ./...                 # All tests
go test -v -race -tags=integration ./internal/server/...  # Integration tests only
```

### Frontend Tests

```bash
cd frontend
npm test -- --coverage --watchAll=false
npm test -- App.test.integration.js --coverage  # Integration tests
```

See [Testing Guide](TESTING.md) for more details.

## API Endpoints

### Workflows

- **POST** `/api/workflows` - Upload a workflow YAML file
- **GET** `/api/workflows` - List all workflows
- **POST** `/api/workflows/{name}/trigger` - Trigger a workflow execution

### Runs

- **GET** `/api/runs` - List all workflow runs
- **GET** `/api/runs/{id}` - Get details of a specific run

## Workflow YAML Format

```yaml
name: My Workflow

on:
  push:
    branches:
      - main

jobs:
  job-name:
    runs-on: ubuntu  # or alpine
    steps:
      - name: Step 1
        run: |
          echo "Hello World"
          
      - name: Step 2
        run: |
          # Your commands here
```

## Usage

### 1. Upload a Workflow

Create a YAML workflow file (see example-workflow.yml) and upload it through the UI:
- Click "Upload Workflow YAML"
- Select your `.yml` or `.yaml` file
- The workflow will appear in the workflows list

### 2. Trigger a Workflow

- Click the play button next to a workflow
- The workflow will execute in a Docker container
- View real-time status in the "Recent Runs" section

### 3. View Run Details

- Click on any run in the "Recent Runs" list
- See job status, steps, and output logs

## Features

✅ YAML-based workflow definitions
✅ Docker container isolation
✅ Real-time status updates
✅ Job and step execution tracking
✅ Output logs capture
✅ Multiple concurrent runs
✅ REST API

## Future Enhancements

Consider adding:

1. **Authentication & Authorization** - User management and permissions
2. **Webhook Integration** - Trigger workflows from Git events
3. **Secrets Management** - Secure storage for API keys and tokens
4. **Artifact Storage** - Save build outputs
5. **Matrix Builds** - Run jobs across multiple configurations
6. **Conditional Execution** - Run steps based on conditions
7. **Caching** - Cache dependencies between runs
8. **Email Notifications** - Alert on workflow completion
9. **Database** - Replace in-memory storage with PostgreSQL/MySQL
10. **Distributed Runners** - Scale execution across multiple machines

## Troubleshooting

**Docker connection issues:**
```bash
# Check if Docker daemon is running
docker ps

# On Linux, add user to docker group
sudo usermod -aG docker $USER
```

**CORS errors:**
- Ensure the backend CORS middleware allows your frontend origin
- Check that API_URL in frontend matches your backend URL

**Port already in use:**
```bash
# Change port in backend
export PORT=8081
go run main.go

# Update API_URL in frontend accordingly
```

## Production Deployment

For production use:

1. **Build the Go binary:**
   ```bash
   go build -o cicd-server main.go
   ```

2. **Build React for production:**
   ```bash
   npm run build
   ```

3. **Use a reverse proxy** (Nginx/Caddy) to serve frontend and proxy API requests

4. **Add TLS/SSL** certificates

5. **Set up proper logging and monitoring**

6. **Use a production database** (PostgreSQL recommended)

7. **Implement authentication** (JWT tokens, OAuth2)

## License

MIT License - Feel free to modify and use for your team!

## Support

For issues or questions, refer to the inline code comments or extend the functionality as needed for your use case.