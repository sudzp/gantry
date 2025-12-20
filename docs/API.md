# Gantry API Reference

## Base URL
http://localhost:8080/api

## Endpoints

### Workflows

#### Upload Workflow
POST /api/workflows
Content-Type: text/yaml
[YAML workflow content]

**Response:**
```json
{
  "message": "Workflow uploaded successfully",
  "name": "Build and Test"
}
```

#### List Workflows
GET /api/workflows

**Response:**
```json
[
  {
    "name": "Build and Test",
    "on": {
      "push": {
        "branches": ["main"]
      }
    },
    "jobs": {...}
  }
]
```

#### Trigger Workflow
POST /api/workflows/{name}/trigger

**Response:**
```json
{
  "id": "run-1234567890",
  "workflow_name": "Build and Test",
  "status": "running",
  "started_at": "2025-01-15T10:30:00Z"
}
```

### Runs

#### List Runs
GET /api/runs

**Response:**
```json
[
  {
    "id": "run-1234567890",
    "workflow_name": "Build and Test",
    "status": "success",
    "started_at": "2025-01-15T10:30:00Z",
    "completed_at": "2025-01-15T10:35:00Z"
  }
]
```

#### Get Run Details
GET /api/runs/{id}

**Response:**
```json
{
  "id": "run-1234567890",
  "workflow_name": "Build and Test",
  "status": "success",
  "jobs": {
    "build": {
      "runs_on": "ubuntu",
      "status": "success",
      "output": "Build logs here...",
      "steps": [...]
    }
  }
}
```
