# padofetch

Project information and statistics tool. Like `neofetch` but for your projects.

Part of the pado ecosystem - companion to `pd` for detailed project insights.

## Installation

```bash
cargo install padofetch
```

## Usage

```bash
# Show comprehensive project information (default)
padofetch
padofetch info

# Check project health
padofetch health

# Show dependencies
padofetch deps

# Check for outdated dependencies
padofetch outdated
```

## Features

### Project Information

```bash
padofetch info
```

Shows:
- **Project metadata**: name, root, type(s)
- **Language statistics**: Line counts per language with visual bars
- **Git information**: Contributors, commit counts, last commit time
- **File statistics**: Total files, code/comments/blanks breakdown

Example output:
```
Project Information

Name: myproject
Root: /home/user/projects/myproject
Type(s): rust, git
Total files: 142

Language Statistics:
  Rust         ████████████████████  82.3%  (2345 lines)
  TOML         ██░░░░░░░░░░░░░░░░░░  10.2%  (123 lines)
  Markdown     █░░░░░░░░░░░░░░░░░░░   5.1%  (234 lines)

Total lines: 2791 (2345 code, 123 comments, 323 blanks)

Git Information:
  Total commits: 235
  Contributors: 3
    - Alice (145 commits)
    - Bob (67 commits)
    - Carol (23 commits)
  Last commit: 2 hours ago
```

### Health Check

```bash
padofetch health
```

Checks:
- Git repository setup
- `.gitignore` presence
- Language-specific checks (lock files, virtual environments, etc.)
- Documentation (README)
- License file

Example:
```
🩺 Project Health Check

Project: myproject
Type(s): rust

Healthy:
  ✓ Git repository
  ✓ .gitignore present
  ✓ Cargo.lock present
  ✓ README present
  ✓ LICENSE present

✓ Project looks healthy!
```

### Dependencies

```bash
padofetch deps
```

Shows project dependencies based on detected project type:
- **Rust**: Parses `Cargo.toml`
- **Node**: Suggests using `npm list` or `yarn list`
- **Python**: Shows `requirements.txt` if available

### Outdated Check

```bash
padofetch outdated
```

Checks for outdated dependencies using build system tools:
- **Rust**: `cargo outdated`
- **Node**: `npm outdated`, `yarn outdated`, etc.
- **Python**: `pip list --outdated`
- **Go**: `go list -u -m all`

## Shell Integration

If you use `pd` (pado), shell functions are automatically available:

```bash
# From anywhere in a project
pdinfo      # Alias for padofetch info
pdhealth    # Alias for padofetch health
```

## Project Type Support

Detects and analyzes:
- Rust (Cargo)
- Node.js (npm, yarn, pnpm, bun)
- Python (pip, Poetry)
- Go
- Java (Maven, Gradle)
- Ruby (Bundler)
- PHP (Composer)
- C/C++ (Make, CMake)
- Docker projects

## Output

All output is designed to be:
- **Human-readable**: Pretty formatting with colors and bars
- **Informative**: Shows the most relevant information first
- **Fast**: Minimal overhead, uses efficient libraries (tokei, gitoxide)

## Dependencies

- **tokei**: Fast language statistics
- **gitoxide**: Efficient Git repository analysis

## Why "padofetch"?

A portmanteau of "pado" (tide) and "fetch" - fetch information about your project like waves bringing messages from the sea.

## See Also

- [pado](https://crates.io/crates/pado) - Project management CLI
- [libpado](https://crates.io/crates/libpado) - Project detection library

## License

MIT
