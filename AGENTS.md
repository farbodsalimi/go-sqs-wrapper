# Claude Code Configuration

## Project Context
This is a Go SQS wrapper library that provides simplified access to AWS SQS operations. The project was created for Go 1.12 and needs modernization.

## Build and Test Commands
```bash
# Build the project
make build

# Run tests
make test

# Test with coverage
make test-cover

# Clean build artifacts
make clean

# Run the binary
make run
```

## Modernization Recommendations

### 1. Go Version and Dependencies (Critical)
- **Current**: Go 1.12 (released March 2019, EOL)
- **Recommendation**: Upgrade to Go 1.21+ for modern features and security
- **AWS SDK**: Currently using v1.34.0, should upgrade to AWS SDK Go v2 for better performance and modern APIs

### 2. Code Structure Issues

#### Global State Management
- **Issue**: Global `clientSession` variable in `pkg/backend/sqs.go:14`
- **Impact**: Not thread-safe, difficult to test, violates dependency injection
- **Fix**: Pass SQS client as struct field, use dependency injection

#### Error Handling
- **Issue**: Inconsistent error handling, some functions don't return errors
- **Examples**: 
  - `sqs.go:77` FetchMessages() logs but doesn't return errors
  - `sqs.go:103` DeleteMessages() logs and returns on first error
- **Fix**: Return errors consistently, let callers handle them

#### Method Receivers
- **Issue**: Mix of value and pointer receivers
- **Example**: `sqs.go:27` Init() uses value receiver but modifies struct
- **Fix**: Use pointer receivers consistently for structs that are modified

### 3. API Design Improvements

#### Constructor Pattern
- **Current**: `new(backend.SQSQueue).Init()` is awkward
- **Better**: `backend.NewSQSQueue()` constructor function

#### Configuration
- **Current**: Hardcoded environment variable loading
- **Better**: Accept configuration struct, allow multiple initialization methods

#### Context Support
- **Missing**: No context.Context support for cancellation/timeouts
- **Fix**: Add context to all AWS operations

### 4. Concurrency Issues

#### Async Processing
- **Issue**: `workers/futures.go:22` blocks on channel reads, not truly async
- **Fix**: Use proper worker pools with context cancellation

#### Stop Mechanism
- **Issue**: `futures.go:28` uses syscall.Kill() which is inappropriate
- **Fix**: Use context cancellation or channels for graceful shutdown

### 5. Modern Go Patterns

#### Missing Interfaces
- **Add**: Define interfaces for SQS operations to improve testability
- **Example**: `QueueReader`, `QueueWriter`, `QueueManager` interfaces

#### Generics (Go 1.18+)
- **Opportunity**: Message processing could benefit from type-safe generics
- **Example**: `ProcessMessage[T any](msg T) error` pattern

#### Structured Logging
- **Current**: Basic `log.Println()` usage
- **Better**: Use structured logging (slog, logrus, or zerolog)

### 6. Testing Improvements
- **Missing**: Comprehensive unit tests with mocking
- **Add**: Integration tests with localstack or AWS SQS mock
- **Tool**: Use testify/mock or gomock for AWS SDK mocking

### 7. Documentation
- **Add**: Go doc comments for all exported functions
- **Improve**: README examples with error handling
- **Add**: Architecture documentation

### 8. Performance Optimizations
- **Connection Pooling**: Reuse HTTP connections
- **Batch Operations**: Use batch delete for multiple messages
- **Memory**: Reduce allocations in hot paths

### 9. Security
- **Credentials**: Support IAM roles, not just access keys
- **Validation**: Input validation for queue URLs and message content
- **TLS**: Ensure secure connections (AWS SDK default)

### 10. CLI Improvements
- **Framework**: Replace google/subcommands with cobra for better UX
- **Config**: Support config files (YAML/JSON) in addition to env vars
- **Validation**: Better validation and help text

## Implementation Priority
1. **High**: Go version upgrade, AWS SDK v2, fix global state
2. **Medium**: Error handling, context support, proper constructors
3. **Low**: Generics, advanced testing, CLI improvements

## Breaking Changes Expected
Upgrading this codebase will require breaking changes to:
- Constructor patterns
- Error handling signatures  
- Configuration methods
- AWS SDK types (v1 → v2)

Consider versioning strategy (v2.0.0) for the modernized version.