use std::path::Path;

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum BuildSystem {
    Cargo,
    Npm,
    Yarn,
    Pnpm,
    Bun,
    Make,
    CMake,
    Maven,
    Gradle,
    Poetry,
    Pip,
    Uv,
    Go,
    Mix,
    Sbt,
    Swift,
    Dotnet,
    Stack,
    Cabal,
    Dune,
    Zig,
    Terraform,
    Luarocks,
    Nix,
    Unknown,
}

impl BuildSystem {
    pub fn detect(root: &Path) -> Self {
        use BuildSystem::*;

        // NOTE: Order matters
        static FIXED: &[(BuildSystem, &[&str])] = &[
            (Cargo, &["Cargo.toml"]),
            (Uv, &["uv.lock"]),
            (Poetry, &["pyproject.toml"]),
            (Pip, &["requirements.txt", "setup.py"]),
            (Maven, &["pom.xml"]),
            (Gradle, &["build.gradle", "build.gradle.kts"]),
            (Sbt, &["build.sbt"]),
            (Go, &["go.mod"]),
            (Mix, &["mix.exs"]),
            (Swift, &["Package.swift"]),
            (Stack, &["stack.yaml"]),
            (Dune, &["dune", "dune-project"]),
            (Zig, &["build.zig"]),
            (Nix, &["flake.nix", "default.nix", "shell.nix"]),
            (CMake, &["CMakeLists.txt"]),
            (Make, &["Makefile", "makefile"]),
        ];

        for (system, files) in FIXED {
            if files.iter().any(|f| root.join(f).exists()) {
                return system.clone();
            }
        }

        if root.join("package.json").exists() {
            return if root.join("bun.lockb").exists() {
                Bun
            } else if root.join("pnpm-lock.yaml").exists() {
                Pnpm
            } else if root.join("yarn.lock").exists() {
                Yarn
            } else {
                Npm
            };
        }

        if Self::has_extension(root, &[".csproj", ".fsproj", ".sln"])
            || root.join("global.json").exists()
        {
            return Dotnet;
        }

        if Self::has_extension(root, &[".cabal"]) {
            return Cabal;
        }

        if Self::has_extension(root, &[".tf"]) || root.join(".terraform").exists() {
            return Terraform;
        }

        if Self::has_extension(root, &[".rockspec"]) {
            return Luarocks;
        }

        Unknown
    }

    fn has_extension(root: &Path, exts: &[&str]) -> bool {
        std::fs::read_dir(root)
            .ok()
            .map(|entries| {
                entries.filter_map(|e| e.ok()).any(|entry| {
                    entry
                        .file_name()
                        .to_str()
                        .map(|name| exts.iter().any(|ext| name.ends_with(ext)))
                        .unwrap_or(false)
                })
            })
            .unwrap_or(false)
    }

    pub fn build_command(&self) -> Option<&str> {
        match self {
            BuildSystem::Cargo => Some("cargo build"),
            BuildSystem::Npm => Some("npm run build"),
            BuildSystem::Yarn => Some("yarn build"),
            BuildSystem::Pnpm => Some("pnpm build"),
            BuildSystem::Bun => Some("bun build"),
            BuildSystem::Make => Some("make"),
            BuildSystem::CMake => Some("cmake --build build"),
            BuildSystem::Maven => Some("mvn compile"),
            BuildSystem::Gradle => Some("gradle build"),
            BuildSystem::Poetry => Some("poetry build"),
            BuildSystem::Uv => Some("uv build"),
            BuildSystem::Go => Some("go build"),
            BuildSystem::Mix => Some("mix compile"),
            BuildSystem::Sbt => Some("sbt compile"),
            BuildSystem::Swift => Some("swift build"),
            BuildSystem::Dotnet => Some("dotnet build"),
            BuildSystem::Stack => Some("stack build"),
            BuildSystem::Cabal => Some("cabal build"),
            BuildSystem::Dune => Some("dune build"),
            BuildSystem::Zig => Some("zig build"),
            BuildSystem::Terraform => Some("terraform plan"),
            BuildSystem::Nix => Some("nix build"),
            BuildSystem::Pip | BuildSystem::Luarocks | BuildSystem::Unknown => None,
        }
    }

