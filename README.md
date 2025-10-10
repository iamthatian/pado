# Parkour

Parkour is a fast project navigator inspired by Emacs Projectile,
for people who work across a bunch of projects.

## Highlights

- Detects project roots for dozens of ecosystems (Rust, Node, Python, Go, Java,
  .NET, and more).
- Keeps a lightweight history of projects you’ve opened.
- Runs build, test, or run command by autodetecting the active build system or
  honoring per-project overrides (`.pk.toml`).
- Generates language and Git contribution statistics directly in the terminal.
- Offers rich project health checks, dependency summaries, and outdated
  package audits.
- Ships with optional shell helpers via `eval "$(pk init)"`.

## Installation

### Prerequisites

- Rust 1.80+ and `cargo`.
- Optional CLI tools used by specific commands:
  - [`fzf`](https://github.com/junegunn/fzf) for interactive selection.
  - [`tree`](https://treecommand.sourceforge.net/) for `pk tree`.
  - [`rg`](https://github.com/BurntSushi/ripgrep) for `pk search`.

### From source

```bash
cargo install --path .
```

## Quick start

1. Add the shell helpers to your profile (Bash example):

   ```bash
   eval "$(pk init)"
   ```

   The script shell functions to use along with `pk`.

2. Jump to the nearest project root:

   ```bash
   pk               # prints the detected project root
   eval "$(pk cd)"  # change into that root and track the visit
   ```

3. Track and inspect projects:

   ```bash
   pk add                # add the current project to the tracked list
   pk list               # show tracked projects (table view)
   pk list --starred     # focus on favorites
   pk switch --recent    # choose a recent project via fzf
   pk stats              # see aggregate usage stats
   ```

## Project tracking workflow

Key commands:

- `pk add [PATH]` – register a project (auto-stars if configured).
- `pk star [PATH]` / `pk star --unstar` – toggle favorites.
- `pk list --format paths|json|table` – list tracked repositories with
  customizable output and sorting (`--sort-by name|access|time`).
- `pk switch [--recent|--starred]` – interactively jump to a tracked project.
- `pk discover <PATH> --depth <n>` – scan directories for repositories and add
  them automatically.
- `pk recent --limit <n>` – show your latest visits
  formatting.
- `pk cleanup` – purge entries that no longer exist on disk.

## Project insight commands

Run these from within a project (or after `eval "$(pk cd)"`):

- `pk info` – overview with detected project type, top languages (via `tokei`),
  Git contributor summary, and repository health hints.
- `pk type` – print just the detected project type slug.
- `pk health` – highlight missing essentials like `.gitignore`, README, and
  license files.
- `pk deps` – quick dependency summary (Cargo / Python requirements).
- `pk outdated` – delegate to the appropriate tool (`cargo outdated`, `npm
  outdated`, `pip list --outdated`, etc.).
- `pk files [--pattern "*suffix"]` – list files within the project,
  respecting `.gitignore`.
- `pk find <pattern>` – feed matching files into `fzf`.
- `pk search <query>` – run ripgrep over the project.
- `pk tree` – display the directory tree from the root.

## Build, test, and run automation

Parkour detects build systems and applies sensible defaults.

- Rust (`cargo`)
- Node (`npm`, `yarn`, `pnpm`, `bun`)
- Python (`uv`, `poetry`, `pip`)
- Java (`maven`, `gradle`)
- Go
- Elixir (`mix`)
- Scala (`sbt`)
- Swift (`swift`)
- .NET (`dotnet`)
- Haskell (`stack`, `cabal`)
- OCaml (`dune`)
- Zig
- Terraform
- Nix
- C/C++ (`cmake`, `make`)

Use the following commands to run the detected workflows:

- `pk build`
- `pk test`
- `pk run`
- `pk exec <name>` – execute a custom command defined in `.pk.toml`.
- `pk exec-all <cmd...> [--tag <type>]` – run a command across every tracked
  project (optionally filtered by detected project type).

### Project-specific overrides

Place a `.pk.toml` file in a project root to define custom commands:

```toml
[commands]
build = "npm run build"
test = "pnpm test -- --watch=false"
run = "npm start"
lint = "pnpm lint"
```

`pk exec lint` now works alongside `pk build`, `pk test`, and `pk run`.

## Configuration

Global settings live at `~/.config/parkour/config.toml` (Linux),
`~/Library/Application Support/parkour/config.toml` (macOS), or the equivalent
`dirs::config_dir` location on your platform. Manage the file with:

```bash
pk config --path   # show the resolved path
pk config --show   # dump the current configuration
pk config --edit   # open in $EDITOR (creates the file with defaults)
```

Defaults you can safely tweak today:

```toml
[markers]
additional = [".project-marker"]   # extra files/folders that mark a project root

[defaults]
sort_order = "access"              # name | access | time
output_format = "paths"            # table | paths | json
recent_limit = 20

[prompt]
format = "[{type}] {name}"
show_full_path = true

[behavior]
auto_star_on_add = true
```

## Optional shell tooling

`pk init` initializes scripts for your current shell (Bash, Zsh, or Fish) to
provide:

- `pk_prompt` for embedding project context in your `$PS1`.
- Guarded helpers that warn if supporting tools (`fzf`, `fd`, `rg`) are missing.

Embed the script with `eval "$(pk init)"` or save it into your shell startup
file manually.

## Development

We welcome issues and pull requests especially around expanding project
detectors, improving config ergonomics, or integrating additional CLIs. Check
the `tests/` directory for examples that cover the core behavior.
