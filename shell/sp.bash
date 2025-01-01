function __sp_pwd() {
    \builtin pwd -L
}

function __sp_cd() {
    # shellcheck disable=SC2164
    \builtin cd -- "$@"
}

function __sp_compile() {
    command sp compile
}

function __sp_cd_root() {
    __sp_cd $(pk)
}

function __sp_grep_edit() {
    rm -f /tmp/rg-fzf-{r,f}
    RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case "
    INITIAL_QUERY="${*:-}"
    fzf --ansi --disabled --query "$INITIAL_QUERY" \
        --bind "start:reload:$RG_PREFIX {q}" \
        --bind "change:reload:sleep 0.1; $RG_PREFIX {q} || true" \
        --bind 'ctrl-t:transform:[[ ! $FZF_PROMPT =~ ripgrep ]] &&
        echo "rebind(change)+change-prompt(1. ripgrep> )+disable-search+transform-query:echo \{q} > /tmp/rg-fzf-f; cat /tmp/rg-fzf-r" ||
        echo "unbind(change)+change-prompt(2. fzf> )+enable-search+transform-query:echo \{q} > /tmp/rg-fzf-r; cat /tmp/rg-fzf-f"' \
        --color "hl:-1:underline,hl+:-1:underline:reverse" \
        --prompt '1. ripgrep> ' \
        --delimiter : \
        --header 'CTRL-T: Switch between ripgrep/fzf' \
        --preview 'bat --color=always {1} --highlight-line {2}' \
        --preview-window 'up,60%,border-bottom,+{2}+3/3,~3' \
        --bind 'enter:become(vim {1} +{2})'
    }
    # Grep edit
    # function __sp_grep_edit() {
    #     RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case"
    #     INITIAL_QUERY="${*:-}"
    #     __sp_cd ${selection}
    #     command fzf --ansi --disabled --query "$INITIAL_QUERY" --bind="ctrl-c:abort" --bind "start:reload:$RG_PREFIX . $(pk $selection)" --bind "change:reload:sleep 0.1; $RG_PREFIX . || true" --delimiter : --preview 'bat --color=always {1} --highlight-line {2}' --preview-window 'up,60%,border-bottom,+{2}+3/3,~3' --bind 'enter:become($EDITOR {1} +{2})'
    #
    #     # RG_STUFF="rg --column --line-number --no-heading --color=always --smart-case"
    #     # INITIAL_QUERY="${*:-}"
    #     # echo fzf --ansi --disabled --query "$INITIAL_QUERY" --bind="ctrl-c:abort" --bind "start:reload:$RG_STUFF {q} $(pk $1)" --bind "change:reload:sleep 0.1; $RG_STUFF || true" --delimiter : --preview 'bat --color=always {1} --highlight-line {2}' --preview-window 'up,60%,border-bottom,+{2}+3/3,~3' --bind 'enter:become($EDITOR {1} +{2})'
    #     # fzf --ansi --disabled --query "$INITIAL_QUERY" --bind="ctrl-c:abort" --bind "start:reload:$RG_STUFF {q} $(pk $1)" --bind "change:reload:sleep 0.1; $RG_STUFF || true" --delimiter : --preview 'bat --color=always {1} --highlight-line {2}' --preview-window 'up,60%,border-bottom,+{2}+3/3,~3' --bind 'enter:become($EDITOR {1} +{2})'
    # }

function __sp_find_edit() {
    FD_PREFIX="fd . -tf"
    INITIAL_QUERY="${*:-}"
    __sp_cd $1
    fzf --ansi --disabled --query "$INITIAL_QUERY" --bind="ctrl-c:abort" --bind "start:reload:$FD_PREFIX $(pk $1)" --bind "change:reload:sleep 0.1; $RG_PREFIX {q} || true" --delimiter : --preview 'bat --color=always {1} --highlight-line {2}' --preview-window 'up,60%,border-bottom,+{2}+3/3,~3' --bind 'enter:become($EDITOR {1} +{2})'
}

function __sp_find_file() {
    \builtin local result
    result="$(\command fd . $(pk) -tf | fzf -- "$@")"

    if [[ -f ${result} ]]; then;
        $EDITOR "${result}"
    fi

    if [[ -d ${result} ]]; then;
        __sp_cd "${result}"
    fi
}

# if not currently in a project (give option to add current directory/project or to go to a project, in the later case this is ran)
# TODO show selected project on selection
function __sp_find_project() {
    s="$(\command pk list | fzf --bind="ctrl-c:abort" -- "$@")"
    command -v $s
    selection=$(echo "$s")
    # set -E

    if [[ $selection == "" ]]; then
        return
    fi

    menu=$(\command printf "find(edit)\nfind(show)\ngrep(edit)\ngo to project\n" | fzf --bind="ctrl-c:abort")
    command -v $menu
    output=$(echo "$menu")

    if [[ $output == "" ]]; then
        return
    fi

    if [[ $output == "find(edit)" ]]; then
        __sp_find_edit $selection
        return
    elif [[ $output == "find(show)" ]]; then
        command fd . $(pk $selection)
        return
    elif [[ $output == "grep(edit)" ]]; then
        rm -f /tmp/rg-fzf-{r,f}
        RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case "
        INITIAL_QUERY="${*:-}"
        fzf --ansi --disabled --query "$INITIAL_QUERY" \
            --bind "start:reload:$RG_PREFIX {q} $(pk $selection)" \
            --bind "change:reload:sleep 0.1; $RG_PREFIX {q} $(pk $selection) || true"
        return
    elif [[ $output == "go to project" ]]; then
        __sp_cd ${selection}
        return
    else
        return
    fi
}

function spc() {
    __sp_compile "$@"
}

function spf() {
    __sp_find_file "$@"
}

function spg() {
    __sp_grep_edit "$@"
}

# Go to project root
function spr() {
    __sp_cd_root "$@"
}

function spp() {
    __sp_find_project "$@"
}
