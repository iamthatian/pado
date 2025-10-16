use crate::error::PadoError;
use globset::{Glob, GlobSetBuilder};
use ignore::WalkBuilder;
use phf::{Set, phf_set};
use std::{collections::HashSet, io, path::{Path, PathBuf}};

pub static PROJECT_FILES: Set<&'static str> = phf_set! {
    // VCS
    ".git",
    ".hg",
    ".fslckout",
    "_FOSSIL_",
    ".bzr",
    "_darcs",
    ".pijul",
    ".svn",
    ".sl",
    ".jj",
    ".cvsignore",
    ".gitignore",
    ".gitattributes",
    "CVS",
    "GTAGS",
    "TAGS",
    "cscope.out",

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

    // Node.js,JavaScript,TypeScript
    "package.json",
    "package-lock.json",
    "pnpm-workspace.yaml",
    "pnpm-lock.yaml",
    "yarn.lock",
    "bun.lockb",
    "tsconfig.json",
    "jsconfig.json",
    "webpack.config.js",
    "vite.config.js",
    "rollup.config.js",
    "gulpfile.js",
    "Gruntfile.js",
    "angular.json",
    ".angular-cli.json",
    "eslint.config.js",
    ".eslintrc",
    ".eslintrc.js",
    ".eslintrc.json",
    ".prettierrc",
    ".prettierrc.js",
    ".prettierrc.json",
    ".npmrc",

    // Java,Kotlin,Scala
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
    "build.sc",
    "build.mill",
    "project/build.properties",
    ".bloop/bloop.settings.json",

    // C,C++,Build systems
    "CMakeLists.txt",
    "CMakePresets.json",
    "CMakeUserPresets.json",
    "Makefile",
    "makefile",
    "GNUMakefile",
    "configure",
    "configure.ac",
    "configure.in",
    "meson.build",
    "meson_options.txt",
    "build.ninja",
    "SConstruct",
    "xmake.lua",

    // Docker
    "Dockerfile",
    "docker-compose.yml",
    "docker-compose.yaml",
    ".dockerignore",

    // Kubernetes,Helm
    "Chart.yaml",
    "kustomization.yaml",
    ".helm/",

    // CI/CD
    "ansible.cfg",
    "playbook.yml",
    ".circleci/",
    ".github/",
    ".gitlab-ci.yml",
    ".travis.yml",
    ".jenkinsfile",
    "azure-pipelines.yml",
    "Taskfile.yml",

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
    "go.work",

    // Zig
    "zig.mod",
    "build.zig",
    "build.zig.zon",

    // Elixir
    "mix.exs",
    "mix.lock",

    // Swift
    "Package.swift",
    "Package.resolved",

    // .NET,C#
    "global.json",
    "nuget.config",
    "*.csproj",
    "*.sln",
    "*.fsproj",
    "*.vbproj",

    // Haskell
    "stack.yaml",
    "stack.yaml.lock",
    "*.cabal",
    "cabal.project",

    // Terraform
    "main.tf",
    "*.tf",
    ".terragrunt.hcl",
    "terraform.tfstate",
    ".terraform/",

    // Lua
    "*.rockspec",

    // Nix
    "flake.lock",
    "flake.nix",
    "default.nix",
    "shell.nix",

    // Racket
    "info.rkt",

    // Dart
    "pubspec.yaml",

    // Elm
    "elm.json",

    // Julia
    "Project.toml",

    // Emacs
    "Cask",
    "Eask",
    "Eldev",
    "Eldev-local",

    // R
    "DESCRIPTION",

    // Crystal
    "shard.yml",

    // Clojure
    "project.clj",
    ".midje.clj",
    "build.boot",
    "deps.edn",

    // Rails,Grails
    "application.yml",

    // Debian packaging
    "debian/control",

    // Bazel
    "WORKSPACE",

    // Miscellaneous,Editors,Environments
    ".venv/",
    "env/",
    ".vscode/",
    ".idea/",
    ".editorconfig",
    ".code-workspace",
    ".projectile",
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
    C,
    CPP,
    CMake,
    Kotlin,
    Bazel,
    Clojure,
    Dart,
    Elm,
    Julia,
    Crystal,
    R,
    OCamlDune,
    Rails,
    Debian,
    Emacs,
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
            ProjectType::OCamlDune => "ocaml-dune",
            ProjectType::Zig => "zig",
            ProjectType::Terraform => "terraform",
            ProjectType::Lua => "lua",
            ProjectType::Nix => "nix",
            ProjectType::Docker => "docker",
            ProjectType::Git => "git",
            ProjectType::C => "c",
            ProjectType::CPP => "cpp",
            ProjectType::CMake => "cmake",
            ProjectType::Kotlin => "kotlin",
            ProjectType::Bazel => "bazel",
            ProjectType::Clojure => "clojure",
            ProjectType::Dart => "dart",
            ProjectType::Elm => "elm",
            ProjectType::Julia => "julia",
            ProjectType::Crystal => "crystal",
            ProjectType::R => "r",
            ProjectType::Rails => "rails",
            ProjectType::Debian => "debian",
            ProjectType::Emacs => "emacs",
            ProjectType::Unknown => "unknown",
        }
    }
}

