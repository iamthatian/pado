use anyhow::{Context, Result};
use std::env;
use std::io::Write;
use std::process::{Command, Stdio, exit};

pub fn run_tree() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;

    let root = pado::find_project_root(&cwd).context("no project root found")?;

    let status = Command::new("tree")
        .current_dir(root)
        .status()
        .context("failed to execute tree command - is it installed?")?;

    exit(status.code().unwrap_or(1));
}

pub fn run_files(pattern: Option<String>) -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;

    let root = pado::find_project_root(&cwd).context("no project root found")?;

    let files = pado::list_project_files(&root, pattern.as_deref())
        .context("failed to list project files")?;

    for file in files {
        println!("{}", file.display());
    }
    Ok(())
}

pub fn run_find(pattern: String, print: bool) -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;

    let root = pado::find_project_root(&cwd).context("no project root found")?;

    let files = pado::list_project_files(&root, Some(&pattern)).context("failed to find files")?;

    if print {
        for file in files {
            println!("{}", file.display());
        }
    } else {
        let mut fzf = Command::new("fzf")
            .stdin(Stdio::piped())
            .spawn()
            .context("failed to spawn fzf - is it installed?")?;

        if let Some(stdin) = fzf.stdin.as_mut() {
            for file in files {
                writeln!(stdin, "{}", file.display())?;
            }
        }

        let status = fzf.wait()?;
        exit(status.code().unwrap_or(1));
    }
    Ok(())
}

pub fn run_search(query: String) -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;

    let root = pado::find_project_root(&cwd).context("no project root found")?;

    let status = Command::new("rg")
        .arg(&query)
        .current_dir(root)
        .status()
        .context("failed to execute rg - is ripgrep installed?")?;

    exit(status.code().unwrap_or(1));
}
