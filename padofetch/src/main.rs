use anyhow::Result;
use clap::Parser;

mod cli;
mod commands;
mod stats;
mod utils;

use cli::{Cli, Commands};
use commands::{run_deps, run_health, run_info, run_outdated};

fn main() -> Result<()> {
    let cli = Cli::parse();

    let command = cli.command.unwrap_or(Commands::Info);

    match command {
        Commands::Info => run_info(),
        Commands::Health => run_health(),
        Commands::Deps => run_deps(),
        Commands::Outdated => run_outdated(),
    }
}
