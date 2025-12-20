# Gantry Workflow Syntax

## Basic Structure
```yaml
name: Workflow Name

on:
  push:
    branches:
      - main
      - develop

jobs:
  job-name:
    runs-on: ubuntu  # or alpine
    steps:
      - name: Step name
        run: |
          # Your commands here
```

## Supported Fields

### name (required)
The name of your workflow

### on (required)
Trigger configuration (currently only supports `push`)

### jobs (required)
Map of jobs to execute

#### runs-on
Container image to use:
- `ubuntu` - Uses ubuntu:latest
- `alpine` - Uses alpine:latest

#### steps
Array of steps to execute

Each step has:
- `name` - Display name
- `run` - Shell commands to execute

## Examples

### Simple Build
```yaml
name: Simple Build

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu
    steps:
      - name: Build
        run: echo "Building..."
```

### Multi-Job Workflow
```yaml
name: Test and Deploy

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu
    steps:
      - name: Run tests
        run: npm test

  deploy:
    runs-on: alpine
    steps:
      - name: Deploy
        run: echo "Deploying..."
```
