# Man Pages

This directory contains man pages for `pd`.

## Files

- `pd.1` - Main command manual (section 1: User Commands)
- `pd-config.5` - Global configuration file format (section 5: File Formats)
- `pd-project.5` - Project-specific configuration file format (section 5: File Formats)

## Installation

### Manual Installation

```bash
# Copy man pages to standard locations
sudo cp pd.1 /usr/local/share/man/man1/
sudo cp pd-config.5 /usr/local/share/man/man5/
sudo cp pd-project.5 /usr/local/share/man/man5/

# Update man database
sudo mandb  # Linux
sudo makewhatis /usr/local/share/man  # macOS
```

### Installation Script

```bash
# From repository root
./install-man.sh
```

### Viewing

```bash
man pd
man pd-config
man pd-project
```

### Previewing Without Installing

```bash
# From man directory
man ./pd.1
man ./pd-config.5
man ./pd-project.5
```

## Format

Man pages are written in [troff](https://www.gnu.org/software/groff/manual/) format using standard man macros.

## Sections

- **Section 1**: User commands (executable programs or shell commands)
- **Section 5**: File formats and conventions
