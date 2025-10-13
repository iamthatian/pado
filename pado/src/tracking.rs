use chrono::{DateTime, Utc};
use libpado::PadoError;
use libpado::detect_project_type;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs;
use std::path::{Path, PathBuf};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TrackedProject {
    pub path: PathBuf,
    pub name: String,
    #[serde(rename = "type")]
    pub project_type: String,
    pub last_accessed: DateTime<Utc>,
    #[serde(default)]
    pub access_count: usize,
    #[serde(default)]
    pub starred: bool,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProjectList {
    pub projects: HashMap<String, TrackedProject>,
}

impl ProjectList {
    pub fn new() -> Self {
        Self {
            projects: HashMap::new(),
        }
    }

    pub fn add_project(&mut self, path: PathBuf) -> Result<(), PadoError> {
        let path = path.canonicalize().map_err(PadoError::Io)?;
        let name = path
            .file_name()
            .and_then(|n| n.to_str())
            .unwrap_or("unknown")
            .to_string();

        // FIXME: Deprecated project types use info for project
        let project_type = detect_project_type(&path).as_str().to_string();
        let key = path.to_string_lossy().to_string();

        let (access_count, starred) = if let Some(existing) = self.projects.get(&key) {
            (existing.access_count, existing.starred)
        } else {
            (1, false)
        };

        self.projects.insert(
            key,
            TrackedProject {
                path,
                name,
                project_type,
                last_accessed: Utc::now(),
                access_count,
                starred,
            },
        );

        Ok(())
    }

    pub fn remove_project(&mut self, path: &Path) -> bool {
        let key = path.to_string_lossy().to_string();
        self.projects.remove(&key).is_some()
    }

    pub fn clear(&mut self) {
        self.projects.clear();
    }

    pub fn cleanup(&mut self) {
        self.projects.retain(|_, project| project.path.exists());
    }

    pub fn update_access_time(&mut self, path: &Path) {
        let key = path.to_string_lossy().to_string();
        if let Some(project) = self.projects.get_mut(&key) {
            project.last_accessed = Utc::now();
            project.access_count += 1;
        }
    }

    pub fn toggle_star(&mut self, path: &Path) -> bool {
        let key = path.to_string_lossy().to_string();
        if let Some(project) = self.projects.get_mut(&key) {
            project.starred = !project.starred;
            project.starred
        } else {
            false
        }
    }

    pub fn set_star(&mut self, path: &Path, starred: bool) -> bool {
        let key = path.to_string_lossy().to_string();
        if let Some(project) = self.projects.get_mut(&key) {
            project.starred = starred;
            true
        } else {
            false
        }
    }

    pub fn get_projects(&self) -> Vec<&TrackedProject> {
        self.get_projects_sorted_by("time")
    }

    pub fn get_projects_sorted_by(&self, sort_by: &str) -> Vec<&TrackedProject> {
        let mut projects: Vec<_> = self.projects.values().collect();
        match sort_by {
            "access" | "frequency" => {
                projects.sort_by(|a, b| {
                    b.access_count
                        .cmp(&a.access_count)
                        .then_with(|| b.last_accessed.cmp(&a.last_accessed))
                });
            }
            "name" => {
                projects.sort_by(|a, b| a.name.cmp(&b.name));
            }
            "time" | _ => {
                projects.sort_by(|a, b| b.last_accessed.cmp(&a.last_accessed));
            }
        }
        projects
    }

    pub fn get_recent_projects(&self, limit: usize) -> Vec<&TrackedProject> {
        let mut projects: Vec<_> = self.projects.values().collect();
        projects.sort_by(|a, b| b.last_accessed.cmp(&a.last_accessed));
        projects.into_iter().take(limit).collect()
    }

    pub fn get_starred_projects(&self) -> Vec<&TrackedProject> {
        let mut projects: Vec<_> = self.projects.values().filter(|p| p.starred).collect();
        projects.sort_by(|a, b| b.last_accessed.cmp(&a.last_accessed));
        projects
    }
}

pub fn get_project_list_path() -> Result<PathBuf, PadoError> {
    let data_dir = dirs::data_dir()
        .ok_or_else(|| PadoError::InvalidPath("could not determine data directory".to_string()))?;

    let pado_dir = data_dir.join("pado");
    fs::create_dir_all(&pado_dir)?;

    Ok(pado_dir.join("projects.json"))
}

pub fn load_project_list() -> Result<ProjectList, PadoError> {
    let path = get_project_list_path()?;

    if !path.exists() {
        return Ok(ProjectList::new());
    }

    let contents = fs::read_to_string(&path)?;
    serde_json::from_str(&contents)
        .map_err(|e| PadoError::SerializationError(format!("projects.json: {}", e)))
}

pub fn save_project_list(list: &ProjectList) -> Result<(), PadoError> {
    let path = get_project_list_path()?;
    let contents = serde_json::to_string_pretty(list)
        .map_err(|e| PadoError::SerializationError(format!("projects.json: {}", e)))?;

    fs::write(&path, contents)?;
    Ok(())
}
