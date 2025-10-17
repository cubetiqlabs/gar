# Contributing to GoArchive (gar)

Thank you for your interest in contributing to GoArchive (gar)! This document provides comprehensive guidelines for both human developers and AI agents to contribute effectively to this project.

## 📋 Table of Contents

-   [Code of Conduct](#code-of-conduct)
-   [Getting Started](#getting-started)
-   [For Human Developers](#for-human-developers)
-   [For AI Agents](#for-ai-agents)
-   [Development Workflow](#development-workflow)
-   [Coding Standards](#coding-standards)
-   [Testing Guidelines](#testing-guidelines)
-   [Commit Guidelines](#commit-guidelines)
-   [Pull Request Process](#pull-request-process)
-   [Project Architecture](#project-architecture)
-   [Security Considerations](#security-considerations)
-   [Documentation](#documentation)
-   [Community](#community)

---

## Code of Conduct

### Our Pledge

We as members, contributors, and leaders pledge to make participation in GoArchive a harassment-free experience for everyone, regardless of age, body size, visible or invisible disability, ethnicity, sex characteristics, gender identity and expression, level of experience, education, socio-economic status, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Our Standards

**Positive behaviors include:**

-   Using welcoming and inclusive language
-   Being respectful of differing viewpoints and experiences
-   Gracefully accepting constructive criticism
-   Focusing on what is best for the community
-   Showing empathy towards other community members

**Unacceptable behaviors include:**

-   Trolling, insulting/derogatory comments, and personal attacks
-   Public or private harassment
-   Publishing others' private information without permission
-   Other conduct which could reasonably be considered inappropriate

### Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be reported to the project team. All complaints will be reviewed and investigated promptly and fairly.

---

## Getting Started

### Prerequisites

Before contributing, ensure you have:

-   **Go 1.21 or higher** installed
-   **Git** for version control
-   **A GitHub account** for pull requests
-   **Basic understanding** of Go programming
-   **Familiarity** with archive formats (ZIP, TAR, GZIP)

### First-Time Setup

```bash
# 1. Fork the repository on GitHub
# Click the "Fork" button at https://github.com/cubetiqlabs/gar

# 2. Clone your fork
git clone https://github.com/cubetiqlabs/gar.git
cd gar

# 3. Add upstream remote
git remote add upstream https://github.com/cubetiqlabs/gar.git

# 4. Install dependencies
go mod download

# 5. Verify everything works
go test ./...
go build -o gar main.go
./gar -version

# 6. Create a branch for your work
git checkout -b feature/my-feature
```

### Finding Ways to Contribute

1. **Check existing issues** labeled with:

    - `good first issue` - Great for newcomers
    - `help wanted` - We need assistance here
    - `bug` - Something isn't working
    - `enhancement` - New feature or improvement

2. **Propose new features**: Open an issue to discuss before implementing

3. **Improve documentation**: Always welcome and valuable

4. **Fix bugs**: Check the issue tracker for reported bugs

---

## For Human Developers

### What We're Looking For

-   **Bug fixes**: Corrections to existing functionality
-   **Features**: New capabilities (discuss in issue first)
-   **Performance improvements**: Optimizations with benchmarks
-   **Documentation**: Clarifications and additions
-   **Tests**: Increase coverage and edge cases
-   **Security**: Vulnerability fixes (report privately first)

### Development Environment

#### Recommended IDE Setup

**Visual Studio Code:**

```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "package",
    "editor.formatOnSave": true,
    "go.formatTool": "gofmt"
}
```

**GoLand/IntelliJ IDEA:**

-   Enable Go modules support
-   Configure gofmt as formatter
-   Enable golangci-lint inspections

#### Useful Tools

```bash
# Install golangci-lint (comprehensive linter)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install staticcheck (static analysis)
go install honnef.co/go/tools/cmd/staticcheck@latest

# Install goreleaser (for releases)
go install github.com/goreleaser/goreleaser@latest
```

### Making Changes

1. **Write clear, idiomatic Go code**
2. **Add tests for new functionality**
3. **Update documentation** as needed
4. **Run all tests** before submitting
5. **Keep commits atomic** and well-described

### Testing Your Changes

```bash
# Run all tests
go test ./... -v

# Run tests with race detector
go test ./... -race

# Run tests with coverage
go test ./... -cover -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out

# Benchmark tests
go test -bench=. -benchmem

# Test specific package
go test ./pkg/compression -v
```

---

## For AI Agents

### Context About This Project

**Project Name**: GoArchive  
**Language**: Go (Golang)  
**Purpose**: Cross-platform archive manager (WinRAR alternative)  
**Key Features**: Compression, extraction, encryption, parallel processing  
**Architecture**: CLI application with worker pool pattern

### Core Design Principles

1. **Performance First**: Use concurrency, efficient buffering, minimize allocations
2. **Security by Default**: Path traversal prevention, strong encryption (AES-256-GCM)
3. **Cross-Platform**: Must work on Windows, Linux, macOS without platform-specific code
4. **Single Binary**: No external dependencies beyond Go standard library
5. **Simple UX**: Clear error messages, sensible defaults, verbose mode for debugging

### Code Context

#### Current File Structure

```
gar/
├── main.go                   # Entry point, CLI parsing, orchestration
│   ├── compress()            # Creates archives
│   ├── extract()             # Extracts archives with worker pool
│   ├── listArchive()         # Lists contents
│   ├── encryption functions  # AES-256-GCM implementation
│   └── format handlers       # ZIP and TAR.GZ support
└── go.mod                    # Go module definition
```

#### Key Functions and Their Purposes

```go
// main() - Entry point
// - Parses CLI flags
// - Validates input
// - Dispatches to compress/extract/list

// compress(input, output, opts) - Creates archives
// - Walks directory tree
// - Handles single files or directories
// - Applies compression level
// - Adds encryption if password provided

// extract(input, output, opts) - Extracts archives
// - Uses worker pool for parallelism
// - Prevents path traversal attacks
// - Handles encrypted archives
// - Preserves file permissions

// compressZip() - ZIP format compression
// - Uses archive/zip package
// - Configurable compression levels
// - Handles nested directories

// extractZip() - ZIP format extraction
// - Parallel file extraction
// - Security validation
// - Permission preservation

// compressTarGz() - TAR.GZ compression
// - Uses archive/tar + compress/gzip
// - Sequential archiving
// - Efficient for large files

// extractTarGz() - TAR.GZ extraction
// - Sequential processing
// - Type-safe header handling
// - Directory creation

// Encryption functions:
// - newEncryptedWriter() - Wraps writer with AES-256-GCM
// - newEncryptedReader() - Wraps reader with decryption
// - Uses PBKDF2 for key derivation (100k iterations)
```

#### Important Data Structures

```go
// ArchiveOptions - Configuration for operations
type ArchiveOptions struct {
    Format           ArchiveFormat      // ZIP or TAR.GZ
    CompressionLevel CompressionLevel   // Fastest, Normal, Best
    Password         string             // For encryption
    Workers          int                // Parallel worker count
    Verbose          bool               // Debug output
}

// ArchiveFormat - Supported archive types
type ArchiveFormat int
const (
    FormatZip ArchiveFormat = iota
    FormatTarGz
)

// CompressionLevel - Compression settings
type CompressionLevel int
const (
    LevelFastest CompressionLevel = iota
    LevelNormal
    LevelBest
)
```

#### Critical Security Patterns

```go
// ALWAYS validate paths to prevent traversal attacks
destPath := filepath.Join(outputPath, fileName)
if !strings.HasPrefix(filepath.Clean(destPath), filepath.Clean(outputPath)) {
    return fmt.Errorf("illegal file path: %s", fileName)
}

// ALWAYS use crypto/rand, never math/rand for security
salt := make([]byte, 32)
if _, err := rand.Read(salt); err != nil {
    return err
}

// ALWAYS use PBKDF2 for password-based keys
key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
```

### AI Agent Contribution Guidelines

When contributing as an AI agent, follow these specific guidelines:

#### 1. Code Generation Standards

**DO:**

-   ✅ Generate complete, working functions (no placeholders)
-   ✅ Include proper error handling for all operations
-   ✅ Add descriptive comments for complex logic
-   ✅ Use Go idioms (defer, error returns, etc.)
-   ✅ Consider edge cases (empty files, large files, permissions)
-   ✅ Follow existing code style and patterns
-   ✅ Include example usage in comments

**DON'T:**

-   ❌ Generate pseudo-code or incomplete functions
-   ❌ Use `TODO` or `FIXME` comments without implementation
-   ❌ Ignore errors or use `panic()` without good reason
-   ❌ Break existing API contracts
-   ❌ Add external dependencies without discussion
-   ❌ Use deprecated Go features
-   ❌ Generate code without proper context

#### 2. Security Requirements

**CRITICAL - Always implement:**

```go
// Path traversal prevention (REQUIRED for all file operations)
func validatePath(basePath, targetPath string) error {
    cleanBase := filepath.Clean(basePath)
    cleanTarget := filepath.Clean(targetPath)

    if !strings.HasPrefix(cleanTarget, cleanBase) {
        return fmt.Errorf("path traversal detected: %s", targetPath)
    }
    return nil
}

// Proper error handling (REQUIRED)
if err != nil {
    return fmt.Errorf("descriptive context: %w", err)
}

// Resource cleanup (REQUIRED)
file, err := os.Open(path)
if err != nil {
    return err
}
defer file.Close() // ALWAYS defer close
```

**Security Checklist:**

-   [ ] All file paths validated for traversal attacks
-   [ ] All resources properly closed (defer)
-   [ ] All errors properly wrapped and returned
-   [ ] No sensitive data in logs
-   [ ] Crypto operations use crypto/rand
-   [ ] Password handling follows best practices

#### 3. Testing Requirements

When adding functionality, ALWAYS include tests:

```go
// Example test structure
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test.txt",
            want:    "expected output",
            wantErr: false,
        },
        {
            name:    "empty input",
            input:   "",
            want:    "",
            wantErr: true,
        },
        {
            name:    "path traversal attempt",
            input:   "../etc/passwd",
            want:    "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewFeature(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFeature() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("NewFeature() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

#### 4. Documentation Requirements

**For every new function, provide:**

```go
// FunctionName does X and returns Y.
// It handles edge cases A, B, and C.
//
// Parameters:
//   - param1: description of param1
//   - param2: description of param2
//
// Returns:
//   - result: description of return value
//   - error: description of error conditions
//
// Example:
//   result, err := FunctionName("input")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(result)
//
// Security: Validates all file paths to prevent traversal attacks.
// Performance: Uses buffered I/O for files larger than 1MB.
func FunctionName(param1, param2 string) (string, error) {
    // Implementation
}
```

#### 5. Performance Guidelines

**Optimization Priorities:**

1. **Correctness First**: Code must work correctly before optimization
2. **Measure Before Optimizing**: Use benchmarks to identify bottlenecks
3. **Common Patterns**:
    - Use buffered I/O (32KB+ buffers)
    - Avoid allocations in hot paths
    - Reuse buffers with sync.Pool for high-frequency operations
    - Use goroutines for I/O-bound operations
    - Limit goroutines with worker pools (don't spawn unbounded)

**Benchmark Template:**

```go
func BenchmarkFeature(b *testing.B) {
    // Setup
    input := setupTestData()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Code to benchmark
        _ = ProcessFeature(input)
    }
}

func BenchmarkFeatureParallel(b *testing.B) {
    input := setupTestData()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _ = ProcessFeature(input)
        }
    })
}
```

#### 6. Common Patterns in This Codebase

**Worker Pool Pattern:**

```go
var wg sync.WaitGroup
sem := make(chan struct{}, workerCount)

for _, item := range items {
    wg.Add(1)
    sem <- struct{}{} // Acquire

    go func(i Item) {
        defer wg.Done()
        defer func() { <-sem }() // Release

        // Process item
        processItem(i)
    }(item)
}

wg.Wait()
```

**Error Wrapping:**

```go
if err := operation(); err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

**Resource Management:**

```go
file, err := os.Open(path)
if err != nil {
    return fmt.Errorf("open file: %w", err)
}
defer file.Close()

// Use file...
```

#### 7. Validation Checklist for AI-Generated Code

Before submitting code, verify:

-   [ ] **Compiles**: `go build ./...` succeeds
-   [ ] **Tests pass**: `go test ./...` succeeds
-   [ ] **Formatted**: `gofmt -w .` applied
-   [ ] **Linted**: `golangci-lint run` passes
-   [ ] **Race-free**: `go test -race ./...` succeeds
-   [ ] **Secure**: No path traversal, proper error handling
-   [ ] **Documented**: All exported functions have godoc comments
-   [ ] **Tested**: New code has unit tests
-   [ ] **Benchmarked**: Performance-critical code has benchmarks

#### 8. Context Preservation

When working on this project, remember:

**Current State:**

-   Single file architecture (main.go)
-   Supports ZIP and TAR.GZ formats
-   AES-256-GCM encryption implemented
-   Worker pool for parallel extraction
-   CLI-only interface (no GUI)

**Future Roadmap:**

-   7-Zip format support
-   Archive splitting/merging
-   GUI application
-   Progress indicators
-   Configuration file support

**Constraints:**

-   Must remain single-binary
-   No external dependencies (only Go stdlib + x/crypto)
-   Must work on Windows, Linux, macOS
-   CLI interface must remain backward compatible

---

## Development Workflow

### Branch Naming Convention

```
feature/short-description    # New features
bugfix/issue-number         # Bug fixes
hotfix/critical-issue       # Urgent fixes
docs/what-changed           # Documentation only
refactor/component-name     # Code refactoring
test/what-testing           # Test additions
```

### Development Process

```bash
# 1. Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# 2. Create feature branch
git checkout -b feature/my-feature

# 3. Make changes
# ... edit files ...

# 4. Test changes
go test ./... -race -cover

# 5. Format and lint
gofmt -w .
golangci-lint run

# 6. Commit changes
git add .
git commit -m "feat: add new feature"

# 7. Push to your fork
git push origin feature/my-feature

# 8. Create Pull Request on GitHub
```

---

## Coding Standards

### Go Style Guide

Follow the official [Effective Go](https://golang.org/doc/effective_go.html) guidelines.

#### Naming Conventions

```go
// Good
type ArchiveOptions struct { }
func compressZip() { }
const BufferSize = 32 * 1024

// Bad
type archiveOptions struct { }  // Unexported when should be exported
func CompressZIP() { }           // Inconsistent acronym casing
const bufferSize = 32 * 1024     // Should be exported constant
```

#### Error Handling

```go
// Good - descriptive error wrapping
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad - swallowing errors
_ = doSomething()

// Bad - non-descriptive errors
if err := doSomething(); err != nil {
    return err
}
```

#### Comments

```go
// Good - explains WHY and WHAT, not HOW
// validatePath prevents directory traversal attacks by ensuring
// the target path is within the base directory.
func validatePath(base, target string) error {
    // Implementation...
}

// Bad - states the obvious
// This function validates the path
func validatePath(base, target string) error {
    // Implementation...
}
```

#### Function Length

-   Keep functions under 50 lines when possible
-   Extract complex logic into helper functions
-   One function should do one thing well

#### Package Organization

```go
// Current (single file)
main.go

// Future (when file grows beyond 2000 lines)
main.go              // Entry point, CLI
compress.go          // Compression functions
extract.go           // Extraction functions
encrypt.go           // Encryption/decryption
formats/
  ├── zip.go         // ZIP-specific code
  └── targz.go       // TAR.GZ-specific code
```

---

## Testing Guidelines

### Test Coverage Requirements

-   **Minimum coverage**: 70% for new code
-   **Critical paths**: 90%+ coverage (encryption, path validation)
-   **Edge cases**: Must be tested explicitly

### Test Structure

```go
func TestFunction(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"

    // Act
    result, err := Function(input)

    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

### Table-Driven Tests

```go
func TestValidatePath(t *testing.T) {
    tests := []struct {
        name    string
        base    string
        target  string
        wantErr bool
    }{
        {"valid path", "/tmp", "/tmp/file.txt", false},
        {"traversal attempt", "/tmp", "/etc/passwd", true},
        {"relative traversal", "/tmp", "/tmp/../etc/passwd", true},
        {"windows traversal", "C:\\temp", "C:\\Windows\\System32", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validatePath(tt.base, tt.target)
            if (err != nil) != tt.wantErr {
                t.Errorf("validatePath() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Tests

```go
func TestEndToEnd(t *testing.T) {
    // Create temp directory
    tmpDir := t.TempDir()

    // Setup test data
    testFile := filepath.Join(tmpDir, "test.txt")
    if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
        t.Fatal(err)
    }

    // Compress
    archivePath := filepath.Join(tmpDir, "test.zip")
    opts := ArchiveOptions{Format: FormatZip}
    if err := compress(tmpDir, archivePath, opts); err != nil {
        t.Fatalf("compress failed: %v", err)
    }

    // Extract
    extractDir := filepath.Join(tmpDir, "extracted")
    if err := extract(archivePath, extractDir, opts); err != nil {
        t.Fatalf("extract failed: %v", err)
    }

    // Verify
    extractedFile := filepath.Join(extractDir, "test.txt")
    content, err := os.ReadFile(extractedFile)
    if err != nil {
        t.Fatalf("read extracted file: %v", err)
    }
    if string(content) != "test" {
        t.Errorf("content mismatch: got %q, want %q", content, "test")
    }
}
```

---

## Commit Guidelines

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

#### Types

-   `feat`: New feature
-   `fix`: Bug fix
-   `docs`: Documentation changes
-   `style`: Code style changes (formatting, missing semicolons, etc.)
-   `refactor`: Code refactoring
-   `test`: Adding or updating tests
-   `chore`: Maintenance tasks
-   `perf`: Performance improvements
-   `ci`: CI/CD changes

#### Examples

```bash
# Good commits
git commit -m "feat(compression): add bzip2 support"
git commit -m "fix(extract): prevent path traversal in ZIP extraction"
git commit -m "docs(readme): update installation instructions"
git commit -m "perf(compress): use buffer pool to reduce allocations"

# With body
git commit -m "feat(encryption): implement AES-256-GCM

- Add newEncryptedWriter and newEncryptedReader
- Use PBKDF2 with 100k iterations for key derivation
- Include random salt and nonce in output
- Add comprehensive tests for encryption/decryption

Closes #42"
```

### Commit Best Practices

-   **One logical change per commit**
-   **Write descriptive commit messages**
-   **Reference issues**: Use `Closes #123` or `Fixes #123`
-   **Keep commits atomic**: Each commit should compile and pass tests
-   **Rebase before merging**: Keep history clean

---

## Pull Request Process

### Before Submitting

1. **Sync with upstream**: Rebase on latest main
2. **Run all tests**: `go test ./... -race`
3. **Check formatting**: `gofmt -w .`
4. **Run linter**: `golangci-lint run`
5. **Update documentation**: README, comments, etc.
6. **Add tests**: For all new functionality

### PR Template

When creating a PR, include:

```markdown
## Description

Brief description of what this PR does.

## Type of Change

-   [ ] Bug fix
-   [ ] New feature
-   [ ] Breaking change
-   [ ] Documentation update

## Testing

Describe how you tested this change:

-   [ ] Unit tests added/updated
-   [ ] Integration tests added/updated
-   [ ] Manual testing performed

## Checklist

-   [ ] Code follows project style guidelines
-   [ ] Self-review completed
-   [ ] Comments added for complex logic
-   [ ] Documentation updated
-   [ ] Tests pass locally
-   [ ] No new warnings introduced

## Related Issues

Closes #(issue number)

## Screenshots (if applicable)
```

### Review Process

1. **Automated checks**: CI must pass (tests, linting)
2. **Code review**: At least one maintainer approval required
3. **Discussion**: Address all review comments
4. **Approval**: Maintainer approves and merges

### After Merge

-   Delete your feature branch
-   Update your local main branch
-   Start next feature from clean main

---

## Project Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────┐
│           CLI Interface                 │
│  (Parse flags, validate, orchestrate)   │
└──────────────┬──────────────────────────┘
               │
               ├──────────────┬────────────┐
               │              │            │
               ▼              ▼            ▼
    ┌──────────────┐  ┌──────────┐  ┌─────────┐
    │  Compression │  │ Extract  │  │  List   │
    │   Handler    │  │ Handler  │  │ Handler │
    └──────┬───────┘  └────┬─────┘  └────┬────┘
           │               │              │
           ▼               ▼              ▼
    ┌──────────────────────────────────────┐
    │        Format Handlers               │
    │   (ZIP, TAR.GZ implementations)      │
    └──────────────┬───────────────────────┘
                   │
                   ├─────────┬──────────┐
                   │         │          │
                   ▼         ▼          ▼
           ┌──────────┐ ┌────────┐ ┌──────────┐
           │Encryption│ │ Worker │ │   I/O    │
           │  Layer   │ │  Pool  │ │ Helpers  │
           └──────────┘ └────────┘ └──────────┘
```

### Key Components

1. **CLI Layer**: Flag parsing, validation, user interaction
2. **Operation Layer**: Compress, extract, list operations
3. **Format Layer**: Format-specific implementations (ZIP, TAR.GZ)
4. **Utility Layer**: Encryption, concurrency, I/O helpers

### Design Patterns Used

-   **Worker Pool**: For parallel extraction
-   **Strategy Pattern**: Different compression formats
-   **Decorator Pattern**: Encryption wrapper around I/O
-   **Command Pattern**: CLI actions (compress, extract, list)

---

## Security Considerations

### Security Review Checklist

When reviewing code, ensure:

-   [ ] No path traversal vulnerabilities
-   [ ] All file operations validate paths
-   [ ] Errors don't leak sensitive information
-   [ ] Crypto operations use crypto/rand
-   [ ] No hardcoded secrets or passwords
-   [ ] Input validation on all user data
-   [ ] Resource limits to prevent DoS
-   [ ] Proper cleanup of sensitive data

### Reporting Security Vulnerabilities

**DO NOT** open public issues for security vulnerabilities.

Instead:

1. Email: oss@cubetiqs.com
2. Include: Detailed description and reproduction steps
3. Wait: We'll respond within 48 hours
4. Coordinate: Work with us on fix and disclosure

### Security Testing

```bash
# Run security scanner
gosec ./...

# Check for vulnerabilities in dependencies
go list -json -m all | nancy sleuth

# Test with malicious inputs
go test ./... -fuzz=FuzzExtract -fuzztime=30s
```

---

## Documentation

### Documentation Requirements

All contributions must include appropriate documentation:

#### Code Comments

```go
// Package-level comment
// Package main implements the GoArchive CLI application.
package main

// Public function documentation
// Compress creates an archive from the input path with the given options.
// It returns an error if the input doesn't exist or archiving fails.
func Compress(input, output string, opts ArchiveOptions) error {
    // Implementation
}
```

#### README Updates

Update README.md for:

-   New features or commands
-   Changed behavior
-   New configuration options
-   Performance improvements

#### Example Code

Include examples for new features:

```go
// Example usage:
//
//   opts := ArchiveOptions{
//       Format:  FormatZip,
//       Workers: 4,
//       Verbose: true,
//   }
//   if err := Compress("myfiles/", "archive.zip", opts); err != nil {
//       log.Fatal(err)
//   }
```

---

## Community

### Communication Channels

-   **GitHub Issues**: Bug reports and feature requests
-   **GitHub Discussions**: Questions and general discussion
-   **Pull Requests**: Code contributions
-   **Email**: For private matters (security, etc.)

### Getting Help

-   Read the [README](README.md) first
-   Search existing issues
-   Check [GitHub Discussions](https://github.com/cubetiqlabs/gar.git/discussions)
-   Ask specific, well-formed questions

### Recognition

Contributors are recognized in:

-   CONTRIBUTORS.md file
-   Release notes
-   Annual contributor highlights

---

## License

By contributing to GoArchive, you agree that your contributions will be licensed under the MIT License.

---

## Questions?

If you have questions about contributing:

1. Check this document thoroughly
2. Search existing issues and discussions
3. Open a new discussion on GitHub
4. Tag your question appropriately

Thank you for contributing to GoArchive (gar)! 🎉

---

<div align="center">

**[⬆ Back to Top](#contributing-to-gar)**

Last Updated: October 2025

</div>
