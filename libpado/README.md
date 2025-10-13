# libpado

A minimal Rust library for project root detection and build system discovery.

## Features

- **Project Root Detection**: Walk directory tree to find project roots based on marker files (`.git`, `Cargo.toml`, `package.json`, etc.)
- **Build System Detection**: Automatically identify build systems (Cargo, npm, Maven, Gradle, Poetry, Go, Make, etc.)
- **Project Type Detection**: Classify projects by type (Rust, Node, Python, Go, Java, etc.)
- **File Operations**: List project files with `.gitignore` support
- **Discovery**: Recursively find projects in directories

## Usage

Add to your `Cargo.toml`:

```toml
[dependencies]
libpado = "0.1.0"
```

### Find Project Root

```rust
use libpado::find_project_root;
use std::env;

let cwd = env::current_dir()?;
let root = find_project_root(&cwd)?;
println!("Project root: {}", root.display());
```

### Detect Build System

```rust
use libpado::BuildSystem;
use std::path::Path;

let root = Path::new("/path/to/project");
let build_system = BuildSystem::detect(root);

if let Some((cmd, args)) = build_system.build_command() {
    println!("Build with: {} {}", cmd, args.join(" "));
}
```

### Get Project Information

```rust
use libpado::get_project_info;
use std::env;

let cwd = env::current_dir()?;
let info = get_project_info(&cwd)?;

println!("Name: {}", info.root.file_name().unwrap().to_str().unwrap());
println!("Types: {:?}", info.project_types);
println!("Files: {}", info.file_count);
```

## Project Markers

Supports detection of:
- Version control: `.git`, `.hg`, `.svn`
- Build files: `Cargo.toml`, `package.json`, `go.mod`, `pom.xml`, `build.gradle`, `pyproject.toml`, `Makefile`, etc.
- Editor/IDE: `.projectile`, `.project`, `.idea`, `.vscode`

## Build Systems

Detects and provides commands for:
- Rust (Cargo)
- Node.js (npm, yarn, pnpm, bun)
- Python (Poetry, pip)
- Go
- Java (Maven, Gradle)
- C/C++ (Make, CMake)

## Error Handling

Uses `thiserror` for typed errors:

```rust
use libpado::{find_project_root, PadoError};

match find_project_root(&path) {
    Ok(root) => println!("Found: {}", root.display()),
    Err(PadoError::NoProjectRoot(p)) => {
        eprintln!("No project root found from: {}", p.display());
    }
    Err(e) => eprintln!("Error: {}", e),
}
```

## License

MIT
