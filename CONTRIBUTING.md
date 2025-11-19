# Contributing to goinit

Thank you for considering contributing to goinit! We welcome contributions from the community.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Your environment (OS, Go version, goinit version)
- Any relevant logs or error messages

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- A clear and descriptive title
- Detailed description of the proposed functionality
- Why this enhancement would be useful
- Any examples or mockups if applicable

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following the code style guidelines
3. **Add tests** if applicable
4. **Update documentation** if you're changing functionality
5. **Ensure tests pass** by running `go test ./...`
6. **Commit your changes** with clear commit messages
7. **Push to your fork** and submit a pull request

## Development Setup

### Prerequisites

- Go 1.25.0 or higher
- Git

### Getting Started

1. Fork and clone the repository:

```bash
git clone https://github.com/YOUR_USERNAME/goinit.git
cd goinit
```

2. Install dependencies:

```bash
go mod download
```

3. Build the project:

```bash
go build -o goinit
```

4. Run tests:

```bash
go test ./...
```

## Project Structure

```
goinit/
├── main.go                    # Main CLI application
├── goinit.json               # Architecture templates configuration
├── internal/
│   ├── config/               # Configuration loading and validation
│   ├── downloader/           # Archive download functionality
│   ├── extractor/            # Archive extraction logic
│   ├── processor/            # File processing and placeholder replacement
│   └── validator/            # Input validation
├── README.md
├── LICENSE
└── CONTRIBUTING.md
```

## Code Style Guidelines

### General Guidelines

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Write clear comments for exported functions and types
- Keep functions focused and concise

### Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `go fmt` before committing
- Run `go vet` to catch common mistakes
- Handle errors explicitly

### Example

```go
// Good
func ProcessFile(path string) error {
    content, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read file: %w", err)
    }
    // Process content...
    return nil
}

// Bad
func ProcessFile(path string) error {
    content, _ := os.ReadFile(path) // Ignoring errors
    // Process content...
    return nil
}
```

## Commit Message Guidelines

Write clear and meaningful commit messages:

```
Add feature to support custom templates

- Implement custom template URL parsing
- Add validation for template format
- Update documentation with examples
```

Format:
- Use present tense ("Add feature" not "Added feature")
- First line should be a brief summary (50 chars or less)
- Optionally add a detailed description after a blank line
- Reference issues and pull requests when relevant

## Adding New Architecture Templates

To add a new architecture template:

1. Create the template repository (or archive)
2. Update `goinit.json`:

```json
{
  "availableArchitectures": [
    "layered@v0",
    "your-architecture@v1"
  ],
  "boilerplates": {
    "your-architecture@v1": {
      "name": "Your Architecture Name",
      "description": "Description of your architecture...",
      "archiveUrl": "https://github.com/user/repo/archive/refs/tags/v1.0.0.tar.gz",
      "modulePlaceholder": "mymodule"
    }
  }
}
```

3. Ensure your template:
   - Uses a consistent placeholder for module names
   - Includes a README with usage instructions
   - Has a valid Go module structure
   - Doesn't contain sensitive information

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/config
```

### Manual Testing

1. Build the binary:
```bash
go build -o goinit
```

2. Test the list command:
```bash
./goinit list
```

3. Test project creation:
```bash
./goinit create --project test-project --module github.com/test/app --architecture layered@v0
```

4. Verify the created project:
```bash
cd test-project
go mod tidy
go build
```

## Documentation

When adding features or making changes:

- Update README.md if user-facing functionality changes
- Add inline comments for complex logic
- Update this CONTRIBUTING.md if development process changes
- Consider adding examples for new features

## Release Process

(For maintainers)

1. Update version in relevant files
2. Update CHANGELOG.md with changes
3. Create a git tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
4. Push the tag: `git push origin v1.0.0`
5. Create a GitHub release with release notes

## Questions?

If you have questions about contributing:

- Open an issue with the "question" label
- Check existing documentation
- Review closed issues for similar questions

## Recognition

Contributors will be recognized in:
- GitHub contributors page
- Release notes for significant contributions

Thank you for contributing to goinit!
