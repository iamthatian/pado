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

function pksearch
    set -l root (pk)
    or return
    if command -v rg &> /dev/null
        cd $root && rg $argv
    else
        echo "pksearch requires ripgrep (rg) to be installed" >&2
        return 1
    end
end

function pkedit
    set -l root (pk)
    or return
    if command -v fzf &> /dev/null; and command -v fd &> /dev/null
        set -l file (cd $root && fd --type f | fzf)
        if test -n "$file"
            eval $EDITOR $root/$file
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
