# Contributing Guide

Thank you for your interest in contributing to Solana Insider Monitor! This guide will help you understand how to contribute to the project effectively.

## Ways to Contribute

There are many ways to contribute to Solana Insider Monitor:

- 🐛 **Reporting bugs** - Help us identify issues
- 💡 **Feature suggestions** - Share your ideas for improvements
- 📝 **Documentation improvements** - Help us make the docs more clear and comprehensive
- 🔍 **Code reviews** - Review pull requests from other contributors
- 💻 **Code contributions** - Implement new features or fix bugs

## Getting Started

### Prerequisites

Before you begin, ensure you have the following:

- Go 1.19 or later
- Git
- Basic understanding of Solana blockchain (for certain features)

### Setting Up Your Development Environment

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/yourusername/insider-monitor.git
   cd insider-monitor
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/accursedgalaxy/insider-monitor.git
   ```
4. **Install dependencies**:
   ```bash
   go mod download
   ```
5. **Install pre-commit hooks**:
   ```bash
   pre-commit install
   ```

## Development Workflow

### 1. Branching Strategy

- `main`: Production-ready code
- Feature branches: `feature/your-feature-name`
- Bug fix branches: `fix/bug-description`
- Hotfix branches: `hotfix/urgent-fix`

### 2. Making Changes

1. **Sync your main branch**:
   ```bash
   git checkout main
   git pull upstream main
   ```
2. **Create a new branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes**
4. **Run tests**:
   ```bash
   make test
   ```
5. **Commit using conventional commit messages**:
   ```bash
   git commit -m "feat: add new feature"
   git commit -m "fix: resolve bug"
   git commit -m "docs: update documentation"
   ```
6. **Push your branch**:
   ```bash
   git push origin feature/your-feature-name
   ```
7. **Create a Pull Request** on GitHub

## Pull Request Guidelines

### Creating Pull Requests

When creating a PR, please:

1. **Link related issues** - Reference any issues your PR addresses
2. **Provide context** - Explain what your changes do and why they're needed
3. **Include screenshots** - For UI changes, include before/after screenshots
4. **Update documentation** - Update relevant documentation if necessary

### PR Review Process

1. A maintainer will review your PR within a few days
2. Address any feedback or requested changes
3. Once approved, a maintainer will merge your PR
4. Your contribution will be included in the next release

## Code Standards

### Go Code Guidelines

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Maintain [Effective Go](https://golang.org/doc/effective_go) principles
- Use meaningful variable and function names
- Document all exported functions, types, and constants
- Run `golangci-lint` before committing:
  ```bash
  make lint
  ```

### Testing Guidelines

- Write tests for all new features and bug fixes
- Maintain test coverage (aim for at least 80% coverage)
- Ensure all tests pass before submitting your PR

## Documentation Guidelines

When contributing to documentation:

- Use clear, concise language
- Follow the existing style and format
- Update table of contents when adding new sections
- Preview changes locally using `mkdocs serve`
- Check for broken links and references

## Release Process

Releases follow semantic versioning (MAJOR.MINOR.PATCH):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backwards compatible)
- **PATCH**: Bug fixes

### Creating a Release

1. Ensure all changes are merged to `main`
2. Create and push a new tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. GitHub Actions will automatically:
   - Run tests
   - Build binaries
   - Create a GitHub release
   - Upload artifacts

## Community Guidelines

- Be respectful and inclusive in all interactions
- Focus on constructive feedback
- Help newcomers get started
- Acknowledge others' contributions
- Follow the [Code of Conduct](CODE_OF_CONDUCT.md)

## Getting Help

If you need help with your contribution:

- Join our [Discord community](https://discord.gg/7vY9ZBPdya)
- Open an issue with your question
- Reach out to the maintainers directly

Thank you for contributing to Solana Insider Monitor! Your efforts help make this project better for everyone.
