use criterion::{Criterion, black_box, criterion_group, criterion_main};
use pado::*;
use std::fs;
use tempfile::TempDir;

fn setup_nested_project(depth: usize) -> TempDir {
    let temp_dir = TempDir::new().unwrap();
    let mut current = temp_dir.path().to_path_buf();

    for i in 0..depth {
        current = current.join(format!("level{}", i));
        fs::create_dir(&current).unwrap();
    }

    fs::write(temp_dir.path().join("Cargo.toml"), "").unwrap();

    temp_dir
}

fn setup_multi_project_tree(num_projects: usize) -> TempDir {
    let temp_dir = TempDir::new().unwrap();

    for i in 0..num_projects {
        let proj_path = temp_dir.path().join(format!("project{}", i));
        fs::create_dir(&proj_path).unwrap();
        fs::write(proj_path.join("Cargo.toml"), "").unwrap();

        fs::create_dir(proj_path.join("src")).unwrap();
        fs::create_dir(proj_path.join("tests")).unwrap();
    }

    temp_dir
}

fn setup_file_tree(num_files: usize) -> TempDir {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("Cargo.toml"), "").unwrap();

    let src_dir = temp_dir.path().join("src");
    fs::create_dir(&src_dir).unwrap();

    for i in 0..num_files {
        fs::write(src_dir.join(format!("file{}.rs", i)), "").unwrap();
    }

    temp_dir
}

fn bench_find_project_root_shallow(c: &mut Criterion) {
    let temp_dir = setup_nested_project(5);
    let deepest = temp_dir.path().join("level0/level1/level2/level3/level4");

    c.bench_function("find_project_root_shallow_5", |b| {
        b.iter(|| {
            let _ = find_project_root(black_box(&deepest)).unwrap();
        });
    });
}

fn bench_find_project_root_deep(c: &mut Criterion) {
    let temp_dir = setup_nested_project(20);
    let mut deepest = temp_dir.path().to_path_buf();
    for i in 0..20 {
        deepest = deepest.join(format!("level{}", i));
    }

    c.bench_function("find_project_root_deep_20", |b| {
        b.iter(|| {
            let _ = find_project_root(black_box(&deepest)).unwrap();
        });
    });
}

fn bench_is_project_root(c: &mut Criterion) {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("Cargo.toml"), "").unwrap();

    c.bench_function("is_project_root", |b| {
        b.iter(|| {
            let _ = is_project_root(black_box(temp_dir.path()));
        });
    });
}

fn bench_detect_project_type(c: &mut Criterion) {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("Cargo.toml"), "").unwrap();
    fs::write(temp_dir.path().join("package.json"), "{}").unwrap();

    c.bench_function("detect_project_type", |b| {
        b.iter(|| {
            let _ = detect_project_type(black_box(temp_dir.path()));
        });
    });
}

fn bench_build_system_detect(c: &mut Criterion) {
    let temp_dir = TempDir::new().unwrap();
    fs::write(temp_dir.path().join("Cargo.toml"), "").unwrap();

    c.bench_function("build_system_detect", |b| {
        b.iter(|| {
            let _ = BuildSystem::detect(black_box(temp_dir.path()));
        });
    });
}

fn bench_discover_projects_small(c: &mut Criterion) {
    let temp_dir = setup_multi_project_tree(10);

    c.bench_function("discover_projects_10", |b| {
        b.iter(|| {
            let _ = discover_projects(black_box(temp_dir.path()), 3).unwrap();
        });
    });
}

fn bench_discover_projects_medium(c: &mut Criterion) {
    let temp_dir = setup_multi_project_tree(50);

    c.bench_function("discover_projects_50", |b| {
        b.iter(|| {
            let _ = discover_projects(black_box(temp_dir.path()), 3).unwrap();
        });
    });
}

fn bench_list_project_files_small(c: &mut Criterion) {
    let temp_dir = setup_file_tree(50);

    c.bench_function("list_project_files_50", |b| {
        b.iter(|| {
            let _ = list_project_files(black_box(temp_dir.path()), None).unwrap();
        });
    });
}

fn bench_list_project_files_medium(c: &mut Criterion) {
    let temp_dir = setup_file_tree(200);

    c.bench_function("list_project_files_200", |b| {
        b.iter(|| {
            let _ = list_project_files(black_box(temp_dir.path()), None).unwrap();
        });
    });
}

fn bench_list_project_files_with_pattern(c: &mut Criterion) {
    let temp_dir = setup_file_tree(100);

    c.bench_function("list_project_files_with_pattern", |b| {
        b.iter(|| {
            let _ = list_project_files(black_box(temp_dir.path()), Some("*.rs")).unwrap();
        });
    });
}

criterion_group!(
    benches,
    bench_find_project_root_shallow,
    bench_find_project_root_deep,
    bench_is_project_root,
    bench_detect_project_type,
    bench_build_system_detect,
    bench_discover_projects_small,
    bench_discover_projects_medium,
    bench_list_project_files_small,
    bench_list_project_files_medium,
    bench_list_project_files_with_pattern,
);
criterion_main!(benches);
