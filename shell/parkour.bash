# Parkour shell integration for bash
# Add this to your ~/.bashrc or ~/.bash_profile:
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

pksearch() {
    local root
    root=$(pk) || return
    if command -v rg &> /dev/null; then
        cd "$root" && rg "$@"
    else
        echo "pksearch requires ripgrep (rg) to be installed" >&2
        return 1
    fi
}

pkedit() {
    local root file
    root=$(pk) || return
    if command -v fzf &> /dev/null && command -v fd &> /dev/null; then
        file=$(cd "$root" && fd --type f | fzf)
        [ -n "$file" ] && ${EDITOR:-vi} "$root/$file"
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

# Prompt integration - add to PS1 for project info in prompt
# Example usage:
#   PS1='[\[\e[34m\]$(pk_prompt)\[\e[0m\]] \w \$ '
# Or simpler:
#   PS1='[$(pk_prompt)] \w \$ '
pk_prompt() {
    local info
    info=$(pk prompt 2>/dev/null)
    if [ -n "$info" ]; then
        echo "$info"
    fi
}
