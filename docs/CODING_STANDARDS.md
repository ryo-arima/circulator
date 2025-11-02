# Coding Standards

## Go Code Formatting

### Standard Formatting

- Use `make fmt` or `go fmt` for standard Go formatting
- This follows Go's official formatting guidelines

### Struct Field Alignment

- For improved readability, some files maintain manual struct field alignment
- Files with manual alignment: `pkg/config/logger.go`
- When editing these files, please maintain the alignment pattern
- Comments should start at consistent column positions

### Formatting Commands

```bash
make fmt           # Standard go fmt
make fmt-all       # Standard formatting + import organization
make fmt-keep-align # Skip files with manual alignment
make fmt-check     # Check if files are formatted correctly
```

### Guidelines for logger.go

- Struct fields should be aligned for readability
- After editing, avoid running `go fmt` directly on this file
- Use `make fmt-keep-align` to format other files while preserving alignment
- Review struct alignment during code reviews

### Example of Properly Aligned Struct

```go
type LogEntry struct {
    Timestamp string                 `json:"timestamp"`
    Level     string                 `json:"level"`
    Code      string                 `json:"code"`                    // ログコード（例: PSR-IR）
    Component string                 `json:"component"`               // client, server, agent
    Service   string                 `json:"service"`                 // specific service name
    Message   string                 `json:"message"`
    // ... other fields with consistent alignment
}
```
