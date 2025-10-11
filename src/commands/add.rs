use anyhow::{Context, Result};
use std::env;
use std::path::PathBuf;
use std::process::exit;

pub fn run_add(path: Option<PathBuf>) -> Result<()> {
    let config = pado::GlobalConfig::load().unwrap_or_default();
    let mut list = pado::load_project_list().context("failed to load project list")?;

    let project_path = if let Some(p) = path {
        p
    } else {
        let cwd = env::current_dir().context("failed to get current directory")?;
        pado::find_project_root(&cwd).context("no project root found")?
    };

    list.add_project(project_path.clone())
        .context("failed to add project")?;

    if config.behavior.auto_star_on_add {
        list.set_star(&project_path, true);
    }

    pado::save_project_list(&list).context("failed to save project list")?;

    let star_msg = if config.behavior.auto_star_on_add {
        " ★"
    } else {
        ""
    };
    println!("Added project:{} {}", star_msg, project_path.display());
    Ok(())
}

pub fn run_star(path: Option<PathBuf>, unstar: bool) -> Result<()> {
    let mut list = pado::load_project_list().context("failed to load project list")?;

    let project_path = if let Some(p) = path {
        p
    } else {
        let cwd = env::current_dir().context("failed to get current directory")?;
        pado::find_project_root(&cwd).context("no project root found")?
    };

    if unstar {
        if list.set_star(&project_path, false) {
            pado::save_project_list(&list).context("failed to save project list")?;
            println!("Unstarred project: {}", project_path.display());
        } else {
            eprintln!("Project not found in list: {}", project_path.display());
            exit(1);
        }
    } else {
        let _ = list.add_project(project_path.clone());

        if list.set_star(&project_path, true) {
            pado::save_project_list(&list).context("failed to save project list")?;
            println!("★ Starred project: {}", project_path.display());
        } else {
            eprintln!("Failed to star project: {}", project_path.display());
            exit(1);
        }
    }
    Ok(())
}
