use anyhow::{Context, Result};
use std::env;
use std::process::{Command, exit};

// TODO: Move tailored functionality to projects.rs
pub fn run_info() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;
    let info = parkour::get_project_info(&cwd).context("no project root found")?;

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
    match parkour::get_language_stats(&info.root) {
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
                    let bar = parkour::create_percentage_bar(percentage, 20);
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

    println!("Git Information:");
    match parkour::get_git_stats(&info.root) {
        Ok(Some(git_stats)) => {
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
                    println!(
                        "    ... and {} more",
                        git_stats.contributors.len() - 5
                    );
                }
            }

            if let Some(last_commit) = git_stats.last_commit_time {
                println!("  Last commit: {}", parkour::format_time_ago(last_commit));
            }
        }
        Ok(None) => println!("  Not a git repository"),
        Err(e) => eprintln!("   Failed to get git stats: {}", e),
    }

    println!();
    Ok(())
}

pub fn run_type() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;
    let info = parkour::get_project_info(&cwd).context("no project root found")?;

    if info.project_types.is_empty() {
        println!("unknown");
    } else {
        println!(
            "{}",
            info.project_types
                .iter()
                .map(|t| t.as_str())
                .collect::<Vec<_>>()
                .join(", ")
        );
    }

    Ok(())
}

pub fn run_health() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;
    let info = parkour::get_project_info(&cwd).context("no project root found")?;

    println!("\n🩺 Project Health Check\n");
    println!(
        "Project: {}",
        info.root.file_name().and_then(|n| n.to_str()).unwrap_or("unknown")
    );
    println!("Type(s): {}\n",
        info.project_types.iter().map(|t| t.as_str()).collect::<Vec<_>>().join(", ")
    );

    let mut issues = Vec::new();
    let mut ok_checks = Vec::new();

    if info.root.join(".git").exists() {
        ok_checks.push("✓ Git repository");
    } else {
        issues.push("✗ Not a git repository");
    }

    if info.root.join(".gitignore").exists() {
        ok_checks.push("✓ .gitignore present");
    } else {
        issues.push("⚠ No .gitignore file");
    }

    for project_type in &info.project_types {
        match project_type {
            parkour::ProjectType::Rust => {
                if info.root.join("Cargo.lock").exists() {
                    ok_checks.push("✓ Cargo.lock present");
                }
                if info.root.join("target").exists() {
                    ok_checks.push("✓ target directory exists");
                }
            }
            parkour::ProjectType::Node => {
                if info.root.join("node_modules").exists() {
                    ok_checks.push("✓ node_modules present");
                } else {
                    issues.push("⚠ node_modules missing (run npm install?)");
                }
                if info.root.join("package-lock.json").exists()
                    || info.root.join("yarn.lock").exists()
                {
                    ok_checks.push("✓ Lock file present");
                }
            }
            parkour::ProjectType::Python => {
                if info.root.join(".venv").exists() || info.root.join("venv").exists() {
                    ok_checks.push("✓ Virtual environment found");
                } else {
                    issues.push("⚠ No virtual environment detected");
                }
            }
            _ => {}
        }
    }

    if info.root.join("README.md").exists() || info.root.join("README").exists() {
        ok_checks.push("✓ README present");
    } else {
        issues.push("⚠ No README file");
    }

    if info.root.join("LICENSE").exists() || info.root.join("LICENSE.md").exists() {
        ok_checks.push("✓ LICENSE present");
    } else {
        issues.push("⚠ No LICENSE file");
    }

    if !ok_checks.is_empty() {
        println!("Healthy:");
        for check in ok_checks {
            println!("  {}", check);
        }
        println!();
    }

    if !issues.is_empty() {
        println!("Issues:");
        for issue in &issues {
            println!("  {}", issue);
        }
        println!();
    }

    if issues.is_empty() {
        println!("✓ Project looks healthy!\n");
    }

    Ok(())
}

pub fn run_deps() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;
    let info = parkour::get_project_info(&cwd).context("no project root found")?;

    println!("\nProject Dependencies\n");

    for project_type in &info.project_types {
        match project_type {
            parkour::ProjectType::Rust => {
                let cargo_toml = info.root.join("Cargo.toml");
                if cargo_toml.exists() {
                    println!("Dependencies from Cargo.toml:\n");
                    let contents = std::fs::read_to_string(&cargo_toml)?;
                    let mut in_deps = false;
                    for line in contents.lines() {
                        if line.trim() == "[dependencies]" {
                            in_deps = true;
                            continue;
                        }
                        if line.trim().starts_with('[') && in_deps {
                            break;
                        }
                        if in_deps && !line.trim().is_empty() && !line.trim().starts_with('#') {
                            println!("  {}", line.trim());
                        }
                    }
                }
            }
            parkour::ProjectType::Node => {
                let package_json = info.root.join("package.json");
                if package_json.exists() {
                    println!("Run `npm list` or `yarn list` to view dependencies");
                }
            }
            parkour::ProjectType::Python => {
                let requirements = info.root.join("requirements.txt");
                if requirements.exists() {
                    println!("Dependencies from requirements.txt:\n");
                    let contents = std::fs::read_to_string(&requirements)?;
                    for line in contents.lines() {
                        if !line.trim().is_empty() && !line.trim().starts_with('#') {
                            println!("  {}", line.trim());
                        }
                    }
                } else {
                    println!("Run `pip list` to see installed packages");
                }
            }
            _ => {}
        }
    }

    println!();
    Ok(())
}

pub fn run_outdated() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;
    let info = parkour::get_project_info(&cwd).context("no project root found")?;

    let build_system = parkour::BuildSystem::detect(&info.root);

    println!("\nChecking for outdated dependencies...\n");

    if let Some((cmd, args)) = build_system.outdated_command() {
        let status = Command::new(cmd)
            .args(args)
            .current_dir(&info.root)
            .status()
            .context(format!("failed to execute {} - is it installed?", cmd))?;

        exit(status.code().unwrap_or(1));
    } else {
        println!("Outdated check not implemented for this build system");
        Ok(())
    }
}
