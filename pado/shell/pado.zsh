# Pado shell integration for zsh
# Add this to your ~/.zshrc:
#   eval "$(pd init)"

pdcd() {
    local root
    root=$(pd) || return
    cd "$root" || return
}

pdfind() {
    local root
    root=$(pd) || return
    if command -v fzf &> /dev/null && command -v fd &> /dev/null; then
        cd "$root" && fd --type f | fzf
    else
        echo "pdfind requires fzf and fd to be installed" >&2
        return 1
    fi
}

pdgrep() {
    local root
    root=$(pd) || return
    if command -v rg &> /dev/null; then
        cd "$root" && rg "$@"
    else
        echo "pdsearch requires ripgrep (rg) to be installed" >&2
        return 1
    fi
}

pdsearch() {
    local root result file line
    root=$(pd) || return

    if ! command -v rg &>/dev/null || ! command -v fzf &>/dev/null; then
        echo "pdsearch requires both ripgrep (rg) and fzf" >&2
        return 1
    fi

    cd "$root" || return

    result=$(
        fzf --ansi --disabled --no-sort --delimiter : \
            --bind "change:reload:sleep 0.1; rg --line-number --color=always --no-heading --smart-case {q} || true" \
            --bind "ctrl-r:reload:sleep 0.1; rg --line-number --color=always --no-heading --smart-case {q} || true" \
            --preview 'bat --color=always --style=numbers --highlight-line {2} {1}' \
            --preview-window 'up,60%,border-bottom,+{2}+3/3' \
            --prompt '🔍 Search> ' \
            --height=90% \
            --layout=reverse \
            | tr -d '\r'
    )

    [[ -z "$result" ]] && return 0

    file=$(echo "$result" | cut -d: -f1)
    line=$(echo "$result" | cut -d: -f2)

    if [[ -n "$file" && -n "$line" ]]; then
        case "${EDITOR:-vi}" in
            nvim|vim) "${EDITOR:-vi}" "+${line}" "$file" ;;
            code) code -g "$file:$line" ;;
            emacsclient*) emacsclient -n "+$line:$file" ;;
            *) "${EDITOR:-vi}" "$file" ;;
        esac
    fi
}

pdedit() {
    local root file
    root=$(pd) || return

    if command -v fzf &>/dev/null && command -v fd &>/dev/null; then
        cd "$root" || return

        file=$(fd --type f --strip-cwd-prefix | fzf --no-multi --ansi | tr -d '\r')

        [ -z "$file" ] && return 0

        "${EDITOR:-vi}" "$file"
    else
        echo "pdedit requires fzf and fd to be installed" >&2
        return 1
    fi
}


pdswitch() {
    if command -v fzf &> /dev/null; then
        eval "$(pd switch)"
    else
        echo "pdswitch requires fzf to be installed" >&2
        return 1
    fi
}

