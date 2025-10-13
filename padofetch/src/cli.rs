use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command(name = "padofetch")]
#[command(about = "Show detailed information about projects")]
#[command(version)]
pub struct Cli {
    #[command(subcommand)]
    pub command: Option<Commands>,
}

#[derive(Subcommand)]
pub enum Commands {
    Info,
    Health,
    Deps,
    Outdated,
}
