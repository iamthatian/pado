# pado (pd)

A tidal flow through your projects - minimalist project management CLI inspired by Emacs Projectile.

Like the tide flows naturally between shores, `pd` helps you flow naturally between your projects.

## Philosophy

- **One tool, one job**: Find project roots and manage projects
- **Unix composable**: Works with `fzf`, `rg`, `fd`, and other tools
- **Shell integration**: Rich functionality through shell functions
- **Minimal and fast**: Core binary is small and quick

## Installation

```bash
cargo install pado
```

### Shell Integration

Add to your shell config:

```bash
# bash/zsh
eval "$(pd init)"

# fish
pd init | source
```

## Usage

### Basic Commands

```bash
# Find project root (default behavior)
pd

# Interactive project switcher (requires fzf)
pd switch

# List all tracked projects
pd list

# Show recent projects
pd recent

# Add current directory as tracked project
pd add

# Star/favorite a project
pd star

# Configure pado
pd config --show
```

### Shell Functions

After running `pd init`, you get these shell functions:

```bash
# Jump to project root
pdcd

# Find files interactively (requires fzf + fd)
pdedit

# Search in project (requires rg + fzf)
pdsearch

# Interactive project jump with actions
pdjump

# Show project statistics
pdstats

# Build/test/run (auto-detects build system)
pdbuild
pdtest
pdrun

# Show dependencies
pddeps

# Check for outdated dependencies
pdoutdated
```

### Project Tracking

Projects are automatically tracked when accessed. Manage them with:

```bash
pd list                    # List all projects
pd list --starred          # Show only starred
pd list --sort-by access   # Sort by access frequency
pd recent --limit 5        # Show 5 most recent
pd switch --recent         # Switch to recent project
pd switch --starred        # Switch to starred project
pd remove                  # Remove project (interactive)
pd cleanup                 # Remove missing projects
pd discover ~/dev          # Find all projects in directory
pd clear                   # Clear all tracked projects
```

## Configuration

Configuration is optional. Create `~/.config/pado/config.toml`:

```toml
[defaults]
sort_order = "access"      # "time", "access", or "name"
output_format = "table"    # "table", "paths", or "json"
recent_limit = 10

[prompt]
format = "{name}:{type}"   # Variables: {name}, {type}, {path}
show_full_path = false

[behavior]
auto_star_on_add = false
auto_add_on_cd = true
```

### Project-Specific Config

Add `.pd.toml` in project root for custom commands:

```toml
[commands]
build = "cargo build --release"
test = "cargo nextest run"
run = "cargo watch -x run"
deploy = "./scripts/deploy.sh"
```

Then use:
```bash
pd build     # Runs custom build command
pd exec deploy   # Runs custom deploy command
```

## Shell Prompt Integration

Add project info to your prompt:

```bash
# bash
PS1='[$(pd_prompt)] \w \$ '

# zsh
PROMPT='[$(pd_prompt)] %~ %# '

# fish (add to fish_prompt function)
function fish_prompt
    set_color blue
    echo -n (pd_prompt)
    set_color normal
    echo -n ' '
    echo -n (prompt_pwd)' > '
end
```

## Multi-Project Operations

Execute commands across all projects:

```bash
# Run in all projects
pd exec-all git status

# Run in specific project types
pd exec-all --tag rust cargo check
pd exec-all --tag node npm outdated
```

## Project Statistics

```bash
pd stats    # Show tracked project statistics
```

Output includes:
- Total projects tracked
- Most accessed projects
- Project type distribution
- Recently accessed projects

## Why "pado"?

Portuguese for "tide" - representing the natural ebb and flow between projects.

## Companion Tools

- **padofetch**: Show detailed project information and statistics
  ```bash
  padofetch info    # Language stats, git info
  padofetch health  # Health check
  ```

## License

MIT
