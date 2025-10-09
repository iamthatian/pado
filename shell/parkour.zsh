# Parkour shell integration for zsh
# Add this to your ~/.zshrc:
#   eval "$(pk init)"

pkroot() {
    local root
    root=$(pk) || return
    cd "$root" || return
}

pkcd() {
    eval "$(pk cd)"
}

pkfind() {
    local root
    root=$(pk) || return
    if command -v fzf &> /dev/null && command -v fd &> /dev/null; then
        cd "$root" && fd --type f | fzf
    else
        echo "pkfind requires fzf and fd to be installed" >&2
        return 1
    fi
}

pkgrep() {
    local root
    root=$(pk) || return
    if command -v rg &> /dev/null; then
        cd "$root" && rg "$@"
    else
        echo "pksearch requires ripgrep (rg) to be installed" >&2
        return 1
    fi
}

pksearch() {
    local root result file line
    root=$(pk) || return

    if ! command -v rg &>/dev/null || ! command -v fzf &>/dev/null; then
        echo "pksearch requires both ripgrep (rg) and fzf" >&2
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

pkedit() {
    local root file
    root=$(pk) || return

    if command -v fzf &>/dev/null && command -v fd &>/dev/null; then
        cd "$root" || return

        file=$(fd --type f --strip-cwd-prefix | fzf --no-multi --ansi | tr -d '\r')

        [ -z "$file" ] && return 0

        "${EDITOR:-vi}" "$file"
    else
        echo "pkedit requires fzf and fd to be installed" >&2
        return 1
    fi
}

pktree() {
    pk tree
}

pkswitch() {
    if command -v fzf &> /dev/null; then
        eval "$(pk switch)"
    else
        echo "pkswitch requires fzf to be installed" >&2
        return 1
    fi
}

pkjump() {
    local root action switch_out

    if ! command -v fzf &>/dev/null; then
        echo "pkjump requires fzf to be installed" >&2
        return 1
    fi

    switch_out=$(pk switch --print-only 2>/dev/null || pk switch 2>/dev/null)
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
        "Find files") pkedit ;;
        "Search text") pksearch ;;
        "Show tree") pktree ;;
        *) echo "Unknown action: $action" ;;
    esac
}

pkinfo() {
    pk info
}

pkbuild() {
    pk build
}

pktest() {
    pk test
}

pkrun() {
    pk run
}

pkstats() {
    pk stats
}

pkrec() {
    pk recent "$@"
}

pkstar() {
    pk star
}

pkrecent() {
    if command -v fzf &> /dev/null; then
        eval "$(pk switch --recent)"
    else
        echo "pkrecent requires fzf to be installed" >&2
        return 1
    fi
}

pkstarred() {
    if command -v fzf &> /dev/null; then
        eval "$(pk switch --starred)"
    else
        echo "pkstarred requires fzf to be installed" >&2
        return 1
    fi
}

# Prompt integration - add to PS1/PROMPT for project info in prompt
# Example usage with powerlevel10k or oh-my-zsh themes:
#   PROMPT='%F{blue}$(pk_prompt)%f %~ %# '
# Or use precmd:
#   precmd() { RPROMPT="%F{blue}$(pk prompt 2>/dev/null)%f" }
pk_prompt() {
    local info
    info=$(pk prompt 2>/dev/null)
    if [ -n "$info" ]; then
        echo "$info"
    fi
}
