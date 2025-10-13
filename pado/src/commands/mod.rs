use anyhow::Result;
use clap::{Parser, Subcommand};
use std::path::PathBuf;

pub trait Run {
    fn run(self) -> Result<()>;
}

#[derive(Parser)]
#[command(name = "pd")]
#[command(about = "Find project root and perform project-aware operations")]
#[command(after_help = "Use 'pd <command> --help' for more information on a specific command.")]
#[command(help_template = "\
{about-with-newline}
{usage-heading} {usage}

{before-help}Shell Integration:
  init      Output shell integration script
  prompt    Get prompt string for shell integration

Project Information:
  type      Detect and print project type

Project Management:
  list      List all known projects
  switch    Interactive project switcher
  recent    Show recent projects
  stats     Show project statistics
  star      Star/unstar a project
  add       Add project to known projects list
  remove    Remove project(s) from list
  clear     Clear all known projects
  cleanup   Remove missing/non-existent projects
  discover  Discover projects recursively

Custom Commands:
  exec      Execute a custom command from .pd.toml
  exec-all  Execute command in all tracked projects

Configuration:
  config    Show or edit configuration

{options}{after-help}
")]
pub struct Cli {
    #[command(subcommand)]
    pub command: Option<Commands>,
}

#[derive(Subcommand)]
pub enum Commands {
    #[command(display_order = 1)]
    Init,

    #[command(display_order = 2)]
    Prompt,

    #[command(display_order = 10)]
    Type,

    #[command(display_order = 30)]
    #[command(next_help_heading = "Project Management")]
    List {
        #[arg(long)]
        json: bool,
        #[arg(long, short)]
        verbose: bool,
        #[arg(long)]
        path: bool,
        #[arg(long)]
        sort_by: Option<String>,
        #[arg(long)]
        starred: bool,
    },

    #[command(display_order = 31)]
    Switch {
        #[arg(long)]
        recent: bool,
        #[arg(long)]
        starred: bool,
    },

    #[command(display_order = 32)]
    Recent {
        #[arg(long, short)]
        limit: Option<usize>,
    },

    #[command(display_order = 33)]
    Stats,

    #[command(display_order = 34)]
    Star {
        path: Option<PathBuf>,
        #[arg(long)]
        unstar: bool,
    },

    #[command(display_order = 35)]
    Add { path: Option<PathBuf> },

    #[command(display_order = 36)]
    Remove {
        path: Option<PathBuf>,
        #[arg(long)]
        all: bool,
    },

    #[command(display_order = 37)]
    Clear,

    #[command(display_order = 38)]
    Cleanup,

    #[command(display_order = 39)]
    Discover {
        path: PathBuf,
        #[arg(long, short, default_value = "3")]
        depth: usize,
    },

    #[command(display_order = 40)]
    Exec { command: String },

    #[command(display_order = 41)]
    ExecAll {
        #[arg(trailing_var_arg = true, allow_hyphen_values = true)]
        command: Vec<String>,
        #[arg(long)]
        tag: Option<String>,
    },

    #[command(display_order = 50)]
    Config {
        #[arg(long)]
        path: bool,
        #[arg(long)]
        edit: bool,
        #[arg(long)]
        show: bool,
    },
}

mod add;
mod build_cmd;
mod config_cmd;
mod info;
mod init;
mod list;
mod remove;

impl Run for Cli {
    fn run(self) -> Result<()> {
        match self.command {
            Some(cmd) => cmd.run(),
            None => {
                use anyhow::Context;
                use std::env;

                let cwd = env::current_dir().context("failed to get current directory")?;

                let root = libpado::find_project_root(&cwd).context("no project root found")?;

                if let Ok(mut list) = pado::load_project_list() {
                    let key = root.to_string_lossy().to_string();
                    let exists = list.projects.contains_key(&key);

                    if exists {
                        list.update_access_time(&root);
                    } else {
                        let _ = list.add_project(root.clone());
                    }

                    let _ = pado::save_project_list(&list);
                }

                println!("{}", root.display());
                Ok(())
            }
        }
    }
}

impl Run for Commands {
    fn run(self) -> Result<()> {
        match self {
            Commands::Init => init::run_init(),
            Commands::Prompt => init::run_prompt(),
            Commands::Add { path } => add::run_add(path),
            Commands::Star { path, unstar } => add::run_star(path, unstar),
            Commands::Remove { path, all } => remove::run_remove(path, all),
            Commands::Clear => remove::run_clear(),
            Commands::Cleanup => remove::run_cleanup(),
            Commands::Discover { path, depth } => remove::run_discover(path, depth),
            Commands::List {
                json,
                verbose,
                path,
                sort_by,
                starred,
            } => list::run_list(json, verbose, path, sort_by, starred),
            Commands::Recent { limit } => list::run_recent(limit),
            Commands::Stats => list::run_stats(),
            Commands::Switch { recent, starred } => list::run_switch(recent, starred),
            Commands::Type => info::run_type(),
            Commands::Exec { command } => build_cmd::run_exec(command),
            Commands::ExecAll { command, tag } => build_cmd::run_exec_all(command, tag),
            Commands::Config { path, edit, show } => config_cmd::run_config(path, edit, show),
        }
    }
}
