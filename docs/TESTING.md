# Testing Guide

## 🧪 Overview

Gantry has comprehensive unit and integration tests for both backend and frontend with **60%+ coverage enforced**. This guide covers running tests, writing new tests, and best practices.

**Current Testing:**
- ✅ Backend: Unit tests + 6 integration tests (server operations)
- ✅ Frontend: Integration tests (8 test cases for React components)
- ✅ Security: Gosec (Go) + npm audit (JavaScript)
- ✅ Coverage: Minimum 60% enforced in CI/CD

## Running Tests

### Quick Test

```bash
cd backend
go test ./...
```

### With Coverage

```bash
cd backend
./run-tests.sh
```

This will:
- Run all tests
- Generate coverage report
- Create HTML coverage visualization
- Show test summary

### Run Specific Package

```bash
# Test parser only
go test -v ./internal/parser

# Test storage only
go test -v ./internal/storage

# Test models only
go test -v ./internal/models
```

### Run Specific Test

```bash
go test -v -run TestParse_ValidWorkflow ./internal/parser
```

### With Race Detection

```bash
go test -race ./...
```

### Verbose Output

```bash
go test -v ./...
```

## Frontend Testing

### Run Frontend Tests

```bash
cd frontend
npm test -- --coverage --watchAll=false
```

### Run Integration Tests Only

```bash
npm test -- App.test.integration.js --coverage
```

### Watch Mode (Development)

```bash
npm test -- --watch
```

## Test Structure

### Server Integration Tests (`internal/server/server_test.go`)

Tests core server functionality:

- ✅ Parse and save workflows
- ✅ List workflows
- ✅ Trigger workflow execution
- ✅ Get workflow statistics
- ✅ Get workflow runs (filtering)
- ✅ Delete workflows (with cascade delete)

**Key Feature:** Cascade delete verification ensures runs are deleted when workflow is removed

### Parser Tests (`internal/parser/yaml_test.go`)

Tests YAML parsing and validation:

- ✅ Valid workflow parsing
- ✅ Invalid YAML handling
- ✅ Workflow validation
- ✅ Missing required fields
- ✅ Job order preservation

**Example:**
```go
func TestParse_ValidWorkflow(t *testing.T) {
    yaml := `
name: Test Workflow
on:
  push:
    branches: [main]
jobs:
  test:
    runs-on: ubuntu
    steps:
      - name: Test
        run: echo "test"
`
    p := NewParser()
    wf, err := p.Parse([]byte(yaml))
    
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }
    // More assertions...
}
```

### Storage Tests (`internal/storage/memory_test.go`)

Tests in-memory storage operations:

- ✅ Save/Get workflows
- ✅ List workflows
- ✅ Delete workflows
- ✅ Save/Get runs
- ✅ Update runs
- ✅ Concurrent operations

**Example:**
```go
func TestMemoryStorage_SaveAndGetWorkflow(t *testing.T) {
    store := NewMemoryStorage()
    wf := &models.Workflow{Name: "Test"}
    
    store.SaveWorkflow(wf)
    retrieved, err := store.GetWorkflow("Test")
    
    if err != nil {
        t.Fatalf("Failed to get workflow: %v", err)
    }
    // More assertions...
}
```

### Models Tests (`internal/models/run_test.go`)

Tests thread-safe operations on WorkflowRun:

- ✅ Update job
- ✅ Get job
- ✅ Set status
- ✅ Complete run
- ✅ Clone run
- ✅ Thread safety

**Example:**
```go
func TestWorkflowRun_ThreadSafety(t *testing.T) {
    run := &WorkflowRun{ID: "run-1", Jobs: make(map[string]Job)}
    var wg sync.WaitGroup
    
    // Concurrent writes
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            job := Job{Status: "running"}
            run.UpdateJob(fmt.Sprintf("job-%d", id), job)
        }(i)
    }
    
    wg.Wait()
    // Assertions...
}
```

## Writing New Tests

### Test File Naming

- Place test files next to the code they test
- Name them `*_test.go`
- Example: `parser.go` → `parser_test.go`

### Test Function Naming

Use descriptive names with the pattern:
```
Test<FunctionName>_<Scenario>
```

Examples:
- `TestParse_ValidWorkflow`
- `TestParse_InvalidYAML`
- `TestValidate_MissingName`

### Test Structure

