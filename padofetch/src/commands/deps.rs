use anyhow::{Context, Result};
use std::env;

pub fn run_deps() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;
    let info = libpado::get_project_info(&cwd).context("no project root found")?;

    println!("\nProject Dependencies\n");

    for project_type in &info.project_types {
        match project_type {
            libpado::ProjectType::Rust => {
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
            libpado::ProjectType::Node => {
                let package_json = info.root.join("package.json");
                if package_json.exists() {
                    println!("Run `npm list` or `yarn list` to view dependencies");
                }
            }
            libpado::ProjectType::Python => {
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
