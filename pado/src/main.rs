mod commands;

use std::env;
use std::io::{self, Write};
use std::process::ExitCode;

use clap::Parser;

use commands::{Cli, Run};

pub fn main() -> ExitCode {
    unsafe { env::remove_var("RUST_LIB_BACKTRACE") };
    unsafe { env::remove_var("RUST_BACKTRACE") };

    match Cli::parse().run() {
        Ok(()) => ExitCode::SUCCESS,
        Err(e) => {
            _ = writeln!(io::stderr(), "pd: {e:?}");
            ExitCode::FAILURE
        }
    }
}
