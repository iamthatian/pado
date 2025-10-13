use pado::*;
use tempfile::TempDir;

#[test]
fn test_global_config_load_default() {
    let config = GlobalConfig::default();
    assert_eq!(config.defaults.sort_order, "time");
    assert_eq!(config.defaults.recent_limit, 10);
}

#[test]
fn test_project_config_load_missing() {
    let temp_dir = TempDir::new().unwrap();
    let result = ProjectConfig::load(temp_dir.path());
    assert!(result.is_ok());
    assert!(result.unwrap().is_none());
}

#[test]
fn test_prompt_format() {
    let config = GlobalConfig::default();
    let output = config.format_prompt("myproject", "rust", "/path/to/myproject");
    assert_eq!(output, "myproject:rust");

    let mut custom_config = GlobalConfig::default();
    custom_config.prompt.format = "[{type}] {name}".to_string();
    let output = custom_config.format_prompt("myproject", "rust", "/path/to/myproject");
    assert_eq!(output, "[rust] myproject");
}

#[test]
fn test_custom_markers() {
    let mut config = GlobalConfig::default();
    config.markers.additional.push(".mymarker".to_string());

    let markers = config.get_all_markers();
    assert!(markers.iter().any(|&m| m == ".mymarker"));
    assert!(markers.iter().any(|&m| m == ".git"));
    assert!(markers.iter().any(|&m| m == "Cargo.toml"));
}