pdjump() {
    local root action switch_out

    if ! command -v fzf &>/dev/null; then
        echo "pdjump requires fzf to be installed" >&2
        return 1
    fi

    switch_out=$(pd switch --print-only 2>/dev/null || pd switch 2>/dev/null)
    switch_out=${switch_out//$'\r'/}

    if [[ $switch_out == cd* ]]; then
        root=${switch_out#cd }
        root=${root%\"}
        root=${root#\"}
        root=${root%\'}
        root=${root#\'}
    else
        root=$switch_out
    fi

    if [[ ! -d $root ]]; then
        echo "Directory not found: $root" >&2
        return 1
    fi

    cd "$root" || return

    action=$(printf "Find files\nSearch text\nShow tree\n" |
        fzf --prompt 'Choose action> ' --ansi --height=30% --layout=reverse | tr -d '\r')

    [[ -z $action ]] && return 0

    case $action in
        "Find files") pdedit ;;
        "Search text") pdsearch ;;
        "Show tree") pdtree ;;
        *) echo "Unknown action: $action" ;;
    esac
}

# Project info commands (uses padofetch binary)
pdinfo() {
    padofetch info
}

pdhealth() {
    padofetch health
}

# Build/Test/Run commands - auto-detect build system
pdbuild() {
    local root
    root=$(pd) || return
    cd "$root" || return

    # Check for custom command in .pd.toml
    if [[ -f .pd.toml ]] && command -v toml &>/dev/null; then
        local custom_cmd
        custom_cmd=$(toml get .pd.toml commands.build 2>/dev/null)
        if [[ -n "$custom_cmd" ]]; then
            echo "Running custom build command: $custom_cmd"
            eval "$custom_cmd"
            return
        fi
    fi

    # Auto-detect build system
    if [[ -f Cargo.toml ]]; then
        cargo build
    elif [[ -f package.json ]]; then
        if [[ -f bun.lockb ]]; then
            bun run build
        elif [[ -f pnpm-lock.yaml ]]; then
            pnpm run build
        elif [[ -f yarn.lock ]]; then
            yarn build
        else
            npm run build
        fi
    elif [[ -f go.mod ]]; then
        go build ./...
    elif [[ -f Makefile ]]; then
        make
    elif [[ -f build.gradle || -f build.gradle.kts ]]; then
        ./gradlew build
    elif [[ -f pom.xml ]]; then
        mvn compile
    elif [[ -f pyproject.toml ]]; then
        if command -v poetry &>/dev/null; then
            poetry build
        elif command -v pip &>/dev/null; then
            pip install -e .
        fi
    else
        echo "No build system detected" >&2
        return 1
    fi
}

pdcompile() {
    pdbuild "$@"
}

pdtest() {
    local root
    root=$(pd) || return
    cd "$root" || return

    # Check for custom command in .pd.toml
    if [[ -f .pd.toml ]] && command -v toml &>/dev/null; then
        local custom_cmd
        custom_cmd=$(toml get .pd.toml commands.test 2>/dev/null)
        if [[ -n "$custom_cmd" ]]; then
            echo "Running custom test command: $custom_cmd"
            eval "$custom_cmd"
            return
        fi
    fi

    # Auto-detect test framework
    if [[ -f Cargo.toml ]]; then
        cargo test
    elif [[ -f package.json ]]; then
        if [[ -f bun.lockb ]]; then
            bun test
        elif [[ -f pnpm-lock.yaml ]]; then
            pnpm test
        elif [[ -f yarn.lock ]]; then
            yarn test
        else
            npm test
        fi
    elif [[ -f go.mod ]]; then
        go test ./...
    elif [[ -f Makefile ]]; then
        make test
    elif [[ -f build.gradle || -f build.gradle.kts ]]; then
        ./gradlew test
    elif [[ -f pom.xml ]]; then
        mvn test
    elif [[ -f pyproject.toml ]]; then
        if command -v poetry &>/dev/null; then
            poetry run pytest
        elif command -v pytest &>/dev/null; then
            pytest
        fi
    else
        echo "No test framework detected" >&2
        return 1
    fi
}

pdrun() {
    local root
    root=$(pd) || return
    cd "$root" || return

    # Check for custom command in .pd.toml
    if [[ -f .pd.toml ]] && command -v toml &>/dev/null; then
        local custom_cmd
        custom_cmd=$(toml get .pd.toml commands.run 2>/dev/null)
        if [[ -n "$custom_cmd" ]]; then
            echo "Running custom run command: $custom_cmd"
            eval "$custom_cmd"
            return
        fi
    fi

    # Auto-detect run command
    if [[ -f Cargo.toml ]]; then
        cargo run
    elif [[ -f package.json ]]; then
        if [[ -f bun.lockb ]]; then
            bun start
        elif [[ -f pnpm-lock.yaml ]]; then
            pnpm start
        elif [[ -f yarn.lock ]]; then
            yarn start
        else
            npm start
        fi
    elif [[ -f go.mod ]]; then
        go run .
    elif [[ -f Makefile ]]; then
        make run
    else
        echo "No run command detected" >&2
        return 1
    fi
}

# File operations
pdtree() {
    local root
    root=$(pd) || return
    if command -v tree &>/dev/null; then
        tree "$root"
    else
        echo "pdtree requires tree to be installed" >&2
        return 1
    fi
}

pdfiles() {
    local root
    root=$(pd) || return
    cd "$root" || return

    if command -v fd &>/dev/null; then
        fd --type f
    elif command -v find &>/dev/null; then
        find . -type f -not -path '*/\.*' -not -path '*/node_modules/*' -not -path '*/target/*'
    else
        echo "pdfiles requires fd or find to be installed" >&2
        return 1
    fi
}

# Dependency management
pddeps() {
    local root
    root=$(pd) || return
    cd "$root" || return

    if [[ -f Cargo.toml ]]; then
        cargo tree
    elif [[ -f package.json ]]; then
        if [[ -f bun.lockb ]]; then
            bun pm ls
        elif [[ -f pnpm-lock.yaml ]]; then
            pnpm list
        elif [[ -f yarn.lock ]]; then
            yarn list
        else
            npm list
        fi
    elif [[ -f go.mod ]]; then
        go list -m all
    elif [[ -f requirements.txt || -f pyproject.toml ]]; then
        pip list
    else
        echo "No dependency file detected" >&2
        return 1
    fi
}

pdoutdated() {
    local root
    root=$(pd) || return
    cd "$root" || return

    if [[ -f Cargo.toml ]]; then
        cargo outdated
    elif [[ -f package.json ]]; then
        if [[ -f bun.lockb ]]; then
            bun outdated
        elif [[ -f pnpm-lock.yaml ]]; then
            pnpm outdated
        elif [[ -f yarn.lock ]]; then
            yarn outdated
        else
            npm outdated
        fi
    elif [[ -f go.mod ]]; then
        go list -u -m all
    elif [[ -f requirements.txt || -f pyproject.toml ]]; then
        pip list --outdated
    else
        echo "No dependency file detected" >&2
        return 1
    fi
}


pdrecent() {
    if command -v fzf &> /dev/null; then
        eval "$(pd switch --recent)"
    else
        echo "pdrecent requires fzf to be installed" >&2
        return 1
    fi
}

pdstarred() {
    if command -v fzf &> /dev/null; then
        eval "$(pd switch --starred)"
    else
        echo "pdstarred requires fzf to be installed" >&2
        return 1
    fi
}

# Prompt integration - add to PS1/PROMPT for project info in prompt
# Example usage with powerlevel10k or oh-my-zsh themes:
#   PROMPT='%F{blue}$(pd_prompt)%f %~ %# '
# Or use precmd:
#   precmd() { RPROMPT="%F{blue}$(pd prompt 2>/dev/null)%f" }
pd_prompt() {
    local info
    info=$(pd prompt 2>/dev/null)
    if [ -n "$info" ]; then
        echo "$info"
    fi
}
