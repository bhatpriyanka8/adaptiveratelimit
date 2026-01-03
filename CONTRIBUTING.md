# Contributing

Thank you for your interest in contributing to `adaptiveratelimit`.

The goal of this project is to provide a **simple, safe, and production-ready
adaptive rate limiter**. Contributions should align with that goal.

---

## Issues and Feature Requests

Before opening an issue, please check the existing issues.

### Issues
When reporting an issue, please include:
- expected behavior
- actual behavior
- a minimal reproduction if possible

### Features
For feature requests:
- explain the use case
- describe how it fits the core goal (safe adaptive rate limiting)
- keep scope focused

---

## Code Guidelines

- Follow standard Go style (`gofmt`, `goimports`)
- Keep public APIs minimal and intentional
- Prefer clarity and simplicity
- Keep code modular, with meaningful naming of files and folders and functions
- Tests covered for anything new

---

## Pull Request Process

Direct pushes to main are not allowed.

1. Fork the repository and create a feature branch from `main`
2. Make your changes
3. Run the full test suite:
   ``` bash
   go test -race ./...
   ```
4. Run linting:

``` bash
golangci-lint run ./...
```
5. If modifying core functions, consider updating benchmarks

6. Open a pull request with a clear description

## Testing Checklist

- All tests must pass

- Race detector must pass (-race)

- Benchmarks should be considered for performance-sensitive changes

- Logically cover all the functionalities, not just for the sake of lines

## Note
Always feel free to open up an issue to discuss for any queries.

Please review the design principles before proposing significant changes:

[DESIGN.md](./DESIGN.md)

Thanks!