    pub fn test_command(&self) -> Option<&str> {
        match self {
            BuildSystem::Cargo => Some("cargo test"),
            BuildSystem::Npm => Some("npm test"),
            BuildSystem::Yarn => Some("yarn test"),
            BuildSystem::Pnpm => Some("pnpm test"),
            BuildSystem::Bun => Some("bun test"),
            BuildSystem::Make => Some("make test"),
            BuildSystem::CMake => Some("ctest"),
            BuildSystem::Maven => Some("mvn test"),
            BuildSystem::Gradle => Some("gradle test"),
            BuildSystem::Poetry => Some("poetry run pytest"),
            BuildSystem::Uv => Some("uv run pytest"),
            BuildSystem::Pip => Some("pytest"),
            BuildSystem::Go => Some("go test ./..."),
            BuildSystem::Mix => Some("mix test"),
            BuildSystem::Sbt => Some("sbt test"),
            BuildSystem::Swift => Some("swift test"),
            BuildSystem::Dotnet => Some("dotnet test"),
            BuildSystem::Stack => Some("stack test"),
            BuildSystem::Cabal => Some("cabal test"),
            BuildSystem::Dune => Some("dune runtest"),
            BuildSystem::Zig => Some("zig build test"),
            BuildSystem::Terraform => Some("terraform validate"),
            BuildSystem::Nix => Some("nix flake check"),
            BuildSystem::Luarocks | BuildSystem::Unknown => None,
        }
    }

    pub fn run_command(&self) -> Option<&str> {
        match self {
            BuildSystem::Cargo => Some("cargo run"),
            BuildSystem::Npm => Some("npm start"),
            BuildSystem::Yarn => Some("yarn start"),
            BuildSystem::Pnpm => Some("pnpm start"),
            BuildSystem::Bun => Some("bun start"),
            BuildSystem::Go => Some("go run ."),
            BuildSystem::Poetry => Some("poetry run python"),
            BuildSystem::Uv => Some("uv run python"),
            BuildSystem::Pip => Some("python"),
            BuildSystem::Maven => Some("mvn exec:java"),
            BuildSystem::Gradle => Some("gradle run"),
            BuildSystem::Mix => Some("mix run"),
            BuildSystem::Sbt => Some("sbt run"),
            BuildSystem::Swift => Some("swift run"),
            BuildSystem::Dotnet => Some("dotnet run"),
            BuildSystem::Stack => Some("stack run"),
            BuildSystem::Cabal => Some("cabal run"),
            BuildSystem::Dune => Some("dune exec"),
            BuildSystem::Zig => Some("zig build run"),
            BuildSystem::Terraform => Some("terraform apply"),
            BuildSystem::Nix => Some("nix run"),
            BuildSystem::Make
            | BuildSystem::CMake
            | BuildSystem::Luarocks
            | BuildSystem::Unknown => None,
        }
    }

    pub fn outdated_command(&self) -> Option<(&'static str, &'static [&'static str])> {
        match self {
            BuildSystem::Cargo => Some(("cargo", &["outdated"])),
            BuildSystem::Npm => Some(("npm", &["outdated"])),
            BuildSystem::Yarn => Some(("yarn", &["outdated"])),
            BuildSystem::Pnpm => Some(("pnpm", &["outdated"])),
            BuildSystem::Bun => Some(("bun", &["outdated"])),
            BuildSystem::Pip => Some(("pip", &["list", "--outdated"])),
            BuildSystem::Poetry => Some(("poetry", &["show", "-o"])),
            BuildSystem::Uv => Some(("uv", &["pip", "list", "--outdated"])),
            BuildSystem::Maven => Some(("mvn", &["versions:display-dependency-updates"])),
            BuildSystem::Gradle => Some(("gradle", &["dependencyUpdates"])),
            BuildSystem::Sbt => Some(("sbt", &["dependencyUpdates"])),
            BuildSystem::Go => Some(("go", &["list", "-u", "-m", "all"])),
            BuildSystem::Mix => Some(("mix", &["hex.outdated"])),
            BuildSystem::Swift => Some(("swift", &["package", "update", "--dry-run"])),
            BuildSystem::Dotnet => Some(("dotnet", &["list", "package", "--outdated"])),
            BuildSystem::Stack => Some(("stack", &["list-dependencies"])),
            BuildSystem::Cabal => Some(("cabal", &["outdated"])),
            BuildSystem::Dune => Some(("opam", &["upgrade", "--dry-run"])),
            BuildSystem::Zig => Some(("zig", &["build", "update-check"])),
            BuildSystem::Terraform => Some(("terraform", &["providers", "lock", "-check"])),
            BuildSystem::Nix => Some(("nix", &["flake", "update", "--dry-run"])),
            BuildSystem::Luarocks => Some(("luarocks", &["list", "--outdated"])),
            BuildSystem::Make | BuildSystem::CMake | BuildSystem::Unknown => None,
        }
    }
}
