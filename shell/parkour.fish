# Parkour shell integration for fish
# Add this to your ~/.config/fish/config.fish:
#   pk init | source

function pkroot
    set -l root (pk)
    or return
    cd $root
end

function pkcd
    eval (pk cd)
end

function pkfind
    set -l root (pk)
    or return
    if command -v fzf &> /dev/null; and command -v fd &> /dev/null
        cd $root && fd --type f | fzf
    else
        echo "pkfind requires fzf and fd to be installed" >&2
        return 1
    end
end

function pkgrep
    set -l root (pk)
    or return
    if command -v rg &> /dev/null
        cd $root && rg $argv
    else
        echo "pksearch requires ripgrep (rg) to be installed" >&2
        return 1
    end
end

function pksearch
    set -l root (pk)
    or return

    if not type -q rg; or not type -q fzf
        echo "pksearch requires both ripgrep (rg) and fzf" >&2
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

function pkedit
    set -l root (pk)
    or return

    if command -v fzf >/dev/null; and command -v fd >/dev/null
        cd $root || return

        set -l file (fd --type f --strip-cwd-prefix | fzf --no-multi --ansi | tr -d '\r')

        if test -n "$file"
            eval "$EDITOR" "$file"
        end
    else
        echo "pkedit requires fzf and fd to be installed" >&2
        return 1
    end
end

function pktree
    pk tree
end

function pkswitch
    if command -v fzf &> /dev/null
        eval (pk switch)
    else
        echo "pkswitch requires fzf to be installed" >&2
        return 1
    end
end

function pkjump
    if not type -q fzf
        echo "pkjump requires fzf to be installed" >&2
        return 1
    end

    set -l switch_out (pk switch --print-only 2>/dev/null; or pk switch 2>/dev/null)
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
            pkedit
        case "Search text"
            pksearch
        case "Show tree"
            pktree
        case '*'
            echo "Unknown action: $action"
    end
end

function pkinfo
    pk info
end

function pkbuild
    pk build
end

function pktest
    pk test
end

function pkrun
    pk run
end

function pkstats
    pk stats
end

function pkrec
    pk recent $argv
end

function pkstar
    pk star
end

function pkrecent
    if command -v fzf &> /dev/null
        eval (pk switch --recent)
    else
        echo "pkrecent requires fzf to be installed" >&2
        return 1
    end
end

function pkstarred
    if command -v fzf &> /dev/null
        eval (pk switch --starred)
    else
        echo "pkstarred requires fzf to be installed" >&2
        return 1
    end
end

# Prompt integration for fish shell
# Example usage - add to your fish_prompt function:
#   function fish_prompt
#       set_color blue
#       echo -n (pk_prompt)
#       set_color normal
#       echo -n ' '
#       echo -n (prompt_pwd)' > '
#   end
function pk_prompt
    set -l info (pk prompt 2>/dev/null)
    if test -n "$info"
        echo -n "[$info]"
    end
end
