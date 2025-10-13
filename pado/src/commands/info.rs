use anyhow::{Context, Result};
use std::env;

pub fn run_type() -> Result<()> {
    let cwd = env::current_dir().context("failed to get current directory")?;
    let info = libpado::get_project_info(&cwd).context("no project root found")?;

    if info.project_types.is_empty() {
        println!("unknown");
    } else {
        println!(
            "{}",
            info.project_types
                .iter()
                .map(|t| t.as_str())
                .collect::<Vec<_>>()
                .join(", ")
        );
    }

    Ok(())
}
