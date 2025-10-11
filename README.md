# Pado

Pado is a fast project navigator inspired by Emacs Projectile,
for people who work across a bunch of projects.

## Highlights

- Detects project roots for dozens of ecosystems (Rust, Node, Python, Go, Java,
  .NET, and more).
- Keeps a lightweight history of projects you’ve opened.
- Runs build, test, or run command by autodetecting the active build system or
  honoring per-project overrides (`.pado.toml`).
- Generates language and Git contribution statistics directly in the terminal.
- Offers rich project health checks, dependency summaries, and outdated
  package audits.
- Ships with optional shell helpers, prompt integration, and
  more) via `eval "$(pd init)"`.

## Installation

### Prerequisites

- Rust 1.80+ and `cargo`.
- Optional CLI tools used by specific commands:
  - [`fzf`](https://github.com/junegunn/fzf) for interactive selection.
  - [`tree`](https://treecommand.sourceforge.net/) for `pd tree`.
  - [`rg`](https://github.com/BurntSushi/ripgrep) for `pd search`.

### From source

```bash
cargo install --path .
```

## Quick start

1. Add the shell helpers to your profile (Bash example):

   ```bash
   eval "$(pd init)"
   ```

   The script shell functions to use along with `pd`.
   The script defines conveniences such as `pdcd`, and `pd_prompt`
   without mutating your shell if `pd` is absent.

2. Jump to the nearest project root:

   ```bash
   pd               # prints the detected project root
   cd "$(pd)"       # change into that root
   ```

3. Track and inspect projects:

   ```bash
   pd add                # add the current project to the tracked list
   pd list               # show tracked projects (table view)
   pd list --starred     # focus on favorites
   pd switch --recent    # choose a recent project via fzf
   pd stats              # see aggregate usage stats
   ```

## Project tracking workflow

Key commands:

- `pd add [PATH]` – register a project (auto-stars if configured).
- `pd star [PATH]` / `pd star --unstar` – toggle favorites.
- `pd list --format paths|json|table` – list tracked repositories with
  customizable output and sorting (`--sort-by name|access|time`).
- `pd switch [--recent|--starred]` – interactively jump to a tracked project.
- `pd discover <PATH> --depth <n>` – scan directories for repositories and add
  them automatically.
- `pd recent --limit <n>` – show your latest visits
  formatting.
- `pd cleanup` – purge entries that no longer exist on disk.

## Project insight commands

Run these from within a project (or after `eval "$(pd cd)"`):

- `pd info` – overview with detected project type, top languages (via `tokei`),
  Git contributor summary, and repository health hints.
- `pd type` – print just the detected project type slug.
- `pd health` – highlight missing essentials like `.gitignore`, README, and
  license files.
- `pd deps` – quick dependency summary (Cargo / Python requirements).
- `pd outdated` – delegate to the appropriate tool (`cargo outdated`, `npm
  outdated`, `pip list --outdated`, etc.).
- `pd files [--pattern "*suffix"]` – list files within the project,
  respecting `.gitignore`.
- `pd find <pattern>` – feed matching files into `fzf`.
- `pd search <query>` – run ripgrep over the project.
- `pd tree` – display the directory tree from the root.

## Build, test, and run automation

Pado detects build systems and applies sensible defaults.

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

- `pd build`
- `pd test`
- `pd run`
- `pd exec <name>` – execute a custom command defined in `.pado.toml`.
- `pd exec-all <cmd...> [--tag <type>]` – run a command across every tracked
  project (optionally filtered by detected project type).

### Project-specific overrides

Place a `.pado.toml` file in a project root to define custom commands:

```toml
[commands]
build = "npm run build"
test = "pnpm test -- --watch=false"
run = "npm start"
lint = "pnpm lint"
```

`pd exec lint` now works alongside `pd build`, `pd test`, and `pd run`.

## Configuration

Global settings live at `~/.config/pado/config.toml` (Linux),
`~/Library/Application Support/pado/config.toml` (macOS), or the equivalent
`dirs::config_dir` location on your platform. Manage the file with:

```bash
pd config --path   # show the resolved path
pd config --show   # dump the current configuration
pd config --edit   # open in $EDITOR (creates the file with defaults)
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

`pd init` initializes scripts for your current shell (Bash, Zsh, or Fish) to
provide:

- `pd_prompt` for embedding project context in your `$PS1`.
- Guarded helpers that warn if supporting tools (`fzf`, `fd`, `rg`) are missing.

Embed the script with `eval "$(pd init)"` or save it into your shell startup
file manually.

## Development

We welcome issues and pull requests especially around expanding project
detectors, improving config ergonomics, or integrating additional CLIs. Check
the `tests/` directory for examples that cover the core behavior.
