use anyhow::{Context, Result};
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

pub fn get_language_stats(root: &Path) -> Result<ProjectStats> {
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

pub fn get_git_stats(root: &Path) -> Result<Option<GitStats>> {
    let repo = match gix::discover(root) {
        Ok(r) => r,
        Err(e) => {
            let err_str = e.to_string();
            if err_str.contains("could not find repository")
                || err_str.contains("not a git repository")
            {
                return Ok(None);
            }
            return Err(anyhow::anyhow!("failed to open repository: {}", e));
        }
    };

    let mut contributor_map: HashMap<String, Contributor> = HashMap::new();
    let mut total_commits = 0;
    let mut last_commit_time: Option<i64> = None;

    let head = repo.head().context("failed to get HEAD")?;

    if let Some(head_ref) = head.try_into_referent() {
        let head_id = head_ref.id();
        let platform = repo.rev_walk([head_id]);
        let iter = platform.all().context("failed to walk commits")?;

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
