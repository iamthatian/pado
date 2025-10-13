use anyhow::{Context, Result};
use std::env;

use crate::stats::{get_git_stats, get_language_stats};
use crate::utils::{create_percentage_bar, format_time_ago};

pub fn run_info() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;
    let info = libpado::get_project_info(&cwd).context("no project root found")?;

    let project_name = info
        .root
        .file_name()
        .and_then(|n| n.to_str())
        .unwrap_or("unknown");

    println!("\nProject Information\n");
    println!("Name: {}", project_name);
    println!("Root: {}", info.root.display());

    if info.project_types.is_empty() {
        println!("Type: unknown");
    } else {
        let type_str = info
            .project_types
            .iter()
            .map(|t| t.as_str())
            .collect::<Vec<_>>()
            .join(", ");
        println!("Type(s): {}", type_str);
    }

    if info.monorepo {
        println!("Monorepo: Yes");
        if !info.subprojects.is_empty() {
            println!("Subprojects:");
            for sub in &info.subprojects {
                println!("  - {}", sub.display());
            }
        } else {
            println!("(no explicit subprojects found)");
        }
    } else {
        println!("Monorepo: No");
    }

    println!("Total files: {}", info.file_count);
    println!();

    println!("Language Statistics:");
    match get_language_stats(&info.root) {
        Ok(stats) => {
            if stats.languages.is_empty() {
                println!("  No source code detected");
            } else {
                for lang in stats.languages.iter().take(10) {
                    let percentage = if stats.total_lines > 0 {
                        (lang.lines as f64 / stats.total_lines as f64) * 100.0
                    } else {
                        0.0
                    };
                    let bar = create_percentage_bar(percentage, 20);
                    println!(
                        "  {:<12} {} {:>5.1}%  ({} lines)",
                        lang.name, bar, percentage, lang.lines
                    );
                }

                println!(
                    "\nTotal lines: {} ({} code, {} comments, {} blanks)\n",
                    stats.total_lines, stats.total_code, stats.total_comments, stats.total_blanks
                );
            }
        }
        Err(e) => eprintln!("   Failed to analyze languages: {}", e),
    }

    match get_git_stats(&info.root) {
        Ok(Some(git_stats)) => {
            println!("Git Information:");
            println!("  Total commits: {}", git_stats.total_commits);

            if !git_stats.contributors.is_empty() {
                println!("  Contributors: {}", git_stats.contributors.len());
                for contributor in git_stats.contributors.iter().take(5) {
                    println!(
                        "    - {} ({} commits)",
                        contributor.name, contributor.commits
                    );
                }
                if git_stats.contributors.len() > 5 {
                    println!("    ... and {} more", git_stats.contributors.len() - 5);
                }
            }

            if let Some(last_commit) = git_stats.last_commit_time {
                println!("  Last commit: {}", format_time_ago(last_commit));
            }
        }
        Ok(None) | Err(_) => {}
    }

    println!();
    Ok(())
}
