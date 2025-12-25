# Gantry CI/CD Platform - AI Coding Guidelines

**Gantry** is a lightweight, self-hosted CI/CD platform built with Go (backend) and React (frontend). It enables teams to define, trigger, and monitor workflows using YAML syntax similar to GitHub Actions.

## Architecture Overview

### Core Components

- **Backend** (`backend/internal/`): Go REST API that orchestrates workflow execution
  - **Parser** (`parser/yaml.go`): Validates & parses GitHub Actions-like YAML workflows, preserving job order
  - **Executor** (`executor/docker.go`): Runs jobs in isolated Docker containers via the Docker daemon
  - **Storage** (`storage/memory.go`, `storage/mongodb.go`): Pluggable persistence (in-memory or MongoDB)
  - **Server** (`server/server.go`): Coordinates all components, exposes REST API
  - **Models** (`models/workflow.go`, `models/run.go`): Domain types (Workflow, WorkflowRun, Job)

- **Frontend** (`frontend/src/`): React 18 UI for uploading workflows and monitoring runs
  - Composition: Header → WorkflowList, RunList, RunDetails (selectable detail views)
  - API layer: `services/apiService.js` abstracts HTTP calls to backend
  - Styling: Tailwind-inspired utility classes in `index.css` (NOT actual Tailwind, bespoke utilities)

### Data Flow

1. User uploads YAML file via frontend → `POST /api/workflows`
2. Backend parses YAML, stores workflow definition, returns name
3. User clicks "Run" on a workflow → `POST /api/workflows/{name}/trigger`
4. Backend creates WorkflowRun, executes jobs sequentially in Docker
5. Real-time UI updates via polling `GET /api/runs` and `GET /api/runs/{id}`
6. Completed logs/status persist in storage layer

## Key Patterns & Conventions

### YAML Workflow Format
- Uses GitHub Actions syntax: `on`, `jobs`, `steps`, `run`, `uses`
- Jobs execute **sequentially** in YAML definition order (preserved via `yaml.v3` Node parsing)
- Each `run:` command executes in a fresh Docker container
- Example: [workflows/example-workflow.yml](../workflows/example-workflow.yml)

### Thread Safety
- **WorkflowRun** uses `sync.RWMutex` for job updates (see [models/run.go](../backend/internal/models/run.go#L29-L48))
- Storage implementations must handle concurrent reads/writes (MongoDB safe by design; memory storage uses mutexes)

### Error Handling
- Backend returns HTTP errors with descriptive messages (e.g., "Failed to parse workflow: invalid YAML")
- Frontend catches API errors and displays user-friendly alerts
- Tests verify error paths explicitly

### Storage Abstraction
- **Storage interface** ([storage/storage.go](../backend/internal/storage/storage.go)) defines: SaveWorkflow, GetWorkflow, SaveRun, GetRun, ListWorkflows, ListRuns
- Switch via `STORAGE_TYPE=mongodb` env var; defaults to in-memory
- MongoDB uses collections: `workflows`, `runs`

## Development Workflows

### Running Backend
```bash
cd backend
go mod download
go run ./cmd/server/main.go
# Server runs on :8080, listens for HTTP requests
```

### Running Frontend
```bash
cd frontend
npm install
npm start
# React dev server on :3000, proxies `/api/*` to backend
```

### Testing Backend
```bash
cd backend
./run-tests.sh          # Runs all tests + generates coverage.html
go test -v ./...        # Quick test run
go test -race ./...     # Detect data races
go test -v -run TestName ./package  # Specific test
```

- Tests mirror source structure (e.g., `yaml_test.go` tests `yaml.go`)
- Mock storage in executor tests; verify Docker calls in isolation
- Target: ~85% coverage

### Using Docker Compose
```bash
docker-compose up  # Starts MongoDB, backend, frontend
# MongoDB: :27017, Backend: :8080, Frontend: :3000
```

## Common Tasks

### Adding a New API Endpoint
1. Define handler in [api/handlers.go](../backend/internal/api/handlers.go)
2. Register route in [api/routes.go](../backend/internal/api/routes.go) via gorilla/mux
3. Call server methods (e.g., `h.server.TriggerWorkflow()`)
4. Return JSON with `json.NewEncoder(w).Encode(data)`

### Extending Workflow Features
1. Update YAML model in [models/workflow.go](../backend/internal/models/workflow.go)
2. Update parser in [parser/yaml.go](../backend/internal/parser/yaml.go) to handle new fields
3. Add parser tests validating new field in [parser/yaml_test.go](../backend/internal/parser/yaml_test.go)
4. Update executor or models as needed

### Adding Frontend Component
- Create in `src/components/` following pattern: `<ComponentName>.jsx`
- Use `apiService` for backend calls; follow the fetch/error pattern in App.jsx
- Style with existing utility classes from `index.css`
- Import lucide-react icons for consistent UI

## Critical Dependencies & Integrations

- **Docker Daemon**: Backend requires access to `/var/run/docker.sock` to spawn job containers
- **Go Packages**: `gopkg.in/yaml.v3` (YAML parsing), `github.com/gorilla/mux` (routing), `go.mongodb.org/mongo-driver` (optional MongoDB)
- **Node Packages**: `react`, `react-dom`, `lucide-react` (icons)
- **API Contract**: Backend serves `application/json`; frontend sends YAML as `text/yaml` on upload

## Testing Philosophy

- **Unit-focused**: Test components in isolation (parser ≠ executor ≠ storage)
- **No external deps**: Use in-memory storage for backend tests; mock API calls in frontend
- **Error cases**: Explicitly test malformed inputs, missing fields, Docker failures
- **Integration tests** (less common): docker-compose can validate end-to-end flow manually

## When Debugging

- **Backend**: Check logs for parsing errors, Docker connection issues, storage errors. Enable verbose logging via `-v` test flag.
- **Frontend**: Check browser console for API call failures, undefined state. Use React DevTools.
- **Data Races**: `go test -race ./...` catches concurrent access bugs in Server/WorkflowRun.
- **Workflow Execution**: Inspect `RunDetails` component for job logs; check Docker daemon connectivity.

---

**Last Updated**: 2025-12-25  
**Related Docs**: [API.md](../docs/API.md), [TESTING.md](../docs/TESTING.md), [WORKFLOWS.md](../docs/WORKFLOWS.md), [SETUP.md](../docs/SETUP.md)
