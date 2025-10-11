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

# function pdtree
#     pd tree
# end
#
#
# function pdcd
#     eval (pd cd)
# end
#
# function pdinfo
#     pd info
# end
#
# function pdbuild
#     pd build
# end
#
# function pdtest
#     pd test
# end
#
# function pdrun
#     pd run
# end
#
# function pdstats
#     pd stats
# end
#
# function pdrec
#     pd recent $argv
# end
#
# function pdstar
#     pd star
# end

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
