# Contributing to Bridge Bid Tutor

Thank you for your interest in contributing to the Bridge Bid Tutor! We appreciate your time and effort in helping improve this project.

## Code of Conduct

This project adheres to the [Contributor Covenant](https://www.contributor-covenant.org/version/2/0/code_of_conduct/). By participating, you are expected to uphold this code.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your changes
4. Make your changes
5. Run tests and linters
6. Commit and push your changes to your fork
7. Create a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Make
- golangci-lint (install with `make deps`)

### Building the Project

```bash
make build
```

### Running Tests

```bash
make test
```

### Running Linters

```bash
make lint
```

To automatically fix linting issues:

```bash
make lint-fix
```

## Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `gofmt` and `goimports` before committing
- Keep functions small and focused
- Write tests for new functionality
- Document exported functions and types

## Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally

## Pull Requests

1. Update the README.md with details of changes if needed
2. Ensure tests pass and coverage is maintained or improved
3. Lint your code before submitting
4. Keep pull requests focused on a single feature or fix
5. Include tests for new functionality

## Reporting Bugs

Please open an issue with:

- A clear title and description
- Steps to reproduce the issue
- Expected vs actual behavior
- Any relevant logs or screenshots

## Feature Requests

Feel free to suggest new features by opening an issue. Please include:

- A clear description of the feature
- Why this feature would be valuable
- Any alternative solutions you've considered

## License

By contributing, you agree that your contributions will be licensed under the project's license.
