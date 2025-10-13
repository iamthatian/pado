#!/usr/bin/env bash
# Install pd man pages

set -e

MAN_DIR="${MAN_DIR:-/usr/local/share/man}"

echo "Installing pd man pages to $MAN_DIR..."

# Create directories if they don't exist
sudo mkdir -p "$MAN_DIR/man1"
sudo mkdir -p "$MAN_DIR/man5"

# Install man pages
sudo cp man/pd.1 "$MAN_DIR/man1/"
sudo cp man/pd-config.5 "$MAN_DIR/man5/"
sudo cp man/pd-project.5 "$MAN_DIR/man5/"

echo "Man pages installed successfully!"

# Update man database
if command -v mandb &> /dev/null; then
    echo "Updating man database (Linux)..."
    sudo mandb
elif command -v makewhatis &> /dev/null; then
    echo "Updating man database (macOS/BSD)..."
    sudo makewhatis "$MAN_DIR"
fi

echo ""
echo "You can now view the man pages with:"
echo "  man pd"
echo "  man pd-config"
echo "  man pd-project"
