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
        if root.join("Cargo.toml").exists() {
            return BuildSystem::Cargo;
        }

        if root.join("package.json").exists() {
            if root.join("bun.lockb").exists() {
                return BuildSystem::Bun;
            } else if root.join("pnpm-lock.yaml").exists() {
                return BuildSystem::Pnpm;
            } else if root.join("yarn.lock").exists() {
                return BuildSystem::Yarn;
            } else {
                return BuildSystem::Npm;
            }
        }

        if root.join("uv.lock").exists() {
            return BuildSystem::Uv;
        }
        if root.join("pyproject.toml").exists() {
            return BuildSystem::Poetry;
        }
        if root.join("requirements.txt").exists() || root.join("setup.py").exists() {
            return BuildSystem::Pip;
        }

        if root.join("pom.xml").exists() {
            return BuildSystem::Maven;
        }
        if root.join("build.gradle").exists() || root.join("build.gradle.kts").exists() {
            return BuildSystem::Gradle;
        }

        if root.join("build.sbt").exists() {
            return BuildSystem::Sbt;
        }

        if root.join("go.mod").exists() {
            return BuildSystem::Go;
        }

        if root.join("mix.exs").exists() {
            return BuildSystem::Mix;
        }

        if root.join("Package.swift").exists() {
            return BuildSystem::Swift;
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
            return BuildSystem::Dotnet;
        }

        if root.join("stack.yaml").exists() {
            return BuildSystem::Stack;
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
                            .map(|name| name.ends_with(".cabal"))
                            .unwrap_or(false)
                    })
            })
            .unwrap_or(false)
            || root.join("cabal.project").exists()
        {
            return BuildSystem::Cabal;
        }

        if root.join("dune").exists() || root.join("dune-project").exists() {
            return BuildSystem::Dune;
        }

        if root.join("build.zig").exists() {
            return BuildSystem::Zig;
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
            return BuildSystem::Terraform;
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
            return BuildSystem::Luarocks;
        }

        if root.join("flake.nix").exists()
            || root.join("default.nix").exists()
            || root.join("shell.nix").exists()
        {
            return BuildSystem::Nix;
        }

        if root.join("CMakeLists.txt").exists() {
            return BuildSystem::CMake;
        }
        if root.join("Makefile").exists() || root.join("makefile").exists() {
            return BuildSystem::Make;
        }

        BuildSystem::Unknown
    }

    pub fn build_command(&self) -> Option<&str> {
        match self {
            BuildSystem::Cargo => Some("cargo build"),
            BuildSystem::Npm => Some("npm run build"),
            BuildSystem::Yarn => Some("yarn build"),
            BuildSystem::Pnpm => Some("pnpm build"),
            BuildSystem::Bun => Some("bun run build"),
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
            BuildSystem::Bun => Some("bun run start"),
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
            BuildSystem::Make | BuildSystem::CMake | BuildSystem::Luarocks | BuildSystem::Unknown => None,
        }
    }
}