#[derive(Debug, Clone)]
pub struct ProjectInfo {
    pub root: PathBuf,
    pub project_types: Vec<ProjectType>,
    pub file_count: usize,
    pub monorepo: bool,
    pub subprojects: Vec<PathBuf>,
}

fn dir_filenames(dir: &Path) -> io::Result<HashSet<String>> {
    let mut names = HashSet::new();
    for entry in dir.read_dir()? {
        let entry = entry?;
        if let Some(name) = entry.file_name().to_str() {
            names.insert(name.to_string());
        }
    }
    Ok(names)
}

pub fn contains_project_file(dir: &Path) -> io::Result<bool> {
    let names = dir_filenames(dir)?;
    Ok(names.iter().any(|n| PROJECT_FILES.contains(n.as_str())))
}

pub fn find_project_root(start: &Path) -> Result<PathBuf, PadoError> {
    for ancestor in start.ancestors() {
        if contains_project_file(ancestor).map_err(PadoError::Io)? {
            return Ok(ancestor.to_path_buf());
        }
    }
    Err(PadoError::NoProjectRoot(start.to_path_buf()))
}

pub fn is_project_root(path: &Path) -> bool {
    contains_project_file(path).unwrap_or(false)
}

fn has(root: &Path, file: &str) -> bool {
    root.join(file).exists()
}

fn has_any(root: &Path, files: &[&str]) -> bool {
    files.iter().any(|f| has(root, f))
}

fn has_any_dir(root: &Path, dirs: &[&str]) -> bool {
    dirs.iter().any(|d| root.join(d).is_dir())
}

fn has_pattern(dir: &Path, pattern: &str) -> bool {
    let glob = match Glob::new(pattern) {
        Ok(g) => g,
        Err(_) => return false,
    };
    let matcher = glob.compile_matcher();

    std::fs::read_dir(dir)
        .ok()
        .map(|entries| {
            entries.flatten().any(|e| {
                e.file_name()
                    .to_str()
                    .map(|name| matcher.is_match(name))
                    .unwrap_or(false)
            })
        })
        .unwrap_or(false)
}

