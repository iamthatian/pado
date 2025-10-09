pub mod error;
pub mod config;
pub mod project;
pub mod tracking;
pub mod build;
pub mod stats;

pub use error::ParkourError;
pub use config::{GlobalConfig, ProjectConfig, get_global_config_path};
pub use project::{
    PROJECT_FILES, ProjectType, ProjectInfo,
    contains_project_file, contains_project_file_with_config,
    find_project_root, find_project_root_with_config,
    is_project_root, detect_project_type,
    list_project_files, get_project_info, discover_projects,
    glob_match,
};
pub use tracking::{
    TrackedProject, ProjectList,
    get_project_list_path, load_project_list, save_project_list,
};
pub use build::BuildSystem;
pub use stats::{
    LanguageStats, ProjectStats, Contributor, GitStats,
    get_language_stats, get_git_stats, format_time_ago, create_percentage_bar,
};
