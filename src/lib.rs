pub mod build;
pub mod config;
pub mod error;
pub mod project;
pub mod stats;
pub mod tracking;

pub use build::BuildSystem;
pub use config::{GlobalConfig, ProjectConfig, get_global_config_path};
pub use error::ParkourError;
pub use project::{
    PROJECT_FILES, ProjectInfo, ProjectType, contains_project_file,
    contains_project_file_with_config, detect_project_type, discover_projects, find_project_root,
    find_project_root_with_config, get_project_info, glob_match, is_project_root,
    list_project_files,
};
pub use stats::{
    Contributor, GitStats, LanguageStats, ProjectStats, create_percentage_bar, format_time_ago,
    get_git_stats, get_language_stats,
};
pub use tracking::{
    ProjectList, TrackedProject, get_project_list_path, load_project_list, save_project_list,
};
