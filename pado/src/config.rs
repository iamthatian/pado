use libpado::PadoError;
use libpado::PROJECT_FILES;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs;
use std::path::{Path, PathBuf};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GlobalConfig {
    #[serde(default)]
    pub markers: MarkersConfig,

    #[serde(default)]
    pub defaults: DefaultsConfig,

    #[serde(default)]
    pub indexing: IndexingConfig,

    #[serde(default)]
    pub prompt: PromptConfig,

    #[serde(default)]
    pub behavior: BehaviorConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct MarkersConfig {
    #[serde(default)]
    pub additional: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DefaultsConfig {
    #[serde(default = "default_sort_order")]
    pub sort_order: String,

    #[serde(default = "default_recent_limit")]
    pub recent_limit: usize,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct IndexingConfig {
    #[serde(default = "default_indexing_method")]
    pub method: String,

    #[serde(default)]
    pub ignore_patterns: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PromptConfig {
    #[serde(default = "default_prompt_format")]
    pub format: String,

    #[serde(default)]
    pub show_full_path: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BehaviorConfig {
    #[serde(default)]
    pub auto_star_on_add: bool,

    #[serde(default = "default_auto_add")]
    pub auto_add_on_cd: bool,

    #[serde(default = "default_require_confirmation")]
    pub require_confirmation_on_clear: bool,
}

fn default_sort_order() -> String {
    "time".to_string()
}
fn default_recent_limit() -> usize {
    10
}
fn default_indexing_method() -> String {
    "native".to_string()
}
fn default_prompt_format() -> String {
    "{name}:{type}".to_string()
}
fn default_auto_add() -> bool {
    true
}
fn default_require_confirmation() -> bool {
    true
}

impl Default for DefaultsConfig {
    fn default() -> Self {
        Self {
            sort_order: default_sort_order(),
            recent_limit: default_recent_limit(),
        }
    }
}

impl Default for IndexingConfig {
    fn default() -> Self {
        Self {
            method: default_indexing_method(),
            ignore_patterns: Vec::new(),
        }
    }
}

impl Default for PromptConfig {
    fn default() -> Self {
        Self {
            format: default_prompt_format(),
            show_full_path: false,
        }
    }
}

impl Default for BehaviorConfig {
    fn default() -> Self {
        Self {
            auto_star_on_add: false,
            auto_add_on_cd: default_auto_add(),
            require_confirmation_on_clear: default_require_confirmation(),
        }
    }
}

impl Default for GlobalConfig {
    fn default() -> Self {
        Self {
            markers: MarkersConfig::default(),
            defaults: DefaultsConfig::default(),
            indexing: IndexingConfig::default(),
            prompt: PromptConfig::default(),
            behavior: BehaviorConfig::default(),
        }
    }
}

impl GlobalConfig {
    pub fn load() -> Result<GlobalConfig, PadoError> {
        let config_path = get_global_config_path()?;

        if !config_path.exists() {
            return Ok(GlobalConfig::default());
        }

        let contents = fs::read_to_string(&config_path)?;
        let config: GlobalConfig = toml::from_str(&contents)
            .map_err(|e| PadoError::ConfigParse(format!("global config: {}", e)))?;

        Ok(config)
    }

    pub fn save(&self) -> Result<(), PadoError> {
        let config_path = get_global_config_path()?;

        if let Some(parent) = config_path.parent() {
            fs::create_dir_all(parent)?;
        }

        let contents = toml::to_string_pretty(self)
            .map_err(|e| PadoError::SerializationError(format!("global config: {}", e)))?;

        fs::write(&config_path, contents)?;
        Ok(())
    }

    pub fn get_all_markers(&self) -> Vec<&str> {
        let mut markers: Vec<&str> = PROJECT_FILES.iter().copied().collect();
        markers.extend(self.markers.additional.iter().map(|s| s.as_str()));
        markers
    }

    pub fn format_prompt(&self, name: &str, project_type: &str, path: &str) -> String {
        self.prompt
            .format
            .replace("{name}", name)
            .replace("{type}", project_type)
            .replace("{path}", path)
    }
}

pub fn get_global_config_path() -> Result<PathBuf, PadoError> {
    let config_dir = dirs::config_dir().ok_or_else(|| {
        PadoError::InvalidPath("could not determine config directory".to_string())
    })?;

    let pado_config_dir = config_dir.join("pado");
    Ok(pado_config_dir.join("config.toml"))
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProjectConfig {
    #[serde(default)]
    pub commands: HashMap<String, String>,
}

impl ProjectConfig {
    pub fn load(root: &Path) -> Result<Option<ProjectConfig>, PadoError> {
        let config_path = root.join(".pado.toml");
        if !config_path.exists() {
            return Ok(None);
        }

        let contents = fs::read_to_string(&config_path)?;
        let config: ProjectConfig = toml::from_str(&contents)
            .map_err(|e| PadoError::ConfigParse(format!(".pado.toml: {}", e)))?;

        Ok(Some(config))
    }

    pub fn get_command(&self, name: &str) -> Option<&str> {
        self.commands.get(name).map(|s| s.as_str())
    }
}
