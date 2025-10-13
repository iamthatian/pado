use anyhow::{Context, Result};
use std::env;

pub fn run_health() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;
    let info = libpado::get_project_info(&cwd).context("no project root found")?;

    println!("\n🩺 Project Health Check\n");
    println!(
        "Project: {}",
        info.root
            .file_name()
            .and_then(|n| n.to_str())
            .unwrap_or("unknown")
    );
    println!(
        "Type(s): {}\n",
        info.project_types
            .iter()
            .map(|t| t.as_str())
            .collect::<Vec<_>>()
            .join(", ")
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
            libpado::ProjectType::Rust => {
                if info.root.join("Cargo.lock").exists() {
                    ok_checks.push("✓ Cargo.lock present");
                }
                if info.root.join("target").exists() {
                    ok_checks.push("✓ target directory exists");
                }
            }
            libpado::ProjectType::Node => {
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
            libpado::ProjectType::Python => {
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
