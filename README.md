# goinit

A command-line tool to quickly initialize Go projects with pre-configured architectures and templates.

## Features

- Initialize Go projects with predefined architecture templates
- Automatic module name replacement throughout the project
- Support for multiple architecture patterns
- Simple CLI interface
- Detailed processing logs

## Installation

### From Source

```bash
go install github.com/openframebox/goinit@latest
```

Or build from source:

```bash
git clone https://github.com/openframebox/goinit.git
cd goinit
go build -o goinit
```

Then move the binary to your PATH:

```bash
# On macOS/Linux
sudo mv goinit /usr/local/bin/

# Or add to your PATH in ~/.bashrc or ~/.zshrc
export PATH=$PATH:/path/to/goinit
```

### Verify Installation

```bash
goinit --help
```

Check your installed version:

```bash
goinit version
# Output: goinit version 1.0.0
```

## Usage

### Check Version

Display the current version of goinit:

```bash
goinit version
```

### List Available Architectures

View all available project templates:

```bash
goinit list
```

Output:
```
Available Architectures:

┌─ layered@v0
│  Name: Layered Architecture (Pre-release)
│  Description:
│    Simple go project with layered architecture. Focus on simplicity and
│    ease of use for small-medium projects.
└─

Usage:
  goinit create --project <dir> --module <name> --architecture <arch>
```

### Create a New Project

Create a new Go project with a specific architecture:

```bash
goinit create --project myapp --module github.com/username/myapp --architecture layered@v0
```

**Flags:**
- `--project` or `-p`: Directory path where the project will be created (required)
- `--module` or `-m`: Go module name (required)
- `--architecture` or `-a`: Architecture template to use (required)

### Example

```bash
# Create a new project called "myapi"
goinit create \
  --project myapi \
  --module github.com/myusername/myapi \
  --architecture layered@v0

# Output:
# Creating new project: github.com/myusername/myapi
# Architecture: Layered Architecture (Pre-release)
# Target directory: myapi
#
# Downloading template from https://github.com/...
# Downloaded 30539 bytes
# Extracting template...
# Extraction complete
# Processing template files...
#
#   Processing file: .env.example (replaced module name)
#   Processing file: .gitignore
#   Processing file: README.md (replaced module name)
#   ...
#
# ✓ Processed 64 files (28 files with module name replacements)
#
# ✓ Project created successfully at myapi
#
# Next steps:
#   cd myapi
#   go mod tidy
```

After creation, navigate to your project and install dependencies:

```bash
cd myapi
go mod tidy
```

## Available Architectures

### layered@v0

A simple layered architecture pattern focused on simplicity and ease of use for small to medium projects. Includes:

- Configuration management
- Repository pattern
- Service layer
- HTTP routing with Fiber
- Authentication middleware
- Database integration (GORM)
- Event system
- Queue system
- Email support
- Storage abstraction

## Requirements

- Go 1.25.0 or higher
- Internet connection (to download templates)

## How It Works

1. **Load Configuration**: Reads available architectures from embedded `goinit.json`
2. **Validate Input**: Ensures the target directory is empty and architecture exists
3. **Download Template**: Fetches the template archive from GitHub
4. **Extract Files**: Extracts the archive to a temporary location
5. **Process Files**: Recursively processes all files, replacing module placeholders
6. **Copy to Target**: Moves processed files to your project directory
7. **Cleanup**: Removes temporary files

## Configuration

Templates are configured in `goinit.json`:

```json
{
  "availableArchitectures": [
    "layered@v0"
  ],
  "boilerplates": {
    "layered@v0": {
      "name": "Layered Architecture (Pre-release)",
      "description": "Simple go project with layered architecture...",
      "archiveUrl": "https://github.com/openframebox/goinit-layered/archive/refs/tags/v0.0.0.tar.gz",
      "modulePlaceholder": "mymodule"
    }
  }
}
```

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions:

- Open an issue on [GitHub](https://github.com/openframebox/goinit/issues)
- Check existing issues for solutions

## Roadmap

- [ ] Support for custom template repositories
- [ ] Interactive mode for project creation
- [ ] Template versioning support
- [ ] Custom placeholder configuration
- [ ] Project update capabilities

## Credits

Developed by the OpenFrameBox team.