pub fn detect_project_types(root: &Path) -> Vec<ProjectType> {
    let mut types = Vec::new();

    if has(root, "Cargo.toml") {
        types.push(ProjectType::Rust);
    }

    if has(root, "package.json") || has_any(root, &["pnpm-workspace.yaml", "bun.lockb"]) {
        types.push(ProjectType::Node);
    }

    if has_any(
        root,
        &[
            "uv.lock",
            "pyproject.toml",
            "setup.py",
            "requirements.txt",
            "Pipfile",
            "Pipfile.lock",
        ],
    ) {
        types.push(ProjectType::Python);
    }

    if has(root, "go.mod") {
        types.push(ProjectType::Go);
    }

    if has(root, "pom.xml")
        || has_any(
            root,
            &[
                "build.gradle",
                "build.gradle.kts",
                "settings.gradle",
                "settings.gradle.kts",
            ],
        )
    {
        types.push(ProjectType::Java);
    }

    if has(root, "build.sbt") {
        types.push(ProjectType::Scala);
    }

    if has(root, ".bloop/bloop.settings.json") {
        types.push(ProjectType::Kotlin);
    }

    if has(root, "Gemfile") {
        if has(root, "application.yml") {
            types.push(ProjectType::Rails);
        }
        types.push(ProjectType::Ruby);
    }

    if has(root, "composer.json") {
        types.push(ProjectType::PHP);
    }

    if has(root, "mix.exs") {
        types.push(ProjectType::Elixir);
    }

    if has(root, "Package.swift") {
        types.push(ProjectType::Swift);
    }

    if has(root, "global.json") || has_pattern(root, "*.csproj") || has_pattern(root, "*.sln") {
        types.push(ProjectType::Dotnet);
    }

    if has(root, "stack.yaml") || has_pattern(root, "*.cabal") {
        types.push(ProjectType::Haskell);
    }

    if has_any(root, &["dune", "dune-project", "opam"]) {
        types.push(ProjectType::OCaml);
    }

    if has(root, "build.zig") {
        types.push(ProjectType::Zig);
    }

    if has_pattern(root, "*.tf") || has(root, ".terraform") {
        types.push(ProjectType::Terraform);
    }

    if has_pattern(root, "*.rockspec") {
        types.push(ProjectType::Lua);
    }

    if has_any(root, &["flake.nix", "default.nix", "shell.nix"]) {
        types.push(ProjectType::Nix);
    }

    if has_any(
        root,
        &["Dockerfile", "docker-compose.yml", "docker-compose.yaml"],
    ) {
        types.push(ProjectType::Docker);
    }

    if has(root, "WORKSPACE") {
        types.push(ProjectType::Bazel);
    }

    if has_any(root, &["Makefile", "makefile", "GNUMakefile"]) {
        types.push(ProjectType::C);
    }

    if has_any(
        root,
        &[
            "CMakeLists.txt",
            "CMakePresets.json",
            "CMakeUserPresets.json",
        ],
    ) {
        types.push(ProjectType::CMake);
    }

    if has_any(root, &["project.clj", "deps.edn", "build.boot"]) {
        types.push(ProjectType::Clojure);
    }

    if has(root, "pubspec.yaml") {
        types.push(ProjectType::Dart);
    }

    if has(root, "elm.json") {
        types.push(ProjectType::Elm);
    }

    if has(root, "Project.toml") {
        types.push(ProjectType::Julia);
    }

    if has(root, "shard.yml") {
        types.push(ProjectType::Crystal);
    }

    if has(root, "DESCRIPTION") {
        types.push(ProjectType::R);
    }

    if has(root, "debian/control") {
        types.push(ProjectType::Debian);
    }

    if has_any(root, &["Cask", "Eask", "Eldev", "Eldev-local"]) {
        types.push(ProjectType::Emacs);
    }

    if has(root, ".git") {
        types.push(ProjectType::Git);
    }

    types
}

pub fn detect_monorepo(root: &Path) -> bool {
    has_any(
        root,
        &[
            "pnpm-workspace.yaml",
            "lerna.json",
            "turbo.json",
            "nx.json",
            "Cargo.toml",
            "WORKSPACE",
            "go.work",
            "flake.nix",
        ],
    ) && has_any_dir(root, &["packages", "apps", "crates", "modules"])
}

// TODO: Deprecate and make use of detect_project_types
pub fn detect_project_type(root: &Path) -> ProjectType {
    if has(root, "Cargo.toml") {
        return ProjectType::Rust;
    }

    if has(root, "package.json") || has_any(root, &["pnpm-workspace.yaml", "bun.lockb"]) {
        return ProjectType::Node;
    }

    if has_any(
        root,
        &[
            "uv.lock",
            "pyproject.toml",
            "setup.py",
            "requirements.txt",
            "Pipfile",
            "Pipfile.lock",
        ],
    ) {
        return ProjectType::Python;
    }

    if has(root, "go.mod") {
        return ProjectType::Go;
    }

    if has(root, "pom.xml")
        || has_any(
            root,
            &[
                "build.gradle",
                "build.gradle.kts",
                "settings.gradle",
                "settings.gradle.kts",
            ],
        )
    {
        return ProjectType::Java;
    }
    if has(root, "build.sbt") {
        return ProjectType::Scala;
    }
    if has(root, ".bloop/bloop.settings.json") {
        return ProjectType::Kotlin;
    }

    if has(root, "Gemfile") {
        if has(root, "application.yml") {
            return ProjectType::Rails;
        }
        return ProjectType::Ruby;
    }

    if has(root, "composer.json") {
        return ProjectType::PHP;
    }

    if has(root, "mix.exs") {
        return ProjectType::Elixir;
    }

    if has(root, "Package.swift") {
        return ProjectType::Swift;
    }

    if has(root, "global.json") || has_pattern(root, "*.csproj") || has_pattern(root, "*.sln") {
        return ProjectType::Dotnet;
    }

    if has(root, "stack.yaml") || has_pattern(root, "*.cabal") {
        return ProjectType::Haskell;
    }

    if has_any(root, &["dune", "dune-project", "opam"]) {
        return ProjectType::OCaml;
    }

    if has(root, "build.zig") {
        return ProjectType::Zig;
    }

    if has_pattern(root, "*.tf") || has(root, ".terraform") {
        return ProjectType::Terraform;
    }

    if has_pattern(root, "*.rockspec") {
        return ProjectType::Lua;
    }

    if has_any(root, &["flake.nix", "default.nix", "shell.nix"]) {
        return ProjectType::Nix;
    }

    if has_any(
        root,
        &["Dockerfile", "docker-compose.yml", "docker-compose.yaml"],
    ) {
        return ProjectType::Docker;
    }

    if has(root, "WORKSPACE") {
        return ProjectType::Bazel;
    }

    if has_any(root, &["Makefile", "makefile", "GNUMakefile"]) {
        return ProjectType::C;
    }
    if has_any(
        root,
        &[
            "CMakeLists.txt",
            "CMakePresets.json",
            "CMakeUserPresets.json",
        ],
    ) {
        return ProjectType::CMake;
    }

    if has_any(root, &["project.clj", "deps.edn", "build.boot"]) {
        return ProjectType::Clojure;
    }

    if has(root, "pubspec.yaml") {
        return ProjectType::Dart;
    }
    if has(root, "elm.json") {
        return ProjectType::Elm;
    }
    if has(root, "Project.toml") {
        return ProjectType::Julia;
    }
    if has(root, "shard.yml") {
        return ProjectType::Crystal;
    }
    if has(root, "DESCRIPTION") {
        return ProjectType::R;
    }
    if has(root, "debian/control") {
        return ProjectType::Debian;
    }

    if has_any(root, &["Cask", "Eask", "Eldev", "Eldev-local"]) {
        return ProjectType::Emacs;
    }

    if has(root, ".git") {
        return ProjectType::Git;
    }

    ProjectType::Unknown
}

