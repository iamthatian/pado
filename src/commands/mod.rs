use anyhow::Result;
use clap::{Parser, Subcommand};
use std::path::PathBuf;

pub trait Run {
    fn run(self) -> Result<()>;
}

#[derive(Parser)]
#[command(name = "pd")]
#[command(about = "Find project root and perform project-aware operations", long_about = None)]
pub struct Cli {
    #[command(subcommand)]
    pub command: Option<Commands>,
}

#[derive(Subcommand)]
pub enum Commands {
    Tree,
    Init,
    Files {
        #[arg(long, short)]
        pattern: Option<String>,
    },
    Find {
        pattern: String,
        #[arg(long)]
        print: bool,
    },
    Search {
        query: String,
    },
    Info,
    Type,
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
    Switch {
        #[arg(long)]
        recent: bool,
        #[arg(long)]
        starred: bool,
    },
    Recent {
        #[arg(long, short)]
        limit: Option<usize>,
    },
    Stats,
    Star {
        path: Option<PathBuf>,
        #[arg(long)]
        unstar: bool,
    },
    Add {
        path: Option<PathBuf>,
    },
    Remove {
        path: Option<PathBuf>,
        #[arg(long)]
        all: bool,
    },
    Clear,
    Cleanup,
    Discover {
        path: PathBuf,
        #[arg(long, short, default_value = "3")]
        depth: usize,
    },
    Build,
    Test,
    Run,
    Compile,
    Exec {
        command: String,
    },
    ExecAll {
        #[arg(trailing_var_arg = true, allow_hyphen_values = true)]
        command: Vec<String>,
        #[arg(long)]
        tag: Option<String>,
    },
    Prompt,
    Health,
    Deps,
    Outdated,
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
mod files;
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

                let root = pado::find_project_root(&cwd).context("no project root found")?;

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
            Commands::Info => info::run_info(),
            Commands::Type => info::run_type(),
            Commands::Health => info::run_health(),
            Commands::Deps => info::run_deps(),
            Commands::Outdated => info::run_outdated(),
            Commands::Build | Commands::Compile => build_cmd::run_build(),
            Commands::Test => build_cmd::run_test(),
            Commands::Run => build_cmd::run_run(),
            Commands::Exec { command } => build_cmd::run_exec(command),
            Commands::ExecAll { command, tag } => build_cmd::run_exec_all(command, tag),
            Commands::Tree => files::run_tree(),
            Commands::Files { pattern } => files::run_files(pattern),
            Commands::Find { pattern, print } => files::run_find(pattern, print),
            Commands::Search { query } => files::run_search(query),
            Commands::Config { path, edit, show } => config_cmd::run_config(path, edit, show),
        }
    }
}
