use anyhow::{Context, Result};
use std::env;

pub fn run_init() -> Result<()> {
    let shell_path = env::var("SHELL")
        .context("$SHELL environment variable not set")?;

    let shell_name = shell_path
        .split('/')
        .last()
        .context("invalid $SHELL path")?;

    let script = match shell_name {
        "bash" => include_str!("../../shell/parkour.bash"),
        "zsh" => include_str!("../../shell/parkour.zsh"),
        "fish" => include_str!("../../shell/parkour.fish"),
        _ => anyhow::bail!("unsupported shell: {}", shell_name),
    };

    print!("{}", script);
    Ok(())
}

pub fn run_cd() -> Result<()> {
    let cwd = env::current_dir()
        .context("failed to get current directory")?;

    let root = parkour::find_project_root(&cwd)
        .context("no project root found")?;

    let mut list = parkour::load_project_list().unwrap_or_else(|_| parkour::ProjectList::new());
    let _ = list.add_project(root.clone());
    list.update_access_time(&root);
    let _ = parkour::save_project_list(&list);

    println!("cd \"{}\"", root.display());
    Ok(())
}

pub fn run_prompt() -> Result<()> {
    let config = parkour::GlobalConfig::load().unwrap_or_default();
    let cwd = env::current_dir()
        .context("failed to get current directory")?;

    if let Ok(root) = parkour::find_project_root(&cwd) {
        let name = root.file_name()
            .and_then(|n| n.to_str())
            .unwrap_or("unknown");
        let project_type = parkour::detect_project_type(&root);
        let path = if config.prompt.show_full_path {
            root.display().to_string()
        } else {
            name.to_string()
        };

        let output = config.format_prompt(name, project_type.as_str(), &path);
        print!("{}", output);
    }
    Ok(())
}
