use anyhow::{Context, Result};
use std::env;
use std::process::{Command, exit};

pub fn run_info() -> Result<()> {
    let cwd = env::current_dir()
        .context("failed to get current directory")?;

    let root = parkour::find_project_root(&cwd)
        .context("no project root found")?;

    let project_name = root
        .file_name()
        .and_then(|n| n.to_str())
        .unwrap_or("unknown");
    let project_type = parkour::detect_project_type(&root);

    println!("\nProject: {}", project_name);
    println!("Root: {}", root.display());
    println!("Type: {}", project_type.as_str());

    println!("\nLanguages:");
    match parkour::get_language_stats(&root) {
        Ok(stats) => {
            if stats.languages.is_empty() {
                println!("  No code files found");
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

                println!("\nTotal: {} lines ({} code, {} comments, {} blanks)",
                    stats.total_lines,
                    stats.total_code,
                    stats.total_comments,
                    stats.total_blanks
                );
            }
        }
        Err(e) => {
            eprintln!("Warning: Failed to get language statistics: {}", e);
        }
    }

    match parkour::get_git_stats(&root) {
        Ok(Some(git_stats)) => {
            println!("\nGit Information:");
            println!("  Total commits: {}", git_stats.total_commits);

            if !git_stats.contributors.is_empty() {
                println!("  Contributors: {}", git_stats.contributors.len());

                for contributor in git_stats.contributors.iter().take(5) {
                    println!("    - {} ({} commits)", contributor.name, contributor.commits);
                }

                if git_stats.contributors.len() > 5 {
                    println!("    ... and {} more", git_stats.contributors.len() - 5);
                }
            }

            if let Some(last_commit) = git_stats.last_commit_time {
                println!("  Last commit: {}", parkour::format_time_ago(last_commit));
            }
        }
        Ok(None) => {
            println!("\nGit Information: Not a git repository");
        }
        Err(e) => {
            eprintln!("\nWarning: Failed to get git statistics: {}", e);
        }
    }

    println!();
    Ok(())
}

pub fn run_type() -> Result<()> {
    let cwd = env::current_dir()
        .context("failed to get current directory")?;

    let root = parkour::find_project_root(&cwd)
        .context("no project root found")?;

    let project_type = parkour::detect_project_type(&root);
    println!("{}", project_type.as_str());
    Ok(())
}

pub fn run_health() -> Result<()> {
    let cwd = env::current_dir()
        .context("failed to get current directory")?;

    let root = parkour::find_project_root(&cwd)
        .context("no project root found")?;

    let project_type = parkour::detect_project_type(&root);

    println!("\nProject Health Check\n");
    println!("Project: {}", root.file_name().and_then(|n| n.to_str()).unwrap_or("unknown"));
    println!("Type: {}\n", project_type.as_str());

    let mut issues = Vec::new();
    let mut ok_checks = Vec::new();

    if root.join(".git").exists() {
        ok_checks.push("✓ Git repository");
    } else {
        issues.push("✗ Not a git repository");
    }

    if root.join(".gitignore").exists() {
        ok_checks.push("✓ .gitignore present");
    } else {
        issues.push("⚠ No .gitignore file");
    }

    match project_type {
        parkour::ProjectType::Rust => {
            if root.join("Cargo.lock").exists() {
                ok_checks.push("✓ Cargo.lock present");
            }
            if root.join("target").exists() {
                ok_checks.push("✓ Build directory exists");
            }
        }
        parkour::ProjectType::Node => {
            if root.join("node_modules").exists() {
                ok_checks.push("✓ node_modules present");
            } else {
                issues.push("⚠ node_modules missing (run npm install?)");
            }
            if root.join("package-lock.json").exists() || root.join("yarn.lock").exists() {
                ok_checks.push("✓ Lock file present");
            }
        }
        parkour::ProjectType::Python => {
            if root.join(".venv").exists() || root.join("venv").exists() {
                ok_checks.push("✓ Virtual environment found");
            } else {
                issues.push("⚠ No virtual environment detected");
            }
        }
        _ => {}
    }

    if root.join("README.md").exists() || root.join("README").exists() {
        ok_checks.push("✓ README present");
    } else {
        issues.push("⚠ No README file");
    }

    if root.join("LICENSE").exists() || root.join("LICENSE.md").exists() {
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
    let cwd = env::current_dir()
        .context("failed to get current directory")?;

    let root = parkour::find_project_root(&cwd)
        .context("no project root found")?;

    let project_type = parkour::detect_project_type(&root);

    println!("\nProject Dependencies\n");

    match project_type {
        parkour::ProjectType::Rust => {
            let cargo_toml = root.join("Cargo.toml");
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
            let package_json = root.join("package.json");
            if package_json.exists() {
                println!("Run 'npm list' or 'yarn list' to see dependencies");
            }
        }
        parkour::ProjectType::Python => {
            let requirements = root.join("requirements.txt");
            if requirements.exists() {
                println!("Dependencies from requirements.txt:\n");
                let contents = std::fs::read_to_string(&requirements)?;
                for line in contents.lines() {
                    if !line.trim().is_empty() && !line.trim().starts_with('#') {
                        println!("  {}", line.trim());
                    }
                }
            } else {
                println!("Run 'pip list' to see installed packages");
            }
        }
        _ => {
            println!("Dependency listing not implemented for this project type");
        }
    }
    println!();
    Ok(())
}

pub fn run_outdated() -> Result<()> {
    let cwd = env::current_dir()
        .context("failed to get current directory")?;

    let root = parkour::find_project_root(&cwd)
        .context("no project root found")?;

    let build_system = parkour::BuildSystem::detect(&root);

    println!("\nChecking for outdated dependencies...\n");

    let (cmd, args): (&str, Vec<&str>) = match build_system {
        parkour::BuildSystem::Cargo => ("cargo", vec!["outdated"]),
        parkour::BuildSystem::Npm => ("npm", vec!["outdated"]),
        parkour::BuildSystem::Yarn => ("yarn", vec!["outdated"]),
        parkour::BuildSystem::Pnpm => ("pnpm", vec!["outdated"]),
        parkour::BuildSystem::Pip => ("pip", vec!["list", "--outdated"]),
        _ => {
            println!("Outdated check not implemented for this build system");
            return Ok(());
        }
    };

    let status = Command::new(cmd)
        .args(&args)
        .current_dir(&root)
        .status()
        .context(format!("failed to execute {} - is it installed?", cmd))?;

    exit(status.code().unwrap_or(1));
}
