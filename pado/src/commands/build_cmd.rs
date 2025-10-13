use anyhow::{Context, Result};
use std::env;
use std::process::{Command, exit};

pub fn run_exec(command: String) -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;

    let root = libpado::find_project_root(&cwd).context("no project root found")?;

    let config = pado::ProjectConfig::load(&root)
        .context("failed to load .pado.toml")?
        .context(".pado.toml not found in project root")?;

    if let Some(cmd) = config.get_command(&command) {
        println!("Running: {}", cmd);
        let status = Command::new("sh")
            .arg("-c")
            .arg(cmd)
            .current_dir(&root)
            .status()
            .context("failed to execute command")?;
        exit(status.code().unwrap_or(1));
    } else {
        eprintln!("Command '{}' not found in .pado.toml", command);
        exit(1);
    }
}

pub fn run_exec_all(command: Vec<String>, tag: Option<String>) -> Result<()> {
    if command.is_empty() {
        eprintln!("No command specified");
        exit(1);
    }

    let list = pado::load_project_list().context("failed to load project list")?;

    let projects = list.get_projects();
    let cmd = command.join(" ");

    println!("Executing '{}' in {} projects...\n", cmd, projects.len());

    for project in projects {
        if let Some(ref filter_tag) = tag {
            if &project.project_type != filter_tag {
                continue;
            }
        }

        println!("==> {} ({})", project.name, project.path.display());

        let status = Command::new("sh")
            .arg("-c")
            .arg(&cmd)
            .current_dir(&project.path)
            .status();

        match status {
            Ok(s) if s.success() => println!("    ✓ Success\n"),
            Ok(s) => println!("    ✗ Failed with code {:?}\n", s.code()),
            Err(e) => println!("    ✗ Error: {}\n", e),
        }
    }
    Ok(())
}
