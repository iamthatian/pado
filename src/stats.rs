use crate::error::ParkourError;
use std::collections::HashMap;
use std::path::Path;

#[derive(Debug, Clone)]
pub struct LanguageStats {
    pub name: String,
    pub lines: usize,
    pub code: usize,
    pub comments: usize,
    pub blanks: usize,
}

#[derive(Debug)]
pub struct ProjectStats {
    pub languages: Vec<LanguageStats>,
    pub total_lines: usize,
    pub total_code: usize,
    pub total_comments: usize,
    pub total_blanks: usize,
}

pub fn get_language_stats(root: &Path) -> Result<ProjectStats, ParkourError> {
    use tokei::{Config, Languages};

    let config = Config::default();
    let mut languages = Languages::new();
    languages.get_statistics(&[root], &[], &config);

    let mut lang_stats = Vec::new();
    let mut total_lines = 0;
    let mut total_code = 0;
    let mut total_comments = 0;
    let mut total_blanks = 0;

    for (lang_type, language) in languages.iter() {
        let lines = language.lines();
        let code = language.code;
        let comments = language.comments;
        let blanks = language.blanks;

        if lines > 0 {
            lang_stats.push(LanguageStats {
                name: lang_type.to_string(),
                lines,
                code,
                comments,
                blanks,
            });

            total_lines += lines;
            total_code += code;
            total_comments += comments;
            total_blanks += blanks;
        }
    }

    lang_stats.sort_by(|a, b| b.lines.cmp(&a.lines));

    Ok(ProjectStats {
        languages: lang_stats,
        total_lines,
        total_code,
        total_comments,
        total_blanks,
    })
}

#[derive(Debug, Clone)]
pub struct Contributor {
    pub name: String,
    pub email: String,
    pub commits: usize,
}

#[derive(Debug)]
pub struct GitStats {
    pub contributors: Vec<Contributor>,
    pub total_commits: usize,
    pub last_commit_time: Option<i64>,
}

pub fn get_git_stats(root: &Path) -> Result<Option<GitStats>, ParkourError> {
    let repo = match gix::discover(root) {
        Ok(r) => r,
        Err(e) => {
            let err_str = e.to_string();
            if err_str.contains("could not find repository")
                || err_str.contains("not a git repository")
            {
                return Ok(None);
            }
            return Err(ParkourError::GitError(format!(
                "failed to open repository: {}",
                e
            )));
        }
    };

    let mut contributor_map: HashMap<String, Contributor> = HashMap::new();
    let mut total_commits = 0;
    let mut last_commit_time: Option<i64> = None;

    let head = repo
        .head()
        .map_err(|e| ParkourError::GitError(format!("failed to get HEAD: {}", e)))?;

    if let Some(head_ref) = head.try_into_referent() {
        let head_id = head_ref.id();

        let platform = repo.rev_walk([head_id]);

        let iter = platform
            .all()
            .map_err(|e| ParkourError::GitError(format!("failed to walk commits: {}", e)))?;

        for commit_result in iter {
            if let Ok(commit_info) = commit_result {
                if let Ok(commit) = commit_info.object() {
                    total_commits += 1;

                    if let Ok(author) = commit.author() {
                        let name = author.name.to_string();
                        let email = author.email.to_string();
                        let key = format!("{}:{}", name, email);

                        if last_commit_time.is_none() {
                            last_commit_time = Some(author.time.seconds);
                        }

                        contributor_map
                            .entry(key)
                            .and_modify(|c| c.commits += 1)
                            .or_insert(Contributor {
                                name: name.clone(),
                                email: email.clone(),
                                commits: 1,
                            });
                    }
                }
            }
        }
    }

    let mut contributors: Vec<_> = contributor_map.into_values().collect();
    contributors.sort_by(|a, b| b.commits.cmp(&a.commits));

    Ok(Some(GitStats {
        contributors,
        total_commits,
        last_commit_time,
    }))
}

pub fn format_time_ago(seconds: i64) -> String {
    use chrono::Utc;

    let now = Utc::now().timestamp();
    let diff = now - seconds;

    if diff < 60 {
        format!("{} seconds ago", diff)
    } else if diff < 3600 {
        format!("{} minutes ago", diff / 60)
    } else if diff < 86400 {
        format!("{} hours ago", diff / 3600)
    } else if diff < 2592000 {
        format!("{} days ago", diff / 86400)
    } else if diff < 31536000 {
        format!("{} months ago", diff / 2592000)
    } else {
        format!("{} years ago", diff / 31536000)
    }
}

pub fn create_percentage_bar(percentage: f64, width: usize) -> String {
    let filled = ((percentage / 100.0) * width as f64).round() as usize;
    let empty = width.saturating_sub(filled);
    format!("{}{}", "█".repeat(filled), "░".repeat(empty))
}
