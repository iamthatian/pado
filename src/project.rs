use crate::config::GlobalConfig;
use crate::error::ParkourError;
use ignore::WalkBuilder;
use phf::{phf_set, Set};
use std::path::{Path, PathBuf};

pub static PROJECT_FILES: Set<&'static str> = phf_set! {
    // VC
    ".projectile",
    ".hg",
    ".svn",
    ".bzr",
    ".cvsignore",
    ".gitignore",
    ".gitattributes",
    ".git",
    // Rust
    "Cargo.toml",
    "Cargo.lock",
    "rust-toolchain",
    "rustfmt.toml",
    "clippy.toml",
    // Python
    "pyproject.toml",
    "setup.py",
    "setup.cfg",
    "requirements.txt",
    "Pipfile",
    "Pipfile.lock",
    "tox.ini",
    "pytest.ini",
    "MANIFEST.in",
    "poetry.lock",
    "pyvenv.cfg",
    "requirements.lock",
    "uv.lock",
    ".python-version",
    // Node.js
    "package.json",
    "pnpm-workspace.yaml",
    "yarn.lock",
    "pnpm-lock.yaml",
    "bun.lockb",
    "package-lock.json",
    "tsconfig.json",
    "jsconfig.json",
    "webpack.config.js",
    "vite.config.js",
    "rollup.config.js",
    "eslint.config.js",
    ".eslintrc",
    ".eslintrc.js",
    ".eslintrc.json",
    ".prettierrc",
    ".prettierrc.js",
    ".prettierrc.json",
    ".npmrc",
    // Java/Kotlin/Scala
    "pom.xml",
    "build.gradle",
    "build.gradle.kts",
    "settings.gradle",
    "settings.gradle.kts",
    "gradlew",
    "gradlew.bat",
    "gradle.properties",
    "mvnw",
    "mvnw.cmd",
    "build.sbt",
    "project/build.properties",
    // C/C++
    "CMakeLists.txt",
    "Makefile",
    "makefile",
    "configure",
    "configure.ac",
    "meson.build",
    "meson_options.txt",
    "build.ninja",
    // Docker
    "Dockerfile",
    "docker-compose.yml",
    "docker-compose.yaml",
    ".dockerignore",
    // Kubernetes/Helm
    ".helm/",
    "Chart.yaml",
    "kustomization.yaml",
    // CI/CD
    "ansible.cfg",
    "playbook.yml",
    ".circleci/",
    ".github/",
    ".gitlab-ci.yml",
    ".travis.yml",
    ".jenkinsfile",
    "azure-pipelines.yml",
    // PHP
    "composer.json",
    "composer.lock",
    "phpunit.xml",
    "phpunit.xml.dist",
    // OCaml
    "dune",
    "dune-project",
    "opam",
    // Ruby
    "Gemfile",
    "Gemfile.lock",
    "Rakefile",
    "config.ru",
    // Go
    "go.mod",
    "go.sum",
    // Zig
    "zig.mod",
    "build.zig",
    // Elixir
    "mix.exs",
    "mix.lock",
    // Swift
    "Package.swift",
    "Package.resolved",
    // .NET/C#
    "*.csproj",
    "*.sln",
    "*.fsproj",
    "*.vbproj",
    "global.json",
    "nuget.config",
    // Haskell
    "stack.yaml",
    "stack.yaml.lock",
    "*.cabal",
    "cabal.project",
    // Terraform
    "main.tf",
    "*.tf",
    ".terraform/",
    ".terragrunt.hcl",
    "terraform.tfstate",
    // Lua
    "*.rockspec",
    // Nix
    "flake.lock",
    "flake.nix",
    "default.nix",
    "shell.nix",
    // Other
    ".venv/",
    "env/",
    ".vscode/",
    ".idea/",
    ".editorconfig",
    ".code-workspace",
};

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum ProjectType {
    Rust,
    Node,
    Python,
    Go,
    Java,
    Ruby,
    PHP,
    Elixir,
    Scala,
    Swift,
    Dotnet,
    Haskell,
    OCaml,
    Zig,
    Terraform,
    Lua,
    Nix,
    Docker,
    Git,
    Unknown,
}

