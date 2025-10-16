use libpado::*;
use std::fs;
use tempfile::TempDir;

#[test]
fn test_find_project_root_basic() {
    let temp_dir = TempDir::new().unwrap();
    let root = temp_dir.path();

    fs::write(root.join("Cargo.toml"), "").unwrap();

    let subdir = root.join("src").join("deep").join("nested");
    fs::create_dir_all(&subdir).unwrap();

    let found_root = find_project_root(&subdir).unwrap();
    assert_eq!(
        found_root.canonicalize().unwrap(),
        root.canonicalize().unwrap()
    );
}

#[test]
fn test_find_project_root_no_project() {
    let temp_dir = TempDir::new().unwrap();
    let path = temp_dir.path().join("not_a_project");
    fs::create_dir_all(&path).unwrap();

    let result = find_project_root(&path);
    assert!(result.is_err());
}

#[test]
fn test_find_project_root_different_markers() {
    let temp_dir = TempDir::new().unwrap();
    fs::create_dir(temp_dir.path().join(".git")).unwrap();
    let subdir = temp_dir.path().join("src");
    fs::create_dir(&subdir).unwrap();
    let found = find_project_root(&subdir).unwrap();
    assert_eq!(
        found.canonicalize().unwrap(),
        temp_dir.path().canonicalize().unwrap()
    );

    let temp_dir2 = TempDir::new().unwrap();
    fs::write(temp_dir2.path().join("package.json"), "{}").unwrap();
    let subdir2 = temp_dir2.path().join("src");
    fs::create_dir(&subdir2).unwrap();
    let found2 = find_project_root(&subdir2).unwrap();
    assert_eq!(
        found2.canonicalize().unwrap(),
        temp_dir2.path().canonicalize().unwrap()
    );
}

#[test]
fn test_find_project_root_nested_projects() {
    let temp_dir = TempDir::new().unwrap();
    let outer_root = temp_dir.path();
    let inner_root = outer_root.join("inner-project");

    fs::write(outer_root.join("Cargo.toml"), "[package]\nname = \"outer\"").unwrap();

    fs::create_dir(&inner_root).unwrap();
    fs::write(inner_root.join("package.json"), "{}").unwrap();

    let inner_subdir = inner_root.join("src");
    fs::create_dir(&inner_subdir).unwrap();

    let found = find_project_root(&inner_subdir).unwrap();
    assert_eq!(
        found.canonicalize().unwrap(),
        inner_root.canonicalize().unwrap()
    );
}

#[test]
fn test_find_project_root_with_projectile_marker() {
    let temp_dir = TempDir::new().unwrap();
    let root = temp_dir.path();

    fs::write(root.join(".projectile"), "").unwrap();

    let subdir = root.join("src");
    fs::create_dir(&subdir).unwrap();

    let found = find_project_root(&subdir).unwrap();
    assert_eq!(found.canonicalize().unwrap(), root.canonicalize().unwrap());
}

#[test]
fn test_is_project_root_true() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("Cargo.toml"), "").unwrap();

    assert!(is_project_root(temp_dir.path()));
}

#[test]
fn test_is_project_root_false() {
    let temp_dir = TempDir::new().unwrap();
    assert!(!is_project_root(temp_dir.path()));
}

#[test]
fn test_is_project_root_git_directory() {
    let temp_dir = TempDir::new().unwrap();
    fs::create_dir(temp_dir.path().join(".git")).unwrap();

    assert!(is_project_root(temp_dir.path()));
}

#[test]
fn test_detect_project_type_rust() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("Cargo.toml"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Rust);
}

#[test]
fn test_detect_project_type_node() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("package.json"), "{}").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Node);
}

#[test]
fn test_detect_project_type_python() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("pyproject.toml"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Python);
}

#[test]
fn test_detect_project_type_go() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("go.mod"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Go);
}

#[test]
fn test_detect_project_type_java() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("pom.xml"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Java);
}

#[test]
fn test_detect_project_type_git_only() {
    let temp_dir = TempDir::new().unwrap();
    fs::create_dir(temp_dir.path().join(".git")).unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Git);
}

#[test]
fn test_detect_project_type_unknown() {
    let temp_dir = TempDir::new().unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Unknown);
}

#[test]
fn test_detect_project_type_python_uv() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("uv.lock"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Python);
}

#[test]
fn test_detect_project_type_elixir() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("mix.exs"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Elixir);
}

#[test]
fn test_detect_project_type_scala() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("build.sbt"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Scala);
}

#[test]
fn test_detect_project_type_swift() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("Package.swift"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Swift);
}

#[test]
fn test_detect_project_type_dotnet() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("app.csproj"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Dotnet);
}

#[test]
fn test_detect_project_type_haskell_stack() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("stack.yaml"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Haskell);
}

#[test]
fn test_detect_project_type_ocaml() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("dune-project"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::OCaml);
}

