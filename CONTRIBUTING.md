# Contribution Guide

Thank you for considering contributing to OSP (Open Source Pilot)!

## Development Process

1. Fork the project
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

## Commit Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `perf`: Performance optimization
- `test`: Test-related changes
- `chore`: Build process or auxiliary tool changes

Examples:
feat: Add automatic weekly report generation
fix: Fix expired authentication issue
docs: Update installation instructions

## Code Style

- Use `gofmt` to format code
- Follow [Effective Go](https://golang.org/doc/effective_go.html) recommendations
- Add necessary comments and documentation
- Ensure test coverage

## Development Setup

1. Clone the project
```bash
git clone https://github.com/yourusername/osp.git

2. Install dependencies
```bash
go mod download

3. Run tests
```bash
go test ./...

4. Build the project
```bash
go build ./cmd/osp

## Checklist Before Submitting a PR

- [ ] Pass all tests
- [ ] Update related documentation
- [ ] Add necessary test cases
- [ ] Follow code conventions
- [ ] Commit message follows guidelines

## Reporting Issues

When reporting issues, please include the following information:

1. Description of the problem
2. Reproduction steps
3. Expected behavior
4. Actual behavior
5. Environment information
   - OSP version
   - Go version
   - Operating system
   - Other relevant information

## Feature Suggestions

We welcome new feature suggestions! Please:

1. Check existing issues to avoid duplicates
2. Provide a detailed description of the new feature
3. Explain use cases
4. Consider implementation options

## Code of Conduct

Please refer to our [Code of Conduct](CODE_OF_CONDUCT.md).

## License

By contributing code, you agree to license it under the [MIT License](LICENSE).