impl ProjectType {
    pub fn as_str(&self) -> &'static str {
        match self {
            ProjectType::Rust => "rust",
            ProjectType::Node => "node",
            ProjectType::Python => "python",
            ProjectType::Go => "go",
            ProjectType::Java => "java",
            ProjectType::Ruby => "ruby",
            ProjectType::PHP => "php",
            ProjectType::Elixir => "elixir",
            ProjectType::Scala => "scala",
            ProjectType::Swift => "swift",
            ProjectType::Dotnet => "dotnet",
            ProjectType::Haskell => "haskell",
            ProjectType::OCaml => "ocaml",
            ProjectType::Zig => "zig",
            ProjectType::Terraform => "terraform",
            ProjectType::Lua => "lua",
            ProjectType::Nix => "nix",
            ProjectType::Docker => "docker",
            ProjectType::Git => "git",
            ProjectType::Unknown => "unknown",
        }
    }
}

pub fn contains_project_file(dir: &Path) -> std::io::Result<bool> {
    contains_project_file_with_config(dir, &GlobalConfig::default())
}

pub fn contains_project_file_with_config(
    dir: &Path,
    config: &GlobalConfig,
) -> std::io::Result<bool> {
    for entry in dir.read_dir()? {
        let entry = entry?;
        if let Some(name) = entry.file_name().to_str() {
            if PROJECT_FILES.contains(name) {
                return Ok(true);
            }
            if config.markers.additional.iter().any(|m| m == name) {
                return Ok(true);
            }
        }
    }
    Ok(false)
}

pub fn find_project_root(start: &Path) -> Result<PathBuf, ParkourError> {
    let config = GlobalConfig::load().unwrap_or_default();
    find_project_root_with_config(start, &config)
}

pub fn find_project_root_with_config(
    start: &Path,
    config: &GlobalConfig,
) -> Result<PathBuf, ParkourError> {
    for ancestor in start.ancestors() {
        if contains_project_file_with_config(ancestor, config).map_err(ParkourError::Io)? {
            return Ok(ancestor.to_path_buf());
        }
    }
    Err(ParkourError::NoProjectRoot(start.to_path_buf()))
}

pub fn is_project_root(path: &Path) -> bool {
    contains_project_file(path).unwrap_or(false)
}

pub fn detect_project_type(root: &Path) -> ProjectType {
    if root.join("Cargo.toml").exists() {
        return ProjectType::Rust;
    }

    if root.join("package.json").exists() {
        return ProjectType::Node;
    }

    if root.join("uv.lock").exists()
        || root.join("pyproject.toml").exists()
        || root.join("setup.py").exists()
        || root.join("requirements.txt").exists()
    {
        return ProjectType::Python;
    }

    if root.join("go.mod").exists() {
        return ProjectType::Go;
    }

    if root.join("pom.xml").exists()
        || root.join("build.gradle").exists()
        || root.join("build.gradle.kts").exists()
    {
        return ProjectType::Java;
    }

    if root.join("build.sbt").exists() {
        return ProjectType::Scala;
    }

    if root.join("Gemfile").exists() {
        return ProjectType::Ruby;
    }

    if root.join("composer.json").exists() {
        return ProjectType::PHP;
    }

    if root.join("mix.exs").exists() {
        return ProjectType::Elixir;
    }

    if root.join("Package.swift").exists() {
        return ProjectType::Swift;
    }

    if std::fs::read_dir(root)
        .ok()
        .map(|entries| {
            entries
                .filter_map(|e| e.ok())
                .any(|entry| {
                    entry
                        .file_name()
                        .to_str()
                        .map(|name| {
                            name.ends_with(".csproj")
                                || name.ends_with(".fsproj")
                                || name.ends_with(".sln")
                        })
                        .unwrap_or(false)
                })
        })
        .unwrap_or(false)
        || root.join("global.json").exists()
    {
        return ProjectType::Dotnet;
    }

    if root.join("stack.yaml").exists()
        || std::fs::read_dir(root)
            .ok()
            .map(|entries| {
                entries
                    .filter_map(|e| e.ok())
                    .any(|entry| {
                        entry
                            .file_name()
                            .to_str()
                            .map(|name| name.ends_with(".cabal"))
                            .unwrap_or(false)
                    })
            })
            .unwrap_or(false)
    {
        return ProjectType::Haskell;
    }

    if root.join("dune").exists() || root.join("dune-project").exists() || root.join("opam").exists() {
        return ProjectType::OCaml;
    }

    if root.join("build.zig").exists() {
        return ProjectType::Zig;
    }

    if std::fs::read_dir(root)
        .ok()
        .map(|entries| {
            entries
                .filter_map(|e| e.ok())
                .any(|entry| {
                    entry
                        .file_name()
                        .to_str()
                        .map(|name| name.ends_with(".tf"))
                        .unwrap_or(false)
                })
        })
        .unwrap_or(false)
        || root.join(".terraform").exists()
    {
        return ProjectType::Terraform;
    }

    if std::fs::read_dir(root)
        .ok()
        .map(|entries| {
            entries
                .filter_map(|e| e.ok())
                .any(|entry| {
                    entry
                        .file_name()
                        .to_str()
                        .map(|name| name.ends_with(".rockspec"))
                        .unwrap_or(false)
                })
        })
        .unwrap_or(false)
    {
        return ProjectType::Lua;
    }

    if root.join("flake.nix").exists()
        || root.join("default.nix").exists()
        || root.join("shell.nix").exists()
    {
        return ProjectType::Nix;
    }

    if root.join("Dockerfile").exists()
        || root.join("docker-compose.yml").exists()
        || root.join("docker-compose.yaml").exists()
    {
        return ProjectType::Docker;
    }

    if root.join(".git").exists() {
        return ProjectType::Git;
    }

    ProjectType::Unknown
}