pub fn list_project_files(root: &Path, pattern: Option<&str>) -> Result<Vec<PathBuf>, PadoError> {
    let mut files = Vec::new();

    let matcher = if let Some(pat) = pattern {
        let mut builder = GlobSetBuilder::new();
        builder.add(Glob::new(pat)?);
        Some(builder.build()?)
    } else {
        None
    };

    for result in WalkBuilder::new(root)
        .hidden(false)
        .ignore(true)
        .git_ignore(true)
        .build()
    {
        let entry = result.map_err(|e| {
            io::Error::new(io::ErrorKind::Other, format!("walking project files: {}", e))
        })?;

        if entry.file_type().map(|ft| ft.is_file()).unwrap_or(false) {
            let path = entry.path().to_path_buf();
            if let Some(matcher) = &matcher {
                if let Some(name) = path.file_name().and_then(|n| n.to_str()) {
                    if matcher.is_match(name) {
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

pub fn get_project_info(root: &Path) -> Result<ProjectInfo, PadoError> {
    let project_types = detect_project_types(root);
    let monorepo = detect_monorepo(root);
    let subprojects = if monorepo {
        discover_subprojects(root, 3)
    } else {
        Vec::new()
    };

    let files = list_project_files(root, None)?;
    let file_count = files.len();

    Ok(ProjectInfo {
        root: root.to_path_buf(),
        project_types,
        file_count,
        monorepo,
        subprojects,
    })
}

pub fn discover_subprojects(root: &Path, max_depth: usize) -> Vec<PathBuf> {
    scan_projects(root, max_depth, |p| {
        contains_project_file(p).unwrap_or(false)
    })
    .unwrap_or_default()
}

pub fn discover_projects(search_path: &Path, max_depth: usize) -> Result<Vec<PathBuf>, PadoError> {
    scan_projects(search_path, max_depth, |p| is_project_root(p))
}

fn scan_projects<F>(
    root: &Path,
    max_depth: usize,
    mut predicate: F,
) -> Result<Vec<PathBuf>, PadoError>
where
    F: FnMut(&Path) -> bool,
{
    use std::collections::HashSet;

    let mut found = Vec::new();
    let mut visited = HashSet::new();

    let walker = WalkBuilder::new(root)
        .hidden(true)
        .ignore(true)
        .git_ignore(true)
        .max_depth(Some(max_depth))
        .build();

    for result in walker {
        let entry = result.map_err(|e| {
            std::io::Error::new(
                std::io::ErrorKind::Other,
                format!("walking project tree: {}", e),
            )
        })?;
        let path = entry.path();

        if path == root {
            continue;
        }

        if path.is_dir() && predicate(path) && visited.insert(path.to_path_buf()) {
            found.push(path.to_path_buf());
        }
    }

    Ok(found)
}
