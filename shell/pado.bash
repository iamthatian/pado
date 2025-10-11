# Pado shell integration for bash
# Add this to your ~/.bashrc or ~/.bash_profile:
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

# pdinfo() {
#     pd info
# }
#
# pdbuild() {
#     pd build
# }
#
# pdtest() {
#     pd test
# }
#
# pdrun() {
#     pd run
# }
#
# pdstats() {
#     pd stats
# }
#
# pdrec() {
#     pd recent "$@"
# }
#
# pdstar() {
#     pd star
# }
# pdtree() {
#     pd tree
# }
#
# pdcd() {
#     eval "$(pd cd)"
# }

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

# Prompt integration - add to PS1 for project info in prompt
# Example usage:
#   PS1='[\[\e[34m\]$(pd_prompt)\[\e[0m\]] \w \$ '
# Or simpler:
#   PS1='[$(pd_prompt)] \w \$ '
pd_prompt() {
    local info
    info=$(pd prompt 2>/dev/null)
    if [ -n "$info" ]; then
        echo "$info"
    fi
}