pub fn list_project_files(
    root: &Path,
    pattern: Option<&str>,
) -> Result<Vec<PathBuf>, ParkourError> {
    let mut files = Vec::new();

    for result in WalkBuilder::new(root)
        .hidden(false)
        .ignore(true)
        .git_ignore(true)
        .build()
    {
        let entry = result
            .map_err(|e| std::io::Error::new(std::io::ErrorKind::Other, format!("walking project files: {}", e)))?;

        if entry.file_type().map(|ft| ft.is_file()).unwrap_or(false) {
            let path = entry.path().to_path_buf();

            if let Some(pat) = pattern {
                if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
                    if glob_match(pat, name) {
                        files.push(path);
                    }
                }
            } else {
                files.push(path);
            }
        }
    }

    Ok(files)
}

pub fn glob_match(pattern: &str, text: &str) -> bool {
    if pattern.starts_with('*') && pattern.len() > 1 {
        text.ends_with(&pattern[1..])
    } else if pattern.ends_with('*') && pattern.len() > 1 {
        text.starts_with(&pattern[..pattern.len() - 1])
    } else if pattern.contains('*') {
        let parts: Vec<&str> = pattern.split('*').collect();
        if parts.len() == 2 {
            text.starts_with(parts[0]) && text.ends_with(parts[1])
        } else {
            pattern == text
        }
    } else {
        pattern == text
    }
}

pub struct ProjectInfo {
    pub root: PathBuf,
    pub project_type: ProjectType,
    pub file_count: usize,
}

pub fn get_project_info(root: &Path) -> Result<ProjectInfo, ParkourError> {
    let project_type = detect_project_type(root);
    let files = list_project_files(root, None)?;
    let file_count = files.len();

    Ok(ProjectInfo {
        root: root.to_path_buf(),
        project_type,
        file_count,
    })
}

pub fn discover_projects(
    search_path: &Path,
    max_depth: usize,
) -> Result<Vec<PathBuf>, ParkourError> {
    use std::collections::HashSet;
    use walkdir::WalkDir;

    let mut projects = Vec::new();
    let mut visited = HashSet::new();

    for entry in WalkDir::new(search_path)
        .max_depth(max_depth)
        .into_iter()
        .filter_entry(|e| {
            let path = e.path();
            if path == search_path || !path.starts_with(search_path) {
                return true;
            }
            if let Some(name) = e.file_name().to_str() {
                !name.starts_with('.')
            } else {
                true
            }
        })
    {
        let entry = entry
            .map_err(|e| std::io::Error::new(std::io::ErrorKind::Other, format!("discovering projects: {}", e)))?;
        let path = entry.path();

        if path.is_dir() && is_project_root(path) {
            // NO duplicates
            if visited.insert(path.to_path_buf()) {
                projects.push(path.to_path_buf());
            }
        }
    }

    Ok(projects)
}
