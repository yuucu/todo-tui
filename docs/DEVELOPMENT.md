# Todo TUI - Development Guide

## Setup

```bash
git clone https://github.com/yuucu/todotui.git
cd todotui
go mod tidy
make build
./bin/todotui sample.todo.txt
```

## Commands

```bash
make build    # Build
make test     # Run tests
make lint     # Lint
make clean    # Clean
make release  # Release build
```

## Branch Naming

Use these prefixes for proper release categorization:

- `feat/` - New features → 🚀 Features
- `fix/` - Bug fixes → 🐛 Bug Fixes
- `docs/` - Documentation → 📚 Documentation
- `tests/` - Tests → 🧪 Tests
- `chore/` - Maintenance → 🔧 Maintenance

Examples:
```bash
git checkout -b feat/search-functionality
git checkout -b fix/crash-on-empty-file
git checkout -b docs/update-readme
```

## Testing

```bash
go test ./...           # All tests
go test -cover ./...    # With coverage
go test ./internal/ui   # Specific package
```

## Contributing

1. Fork repository
2. Create branch with proper prefix
3. Commit changes
4. Push and create PR

## Standards

- Follow Go style (`gofmt`, `go vet`)
- 80%+ test coverage
- Use branch naming convention 