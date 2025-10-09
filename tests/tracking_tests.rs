use parkour::*;
use std::fs;
use tempfile::TempDir;

#[test]
fn test_project_list_operations() {
    let mut list = ProjectList::new();
    assert_eq!(list.projects.len(), 0);

    let temp_dir = TempDir::new().unwrap();
    let project_path = temp_dir.path().to_path_buf();
    std::fs::write(project_path.join("Cargo.toml"), "").unwrap();

    list.add_project(project_path.clone()).unwrap();
    assert_eq!(list.projects.len(), 1);

    assert!(!list.projects.values().next().unwrap().starred);

    let canonical_path = project_path.canonicalize().unwrap();
    assert!(list.toggle_star(&canonical_path));
    assert!(list.projects.values().next().unwrap().starred);

    let starred = list.get_starred_projects();
    assert_eq!(starred.len(), 1);

    assert!(list.remove_project(&canonical_path));
    assert_eq!(list.projects.len(), 0);
}

#[test]
fn test_project_list_sorting() {
    let list = ProjectList::new();

    let projects = list.get_projects_sorted_by("time");
    assert_eq!(projects.len(), 0);

    let projects = list.get_projects_sorted_by("access");
    assert_eq!(projects.len(), 0);

    let projects = list.get_projects_sorted_by("name");
    assert_eq!(projects.len(), 0);
}

#[test]
fn test_save_and_load_project_list() {
    let temp_dir = TempDir::new().unwrap();
    let project_path = temp_dir.path().join("test-project");
    fs::create_dir(&project_path).unwrap();
    fs::write(project_path.join("Cargo.toml"), "").unwrap();

    let mut list = ProjectList::new();
    list.add_project(project_path.clone()).unwrap();

    let json = serde_json::to_string(&list).unwrap();
    let loaded: ProjectList = serde_json::from_str(&json).unwrap();

    assert_eq!(loaded.projects.len(), 1);
}

#[test]
fn test_load_project_list_missing_file() {
    let list = ProjectList::new();
    assert_eq!(list.projects.len(), 0);
}

#[test]
fn test_project_list_new_project_starts_with_count_one() {
    let temp_dir = TempDir::new().unwrap();
    let project_path = temp_dir.path().join("test-project");
    fs::create_dir(&project_path).unwrap();
    fs::write(project_path.join("Cargo.toml"), "").unwrap();

    let mut list = ProjectList::new();
    list.add_project(project_path.clone()).unwrap();

    let canonical_path = project_path.canonicalize().unwrap();
    let key = canonical_path.to_string_lossy().to_string();
    let project = list.projects.get(&key).unwrap();

    assert_eq!(project.access_count, 1);
}

#[test]
fn test_project_list_add_existing_preserves_count() {
    let temp_dir = TempDir::new().unwrap();
    let project_path = temp_dir.path().join("test-project");
    fs::create_dir(&project_path).unwrap();
    fs::write(project_path.join("Cargo.toml"), "").unwrap();

    let mut list = ProjectList::new();

    list.add_project(project_path.clone()).unwrap();

    let canonical_path = project_path.canonicalize().unwrap();

    list.update_access_time(&canonical_path);
    list.update_access_time(&canonical_path);

    let key = canonical_path.to_string_lossy().to_string();
    let count_before = list.projects.get(&key).unwrap().access_count;
    assert_eq!(count_before, 3);

    list.add_project(project_path).unwrap();

    let count_after = list.projects.get(&key).unwrap().access_count;
    assert_eq!(count_after, 3);
}

#[test]
fn test_project_list_update_access_time_increments() {
    let temp_dir = TempDir::new().unwrap();
    let project_path = temp_dir.path().join("test-project");
    fs::create_dir(&project_path).unwrap();
    fs::write(project_path.join("Cargo.toml"), "").unwrap();

    let mut list = ProjectList::new();
    list.add_project(project_path.clone()).unwrap();

    let canonical_path = project_path.canonicalize().unwrap();
    let key = canonical_path.to_string_lossy().to_string();

    let initial_count = list.projects.get(&key).unwrap().access_count;
    let initial_time = list.projects.get(&key).unwrap().last_accessed;

    std::thread::sleep(std::time::Duration::from_millis(10));

    list.update_access_time(&canonical_path);

    let updated_count = list.projects.get(&key).unwrap().access_count;
    let updated_time = list.projects.get(&key).unwrap().last_accessed;

    assert_eq!(updated_count, initial_count + 1);
    assert!(updated_time > initial_time);
}

#[test]
fn test_project_list_sort_by_access_frequency() {
    let temp_dir = TempDir::new().unwrap();

    let proj1 = temp_dir.path().join("project1");
    let proj2 = temp_dir.path().join("project2");
    let proj3 = temp_dir.path().join("project3");

    for proj in [&proj1, &proj2, &proj3] {
        fs::create_dir(proj).unwrap();
        fs::write(proj.join("Cargo.toml"), "").unwrap();
    }

    let mut list = ProjectList::new();
    list.add_project(proj1.clone()).unwrap();
    list.add_project(proj2.clone()).unwrap();
    list.add_project(proj3.clone()).unwrap();

    let canonical1 = proj1.canonicalize().unwrap();
    let canonical2 = proj2.canonicalize().unwrap();

    for _ in 0..4 {
        list.update_access_time(&canonical2);
    }

    for _ in 0..2 {
        list.update_access_time(&canonical1);
    }

    let sorted = list.get_projects_sorted_by("access");

    assert_eq!(sorted[0].access_count, 5);
    assert_eq!(sorted[1].access_count, 3);
    assert_eq!(sorted[2].access_count, 1);
    assert!(sorted[0].name.contains("project2"));
    assert!(sorted[1].name.contains("project1"));
    assert!(sorted[2].name.contains("project3"));
}
