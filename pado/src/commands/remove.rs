use anyhow::{Context, Result};
use std::io::Write;
use std::path::PathBuf;
use std::process::{Command, Stdio, exit};

pub fn run_remove(path: Option<PathBuf>, all: bool) -> Result<()> {
    let mut list = pado::load_project_list().context("failed to load project list")?;

    if let Some(p) = path {
        if list.remove_project(&p) {
            pado::save_project_list(&list).context("failed to save project list")?;
            println!("Removed project: {}", p.display());
        } else {
            eprintln!("Project not found: {}", p.display());
            exit(1);
        }
    } else {
        let projects = list.get_projects();

        if projects.is_empty() {
            eprintln!("No projects tracked yet.");
            exit(1);
        }

        let mut fzf_args = vec!["--height=40%", "--reverse"];
        if all {
            fzf_args.push("--multi");
        }

        let mut fzf = Command::new("fzf")
            .args(&fzf_args)
            .stdin(Stdio::piped())
            .stdout(Stdio::piped())
            .spawn()
            .context("failed to spawn fzf - is it installed?")?;

        if let Some(stdin) = fzf.stdin.as_mut() {
            for project in projects {
                writeln!(
                    stdin,
                    "{}\t{}\t{}",
                    project.name,
                    project.project_type,
                    project.path.display()
                )?;
            }
        }

        let output = fzf.wait_with_output()?;

        if output.status.success() && !output.stdout.is_empty() {
            let selections = String::from_utf8_lossy(&output.stdout);
            let mut removed = 0;

            for line in selections.lines() {
                if let Some(path_str) = line.split('\t').nth(2) {
                    let path = PathBuf::from(path_str.trim());
                    if list.remove_project(&path) {
                        println!("Removed: {}", path.display());
                        removed += 1;
                    }
                }
            }

            if removed > 0 {
                pado::save_project_list(&list).context("failed to save project list")?;
                println!("\nRemoved {} project(s)", removed);
            }
        }
    }
    Ok(())
}

pub fn run_clear() -> Result<()> {
    let mut list = pado::load_project_list().context("failed to load project list")?;

    let count = list.projects.len();
    list.clear();

    pado::save_project_list(&list).context("failed to save project list")?;

    println!("Cleared {} projects", count);
    Ok(())
}

pub fn run_cleanup() -> Result<()> {
    let mut list = pado::load_project_list().context("failed to load project list")?;

    let before = list.projects.len();
    list.cleanup();
    let after = list.projects.len();

    pado::save_project_list(&list).context("failed to save project list")?;

    println!("Removed {} missing projects", before - after);
    Ok(())
}

pub fn run_discover(path: PathBuf, depth: usize) -> Result<()> {
    println!(
        "Discovering projects in {} (depth: {})...",
        path.display(),
        depth
    );

    let projects = libpado::discover_projects(&path, depth).context("failed to discover projects")?;

    let mut list = pado::load_project_list().context("failed to load project list")?;

    let mut added = 0;
    for project in projects {
        let key = project.to_string_lossy().to_string();
        if !list.projects.contains_key(&key) {
            list.add_project(project.clone())?;
            println!("  Found: {}", project.display());
            added += 1;
        }
    }

    pado::save_project_list(&list).context("failed to save project list")?;

    println!("\nDiscovered {} new projects", added);
    Ok(())
}
