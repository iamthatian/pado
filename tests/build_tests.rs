use parkour::*;
use std::fs;
use tempfile::TempDir;

#[test]
fn test_build_system_detection() {
    let temp_dir = TempDir::new().unwrap();

    std::fs::write(temp_dir.path().join("Cargo.toml"), "").unwrap();
    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Cargo);

    assert_eq!(build_sys.build_command(), Some("cargo build"));
    assert_eq!(build_sys.test_command(), Some("cargo test"));
    assert_eq!(build_sys.run_command(), Some("cargo run"));
}

#[test]
fn test_build_system_priority_cargo_over_make() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("Cargo.toml"), "").unwrap();
    fs::write(temp_dir.path().join("Makefile"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Cargo);
}

#[test]
fn test_build_system_priority_cargo_over_npm() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("Cargo.toml"), "").unwrap();
    fs::write(temp_dir.path().join("package.json"), "{}").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Cargo);
}

#[test]
fn test_build_system_priority_node_lockfiles() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("package.json"), "{}").unwrap();

    assert_eq!(BuildSystem::detect(temp_dir.path()), BuildSystem::Npm);

    fs::write(temp_dir.path().join("yarn.lock"), "").unwrap();
    assert_eq!(BuildSystem::detect(temp_dir.path()), BuildSystem::Yarn);

    fs::write(temp_dir.path().join("pnpm-lock.yaml"), "").unwrap();
    assert_eq!(BuildSystem::detect(temp_dir.path()), BuildSystem::Pnpm);

    fs::write(temp_dir.path().join("bun.lockb"), "").unwrap();
    assert_eq!(BuildSystem::detect(temp_dir.path()), BuildSystem::Bun);
}

#[test]
fn test_build_system_priority_python_poetry_over_pip() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("pyproject.toml"), "").unwrap();
    fs::write(temp_dir.path().join("requirements.txt"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Poetry);
}

#[test]
fn test_build_system_priority_java_maven_over_gradle() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("pom.xml"), "").unwrap();
    fs::write(temp_dir.path().join("build.gradle"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Maven);
}

#[test]
fn test_build_system_priority_cmake_over_make() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("CMakeLists.txt"), "").unwrap();
    fs::write(temp_dir.path().join("Makefile"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::CMake);
}

#[test]
fn test_build_system_priority_uv_over_poetry() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("uv.lock"), "").unwrap();
    fs::write(temp_dir.path().join("pyproject.toml"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Uv);
}

#[test]
fn test_build_system_priority_stack_over_cabal() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("stack.yaml"), "").unwrap();
    fs::write(temp_dir.path().join("app.cabal"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Stack);
}

#[test]
fn test_build_system_uv() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("uv.lock"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Uv);
    assert_eq!(build_sys.build_command(), Some("uv build"));
    assert_eq!(build_sys.test_command(), Some("uv run pytest"));
    assert_eq!(build_sys.run_command(), Some("uv run python"));
}

#[test]
fn test_build_system_mix() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("mix.exs"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Mix);
    assert_eq!(build_sys.build_command(), Some("mix compile"));
    assert_eq!(build_sys.test_command(), Some("mix test"));
    assert_eq!(build_sys.run_command(), Some("mix run"));
}

#[test]
fn test_build_system_sbt() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("build.sbt"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Sbt);
    assert_eq!(build_sys.build_command(), Some("sbt compile"));
    assert_eq!(build_sys.test_command(), Some("sbt test"));
    assert_eq!(build_sys.run_command(), Some("sbt run"));
}

#[test]
fn test_build_system_swift() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("Package.swift"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Swift);
    assert_eq!(build_sys.build_command(), Some("swift build"));
    assert_eq!(build_sys.test_command(), Some("swift test"));
    assert_eq!(build_sys.run_command(), Some("swift run"));
}

#[test]
fn test_build_system_dotnet() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("app.csproj"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Dotnet);
    assert_eq!(build_sys.build_command(), Some("dotnet build"));
    assert_eq!(build_sys.test_command(), Some("dotnet test"));
    assert_eq!(build_sys.run_command(), Some("dotnet run"));
}

#[test]
fn test_build_system_stack() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("stack.yaml"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Stack);
    assert_eq!(build_sys.build_command(), Some("stack build"));
    assert_eq!(build_sys.test_command(), Some("stack test"));
    assert_eq!(build_sys.run_command(), Some("stack run"));
}

#[test]
fn test_build_system_dune() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("dune-project"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Dune);
    assert_eq!(build_sys.build_command(), Some("dune build"));
    assert_eq!(build_sys.test_command(), Some("dune runtest"));
    assert_eq!(build_sys.run_command(), Some("dune exec"));
}

#[test]
fn test_build_system_zig() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("build.zig"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Zig);
    assert_eq!(build_sys.build_command(), Some("zig build"));
    assert_eq!(build_sys.test_command(), Some("zig build test"));
    assert_eq!(build_sys.run_command(), Some("zig build run"));
}

#[test]
fn test_build_system_terraform() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("main.tf"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Terraform);
    assert_eq!(build_sys.build_command(), Some("terraform plan"));
    assert_eq!(build_sys.test_command(), Some("terraform validate"));
    assert_eq!(build_sys.run_command(), Some("terraform apply"));
}

#[test]
fn test_build_system_nix() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("flake.nix"), "").unwrap();

    let build_sys = BuildSystem::detect(temp_dir.path());
    assert_eq!(build_sys, BuildSystem::Nix);
    assert_eq!(build_sys.build_command(), Some("nix build"));
    assert_eq!(build_sys.test_command(), Some("nix flake check"));
    assert_eq!(build_sys.run_command(), Some("nix run"));
}
