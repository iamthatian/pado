use anyhow::{Context, Result};
use std::env;
use std::process::{Command, exit};

pub fn run_config(path: bool, edit: bool, show: bool) -> Result<()> {
    let config_path = parkour::get_global_config_path().context("failed to get config path")?;

    if path {
        println!("{}", config_path.display());
    } else if edit {
        let editor = env::var("EDITOR").unwrap_or_else(|_| "vi".to_string());

        if !config_path.exists() {
            if let Some(parent) = config_path.parent() {
                std::fs::create_dir_all(parent)?;
            }
            let default_config = parkour::GlobalConfig::default();
            default_config.save()?;
        }

        let status = Command::new(&editor)
            .arg(&config_path)
            .status()
            .context(format!(
                "failed to execute {} - is $EDITOR set correctly?",
                editor
            ))?;

        exit(status.code().unwrap_or(1));
    } else if show {
        let config = parkour::GlobalConfig::load().unwrap_or_default();
        let config_toml = toml::to_string_pretty(&config).context("failed to serialize config")?;
        println!("{}", config_toml);
    } else {
        if config_path.exists() {
            println!("Config: {}", config_path.display());
            println!("\nUse 'pk config --show' to view");
            println!("Use 'pk config --edit' to edit");
        } else {
            println!("No config file found at: {}", config_path.display());
            println!("\nUse 'pk config --edit' to create one");
        }
    }
    Ok(())
}
