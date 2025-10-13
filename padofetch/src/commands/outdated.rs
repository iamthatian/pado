use anyhow::{Context, Result};
use std::env;
use std::process::{Command, exit};

pub fn run_outdated() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;
    let info = libpado::get_project_info(&cwd).context("no project root found")?;

    let build_system = libpado::BuildSystem::detect(&info.root);

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
