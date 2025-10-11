use std::path::PathBuf;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum PadoError {
    #[error("no project root found from {0}")]
    NoProjectRoot(PathBuf),

    #[error("project not found: {0}")]
    ProjectNotFound(PathBuf),

    #[error("failed to parse configuration: {0}")]
    ConfigParse(String),

    #[error("failed to serialize data: {0}")]
    SerializationError(String),

    #[error("invalid path: {0}")]
    InvalidPath(String),

    #[error("git operation failed: {0}")]
    GitError(String),

    #[error("unsupported operation: {0}")]
    UnsupportedOperation(String),

    #[error("I/O error: {0}")]
    Io(#[from] std::io::Error),
}