#[test]
fn test_detect_project_type_zig() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("build.zig"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Zig);
}

#[test]
fn test_detect_project_type_terraform() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("main.tf"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Terraform);
}

#[test]
fn test_detect_project_type_nix() {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("flake.nix"), "").unwrap();

    let proj_type = detect_project_type(temp_dir.path());
    assert_eq!(proj_type, ProjectType::Nix);
}

#[test]
fn test_discover_projects_basic() {
    let temp_dir = TempDir::new().unwrap();
    let search_path = temp_dir.path();

    let proj1 = search_path.join("project1");
    let proj2 = search_path.join("project2");
    fs::create_dir(&proj1).unwrap();
    fs::create_dir(&proj2).unwrap();
    fs::write(proj1.join("Cargo.toml"), "").unwrap();
    fs::write(proj2.join("package.json"), "{}").unwrap();

    let projects = discover_projects(search_path, 2).unwrap();
    assert_eq!(projects.len(), 2);
}

#[test]
fn test_discover_projects_respects_max_depth() {
    let temp_dir = TempDir::new().unwrap();
    let search_path = temp_dir.path();

    let shallow = search_path.join("shallow");
    fs::create_dir(&shallow).unwrap();
    fs::write(shallow.join("Cargo.toml"), "").unwrap();

    let deep = search_path.join("level1").join("level2").join("deep");
    fs::create_dir_all(&deep).unwrap();
    fs::write(deep.join("package.json"), "{}").unwrap();

    let projects = discover_projects(search_path, 1).unwrap();
    assert_eq!(projects.len(), 1);

    let projects = discover_projects(search_path, 3).unwrap();
    assert_eq!(projects.len(), 2);
}

#[test]
fn test_discover_projects_skips_hidden() {
    let temp_dir = TempDir::new().unwrap();
    let search_path = temp_dir.path();

    let visible = search_path.join("visible");
    fs::create_dir(&visible).unwrap();
    fs::write(visible.join("Cargo.toml"), "").unwrap();

    let hidden = search_path.join(".hidden");
    fs::create_dir(&hidden).unwrap();
    fs::write(hidden.join("package.json"), "{}").unwrap();

    let projects = discover_projects(search_path, 2).unwrap();
    assert_eq!(projects.len(), 1);
    assert!(projects[0].ends_with("visible"));
}

#[test]
fn test_discover_projects_no_duplicates() {
    let temp_dir = TempDir::new().unwrap();
    let search_path = temp_dir.path();

    let proj = search_path.join("project");
    fs::create_dir(&proj).unwrap();
    fs::write(proj.join("Cargo.toml"), "").unwrap();
    fs::write(proj.join("package.json"), "{}").unwrap();
    fs::create_dir(proj.join(".git")).unwrap();

    let projects = discover_projects(search_path, 2).unwrap();
    assert_eq!(projects.len(), 1);
}

#[test]
fn test_list_project_files_basic() {
    let temp_dir = TempDir::new().unwrap();
    let root = temp_dir.path();

    fs::write(root.join("file1.txt"), "").unwrap();
    fs::write(root.join("file2.rs"), "").unwrap();
    fs::create_dir(root.join("src")).unwrap();
    fs::write(root.join("src").join("lib.rs"), "").unwrap();

    let files = list_project_files(root, None).unwrap();
    assert!(files.len() >= 3);
}

#[test]
fn test_list_project_files_with_pattern() {
    let temp_dir = TempDir::new().unwrap();
    let root = temp_dir.path();

    fs::write(root.join("file1.txt"), "").unwrap();
    fs::write(root.join("file2.rs"), "").unwrap();
    fs::write(root.join("file3.rs"), "").unwrap();

    let files = list_project_files(root, Some("*.rs")).unwrap();
    assert!(
        files
            .iter()
            .any(|f| f.to_str().unwrap().contains("file2.rs"))
    );
    assert!(
        files
            .iter()
            .any(|f| f.to_str().unwrap().contains("file3.rs"))
    );
}

#[test]
fn test_list_project_files_respects_gitignore() {
    let temp_dir = TempDir::new().unwrap();
    let root = temp_dir.path();

    fs::create_dir(root.join(".git")).unwrap();

    fs::write(root.join(".gitignore"), "ignored.txt\ntarget/\n").unwrap();

    fs::write(root.join("included.txt"), "").unwrap();
    fs::write(root.join("ignored.txt"), "").unwrap();
    fs::create_dir(root.join("target")).unwrap();
    fs::write(root.join("target").join("debug.txt"), "").unwrap();

    let files = list_project_files(root, None).unwrap();
    let file_names: Vec<_> = files
        .iter()
        .filter_map(|p| p.file_name().and_then(|n| n.to_str()))
        .collect();

    assert!(file_names.contains(&"included.txt"));
    assert!(!file_names.contains(&"ignored.txt"));
    assert!(!file_names.contains(&"debug.txt"));
}
