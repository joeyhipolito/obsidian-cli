# Contributing to obsidian-cli

Thanks for your interest in contributing! This document covers the basics.

## Development Setup

```bash
git clone https://github.com/joeyhipolito/obsidian-cli.git
cd obsidian-cli
make build
```

Requires Go 1.21+. The only external dependency is `modernc.org/sqlite`.

## Making Changes

1. Fork the repo and create a branch from `main`
2. Make your changes
3. Run `make lint` (formats code and runs vet)
4. Run `make test` to ensure tests pass
5. Submit a pull request

## Code Style

- Follow standard Go conventions (`gofmt`)
- Add tests for new commands or search functionality
- Use the existing package structure (`internal/cmd`, `internal/vault`, `internal/index`, etc.)

## Reporting Issues

Open an issue with:
- What you expected to happen
- What actually happened
- Steps to reproduce
- Your Go version and OS

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
