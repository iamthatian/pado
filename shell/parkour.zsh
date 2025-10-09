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
