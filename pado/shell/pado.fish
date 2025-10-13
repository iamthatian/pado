# Pado shell integration for fish
# Add this to your ~/.config/fish/config.fish:
#   pd init | source

function pdcd
    set -l root (pd)
    or return
    cd $root
end

function pdfind
    set -l root (pd)
    or return
    if command -v fzf &> /dev/null; and command -v fd &> /dev/null
        cd $root && fd --type f | fzf
    else
        echo "pdfind requires fzf and fd to be installed" >&2
        return 1
    end
end

function pdgrep
    set -l root (pd)
    or return
    if command -v rg &> /dev/null
        cd $root && rg $argv
    else
        echo "pdsearch requires ripgrep (rg) to be installed" >&2
        return 1
    end
end

function pdsearch
    set -l root (pd)
    or return

    if not type -q rg; or not type -q fzf
        echo "pdsearch requires both ripgrep (rg) and fzf" >&2
        return 1
    end

    cd $root || return

    set -l result (fzf --ansi --disabled --no-sort --delimiter : \
        --bind "change:reload:sleep 0.1; rg --line-number --color=always --no-heading --smart-case {q} || true" \
        --bind "ctrl-r:reload:sleep 0.1; rg --line-number --color=always --no-heading --smart-case {q} || true" \
        --preview 'bat --color=always --style=numbers --highlight-line {2} {1}' \
        --preview-window 'up,60%,border-bottom,+{2}+3/3' \
        --prompt '🔍 Search> ' \
        --height=90% \
        --layout=reverse | tr -d '\r')

    test -z "$result"; and return

    set -l file (echo $result | cut -d: -f1)
    set -l line (echo $result | cut -d: -f2)

    if test -n "$file" -a -n "$line"
        switch $EDITOR
            case nvim vim
                eval "$EDITOR +$line $file"
            case code
                eval "code -g $file:$line"
            case 'emacsclient*'
                eval "emacsclient -n +$line:$file"
            case '*'
                eval "$EDITOR $file"
        end
    end
end

function pdedit
    set -l root (pd)
    or return

    if command -v fzf >/dev/null; and command -v fd >/dev/null
        cd $root || return

        set -l file (fd --type f --strip-cwd-prefix | fzf --no-multi --ansi | tr -d '\r')

        if test -n "$file"
            eval "$EDITOR" "$file"
        end
    else
        echo "pdedit requires fzf and fd to be installed" >&2
        return 1
    end
end

function pdswitch
    if command -v fzf &> /dev/null
        eval (pd switch)
    else
        echo "pdswitch requires fzf to be installed" >&2
        return 1
    end
end

function pdjump
    if not type -q fzf
        echo "pdjump requires fzf to be installed" >&2
        return 1
    end

    set -l switch_out (pd switch --print-only 2>/dev/null; or pd switch 2>/dev/null)
    set switch_out (string replace -r '\r' '' $switch_out)

    if string match -q "cd*" $switch_out
        set -l root (string replace -r '^cd[[:space:]]+' '' $switch_out)
        set root (string trim --chars="'\"" $root)
    else
        set -l root $switch_out
    end

    if not test -d "$root"
        echo "Directory not found: $root" >&2
        return 1
    end

    cd "$root" || return

    set -l action (printf "Find files\nSearch text\nShow tree\n" |
        fzf --prompt 'Choose action> ' --ansi --height=30% --layout=reverse | tr -d '\r')

    test -z "$action"; and return

    switch $action
        case "Find files"
            pdedit
        case "Search text"
            pdsearch
        case "Show tree"
            pdtree
        case '*'
            echo "Unknown action: $action"
    end
end

# Project info commands (uses padofetch binary)
function pdinfo
    padofetch info
end

function pdhealth
    padofetch health
end

# Build/Test/Run commands - auto-detect build system
function pdbuild
    set -l root (pd)
    or return
    cd $root || return

    # Check for custom command in .pd.toml
    if test -f .pd.toml; and type -q toml
        set -l custom_cmd (toml get .pd.toml commands.build 2>/dev/null)
        if test -n "$custom_cmd"
            echo "Running custom build command: $custom_cmd"
            eval $custom_cmd
            return
        end
    end

    # Auto-detect build system
    if test -f Cargo.toml
        cargo build
    else if test -f package.json
        if test -f bun.lockb
            bun run build
        else if test -f pnpm-lock.yaml
            pnpm run build
        else if test -f yarn.lock
            yarn build
        else
            npm run build
        end
    else if test -f go.mod
        go build ./...
    else if test -f Makefile
        make
    else if test -f build.gradle; or test -f build.gradle.kts
        ./gradlew build
    else if test -f pom.xml
        mvn compile
    else if test -f pyproject.toml
        if type -q poetry
            poetry build
        else if type -q pip
            pip install -e .
        end
    else
        echo "No build system detected" >&2
        return 1
    end
