use anyhow::{Context, Result};
use std::env;
use std::process::{Command, exit};

pub fn run_build() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;

    let root = pado::find_project_root(&cwd).context("no project root found")?;

    if let Ok(Some(config)) = pado::ProjectConfig::load(&root) {
        if let Some(cmd) = config.get_command("build") {
            println!("Running custom build command: {}", cmd);
            let status = Command::new("sh")
                .arg("-c")
                .arg(cmd)
                .current_dir(&root)
                .status()
                .context("failed to execute build command")?;
            exit(status.code().unwrap_or(1));
        }
    }

    let build_system = pado::BuildSystem::detect(&root);

    if let Some(cmd) = build_system.build_command() {
        println!("Running: {}", cmd);
        let status = Command::new("sh")
            .arg("-c")
            .arg(cmd)
            .current_dir(&root)
            .status()
            .context("failed to execute build command")?;
        exit(status.code().unwrap_or(1));
    } else {
        eprintln!("No build system detected for this project");
        exit(1);
    }
}

pub fn run_test() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;

    let root = pado::find_project_root(&cwd).context("no project root found")?;

    if let Ok(Some(config)) = pado::ProjectConfig::load(&root) {
        if let Some(cmd) = config.get_command("test") {
            println!("Running custom test command: {}", cmd);
            let status = Command::new("sh")
                .arg("-c")
                .arg(cmd)
                .current_dir(&root)
                .status()
                .context("failed to execute test command")?;
            exit(status.code().unwrap_or(1));
        }
    }

    let build_system = pado::BuildSystem::detect(&root);

    if let Some(cmd) = build_system.test_command() {
        println!("Running: {}", cmd);
        let status = Command::new("sh")
            .arg("-c")
            .arg(cmd)
            .current_dir(&root)
            .status()
            .context("failed to execute test command")?;
        exit(status.code().unwrap_or(1));
    } else {
        eprintln!("No test framework detected for this project");
        exit(1);
    }
}

pub fn run_run() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;

    let root = pado::find_project_root(&cwd).context("no project root found")?;

    if let Ok(Some(config)) = pado::ProjectConfig::load(&root) {
        if let Some(cmd) = config.get_command("run") {
            println!("Running custom run command: {}", cmd);
            let status = Command::new("sh")
                .arg("-c")
                .arg(cmd)
                .current_dir(&root)
                .status()
                .context("failed to execute run command")?;
            exit(status.code().unwrap_or(1));
        }
    }

    let build_system = pado::BuildSystem::detect(&root);

    if let Some(cmd) = build_system.run_command() {
        println!("Running: {}", cmd);
        let status = Command::new("sh")
            .arg("-c")
            .arg(cmd)
            .current_dir(&root)
            .status()
            .context("failed to execute run command")?;
        exit(status.code().unwrap_or(1));
    } else {
        eprintln!("No run command detected for this project");
        exit(1);
    }
}

pub fn run_exec(command: String) -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;

    let root = pado::find_project_root(&cwd).context("no project root found")?;

    let config = pado::ProjectConfig::load(&root)
        .context("failed to load .pd.toml")?
        .context(".pd.toml not found in project root")?;

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
        eprintln!("Command '{}' not found in .pd.toml", command);
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