```go
func TestSomething_Scenario(t *testing.T) {
    // 1. Setup
    input := "test data"
    
    // 2. Execute
    result, err := FunctionUnderTest(input)
    
    // 3. Assert
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### Table-Driven Tests

For testing multiple scenarios:

```go
func TestValidate(t *testing.T) {
    tests := []struct {
        name    string
        input   *Workflow
        wantErr bool
    }{
        {
            name:    "valid workflow",
            input:   &Workflow{Name: "Test", Jobs: ...},
            wantErr: false,
        },
        {
            name:    "missing name",
            input:   &Workflow{Name: "", Jobs: ...},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Testing Best Practices

### 1. Test One Thing

Each test should verify one specific behavior:

❌ Bad:
```go
func TestEverything(t *testing.T) {
    // Tests 10 different things
}
```

✅ Good:
```go
func TestParse_ValidWorkflow(t *testing.T) { ... }
func TestParse_InvalidYAML(t *testing.T) { ... }
func TestValidate_MissingName(t *testing.T) { ... }
```

### 2. Use Descriptive Names

❌ Bad:
```go
func TestParse(t *testing.T) { ... }
```

✅ Good:
```go
func TestParse_ValidWorkflow(t *testing.T) { ... }
```

### 3. Test Error Cases

Always test both success and failure paths:

```go
func TestGetWorkflow_Exists(t *testing.T) { ... }
func TestGetWorkflow_NotFound(t *testing.T) { ... }
```

### 4. Clean Up Resources

```go
func TestWithFile(t *testing.T) {
    f, _ := os.CreateTemp("", "test")
    defer os.Remove(f.Name())  // Cleanup
    
    // Test code...
}
```

### 5. Use t.Helper() for Helper Functions

```go
func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }
}
```

## Mocking

For testing components that depend on external services:

```go
// Mock storage
type MockStorage struct {
    workflows map[string]*Workflow
}

func (m *MockStorage) GetWorkflow(name string) (*Workflow, error) {
    wf, exists := m.workflows[name]
    if !exists {
        return nil, errors.New("not found")
    }
    return wf, nil
}
```

## Coverage Goals

Current enforced minimum: **60%+**

**CI/CD Enforcement:**
- Coverage must be ≥ 60% to pass builds
- Coverage report uploaded to Codecov
- Fails build if threshold not met

**Target coverage by module:**
- Models: **90%+** ✅
- Parser: **90%+** ✅
- Storage: **85%+** ✅
- Server: **75%+** (including integration tests) ✅
- Frontend App: **70%+** (including integration tests) ✅
- Executor: **60%+** 🚧
- API: **70%+** 🚧

## CI/CD Integration

Tests run automatically on:
- Every push to main/develop
- Every pull request
- Before merging

### GitHub Actions Pipeline (10 Jobs)

**Backend:**
1. `backend-lint` - Go formatting, vet, golangci-lint
2. `backend-test` - Unit + integration tests, coverage check (≥60%)
3. `backend-security` - Gosec vulnerability scanning
4. `backend-build` - Go binary compilation

**Frontend:**
5. `frontend-lint` - ESLint, Prettier
6. `frontend-test` - Unit + integration tests, coverage
7. `frontend-security` - npm audit (moderate level)
8. `frontend-build` - Production build

**Gate:**
9. `all-checks-passed` - Requires all 8 jobs to pass before merge

**Test Commands Run:**
```bash
# Backend
go test -v -race ./...
go test -v -race -tags=integration ./internal/server/...

# Frontend
npm test -- --coverage --watchAll=false
npm test -- App.test.integration.js --coverage
```

## Benchmarking

For performance-critical code:

```go
func BenchmarkParse(b *testing.B) {
    yaml := []byte(`name: Test...`)
    p := NewParser()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        p.Parse(yaml)
    }
}
```

Run benchmarks:
```bash
go test -bench=. ./...
```

## Test Data

Store test fixtures in `testdata/` directories:

```
internal/parser/
├── parser.go
├── parser_test.go
└── testdata/
    ├── valid-workflow.yml
    ├── invalid-workflow.yml
    └── complex-workflow.yml
```

Load test data:
```go
data, _ := os.ReadFile("testdata/valid-workflow.yml")
```

## Troubleshooting

### Tests Pass Locally But Fail in CI

- Check for race conditions: `go test -race`
- Verify no absolute paths
- Check for timing dependencies
- Ensure clean state between tests

### Flaky Tests

- Add retries for external dependencies
- Use mocks instead of real services
- Increase timeouts if needed
- Check for goroutine leaks

### Slow Tests

- Use `t.Parallel()` for independent tests
- Mock slow operations
- Run expensive tests separately
- Profile with `go test -cpuprofile`

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Test Fixtures](https://dave.cheney.net/2016/05/10/test-fixtures-in-go)
- [Mocking](https://github.com/golang/mock)

## Contributing Tests

When contributing:

1. **Add tests for new features**
2. **Maintain coverage** - Don't decrease coverage
3. **Test error cases** - Not just happy path
4. **Use meaningful names** - Describe what's being tested
5. **Keep tests fast** - Mock external dependencies
6. **Run locally first** - `go test -race ./...` and `npm test`

---

**Questions?** Open an issue or reach out to the maintainers!