end

function pdcompile
    pdbuild $argv
end

function pdtest
    set -l root (pd)
    or return
    cd $root || return

    # Check for custom command in .pd.toml
    if test -f .pd.toml; and type -q toml
        set -l custom_cmd (toml get .pd.toml commands.test 2>/dev/null)
        if test -n "$custom_cmd"
            echo "Running custom test command: $custom_cmd"
            eval $custom_cmd
            return
        end
    end

    # Auto-detect test framework
    if test -f Cargo.toml
        cargo test
    else if test -f package.json
        if test -f bun.lockb
            bun test
        else if test -f pnpm-lock.yaml
            pnpm test
        else if test -f yarn.lock
            yarn test
        else
            npm test
        end
    else if test -f go.mod
        go test ./...
    else if test -f Makefile
        make test
    else if test -f build.gradle; or test -f build.gradle.kts
        ./gradlew test
    else if test -f pom.xml
        mvn test
    else if test -f pyproject.toml
        if type -q poetry
            poetry run pytest
        else if type -q pytest
            pytest
        end
    else
        echo "No test framework detected" >&2
        return 1
    end
end

function pdrun
    set -l root (pd)
    or return
    cd $root || return

    # Check for custom command in .pd.toml
    if test -f .pd.toml; and type -q toml
        set -l custom_cmd (toml get .pd.toml commands.run 2>/dev/null)
        if test -n "$custom_cmd"
            echo "Running custom run command: $custom_cmd"
            eval $custom_cmd
            return
        end
    end

    # Auto-detect run command
    if test -f Cargo.toml
        cargo run
    else if test -f package.json
        if test -f bun.lockb
            bun start
        else if test -f pnpm-lock.yaml
            pnpm start
        else if test -f yarn.lock
            yarn start
        else
            npm start
        end
    else if test -f go.mod
        go run .
    else if test -f Makefile
        make run
    else
        echo "No run command detected" >&2
        return 1
    end
end

# File operations
function pdtree
    set -l root (pd)
    or return
    if type -q tree
        tree $root
    else
        echo "pdtree requires tree to be installed" >&2
        return 1
    end
end

function pdfiles
    set -l root (pd)
    or return
    cd $root || return

    if type -q fd
        fd --type f
    else if type -q find
        find . -type f -not -path '*/\.*' -not -path '*/node_modules/*' -not -path '*/target/*'
    else
        echo "pdfiles requires fd or find to be installed" >&2
        return 1
    end
end

# Dependency management
function pddeps
    set -l root (pd)
    or return
    cd $root || return

    if test -f Cargo.toml
        cargo tree
    else if test -f package.json
        if test -f bun.lockb
            bun pm ls
        else if test -f pnpm-lock.yaml
            pnpm list
        else if test -f yarn.lock
            yarn list
        else
            npm list
        end
    else if test -f go.mod
        go list -m all
    else if test -f requirements.txt; or test -f pyproject.toml
        pip list
    else
        echo "No dependency file detected" >&2
        return 1
    end
end

function pdoutdated
    set -l root (pd)
    or return
    cd $root || return

    if test -f Cargo.toml
        cargo outdated
    else if test -f package.json
        if test -f bun.lockb
            bun outdated
        else if test -f pnpm-lock.yaml
            pnpm outdated
        else if test -f yarn.lock
            yarn outdated
        else
            npm outdated
        end
    else if test -f go.mod
        go list -u -m all
    else if test -f requirements.txt; or test -f pyproject.toml
        pip list --outdated
    else
        echo "No dependency file detected" >&2
        return 1
    end
end

function pdrecent
    if command -v fzf &> /dev/null
        eval (pd switch --recent)
    else
        echo "pdrecent requires fzf to be installed" >&2
        return 1
    end
end

function pdstarred
    if command -v fzf &> /dev/null
        eval (pd switch --starred)
    else
        echo "pdstarred requires fzf to be installed" >&2
        return 1
    end
end

# Prompt integration for fish shell
# Example usage - add to your fish_prompt function:
#   function fish_prompt
#       set_color blue
#       echo -n (pd_prompt)
#       set_color normal
#       echo -n ' '
#       echo -n (prompt_pwd)' > '
#   end
function pd_prompt
    set -l info (pd prompt 2>/dev/null)
    if test -n "$info"
        echo -n "[$info]"
    end
end
