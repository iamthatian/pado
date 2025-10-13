use anyhow::{Context, Result};
use std::env;

pub fn run_init() -> Result<()> {
    let shell_path = env::var("SHELL").context("$SHELL environment variable not set")?;

    let shell_name = shell_path
        .split('/')
        .last()
        .context("invalid $SHELL path")?;

    let script = match shell_name {
        "bash" => include_str!("../../shell/pado.bash"),
        "zsh" => include_str!("../../shell/pado.zsh"),
        "fish" => include_str!("../../shell/pado.fish"),
        _ => anyhow::bail!("unsupported shell: {}", shell_name),
    };

    print!("{}", script);
    Ok(())
}

// FIXME: just running pd should add it to the tracked list
// pub fn run_cd() -> Result<()> {
//     let cwd = env::current_dir().context("failed to get current directory")?;
//
//     let root = libpado::find_project_root(&cwd).context("no project root found")?;
//
//     let mut list = pado::load_project_list().unwrap_or_else(|_| pado::ProjectList::new());
//     let _ = list.add_project(root.clone());
//     list.update_access_time(&root);
//     let _ = pado::save_project_list(&list);
//
//     println!("cd \"{}\"", root.display());
//     Ok(())
// }

pub fn run_prompt() -> Result<()> {
    let config = pado::GlobalConfig::load().unwrap_or_default();
    let cwd = env::current_dir().context("failed to get current directory")?;

    if let Ok(root) = libpado::find_project_root(&cwd) {
        let name = root
            .file_name()
            .and_then(|n| n.to_str())
            .unwrap_or("unknown");
        let project_type = libpado::detect_project_type(&root);
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
