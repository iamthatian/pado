# shellcheck shell=bash
function __sp_cd() {
    # shellcheck disable=SC2164
    \builtin cd -- "$@"
}

function __sp_compile() {
    command sp compile
}

function __sp_cd_root() {
    __sp_cd $(sp)
}

function __sp_grep_edit() {
    RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case"
    INITIAL_QUERY="${*:-}"
    command fzf --ansi --disabled --query "$INITIAL_QUERY" --bind="ctrl-c:abort" --bind "start:reload:$RG_PREFIX {q} $(sp)" --bind "change:reload:sleep 0.1; $RG_PREFIX {q} || true" --delimiter : --preview 'bat --color=always {1} --highlight-line {2}' --preview-window 'up,60%,border-bottom,+{2}+3/3,~3' --bind 'enter:become($EDITOR {1} +{2})'
}

function __sp_find_file() {
    \builtin local result
    result="$(\command fd . $(sp) | fzf -- "$@")"

    if [[ -f ${result} ]]; then;
        $EDITOR "${result}"
    fi

    if [[ -d ${result} ]]; then;
        __sp_cd "${result}"
    fi
}

# TODO show selected project on selection
function __sp_find_other_project() {
    s="$(\command sp list | fzf --bind="ctrl-c:abort" -- "$@")"
    i=0
    while true
    do
        command -v $s
        d=$(echo "$s")
        if (( $i == 0 )); then
            if [[ $d == "" ]]; then
                break
            fi
            selection=$d
            menu=$(\command printf "find(edit)\nfind(show)\ngrep(edit)\ngo to project\n" | fzf --bind="ctrl-c:abort")
            s=$menu "$@"
        elif (( $i == 1 )); then
            if [[ $d == "" ]]; then
                break
            fi
            action=$d
            if [[ $action == "find(edit)" ]]; then
                s="$(fd . $(sp $selection) | fzf --bind="ctrl-c:abort" -- "$@")"
            elif [[ $action == "find(show)" ]]; then
                command fd . $(sp $selection)
                break
            elif [[ $action == "grep(edit)" ]]; then
                # TODO: This switches to a different directory on query wtf
                RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case"
                INITIAL_QUERY="${*:-}"
                __sp_cd ${selection}
                command fzf --ansi --disabled --query "$INITIAL_QUERY" --bind="ctrl-c:abort" --bind "start:reload:$RG_PREFIX {q} $(sp $selection)" --bind "change:reload:sleep 0.1; $RG_PREFIX {q} || true" --delimiter : --preview 'bat --color=always {1} --highlight-line {2}' --preview-window 'up,60%,border-bottom,+{2}+3/3,~3' --bind 'enter:become($EDITOR {1} +{2})'
                break
            elif [[ $action == "go to project" ]]; then
                __sp_cd ${selection}
                break
            else
                break
            fi
        elif (( $i == 2 )); then
            if [[ $d == "" ]]; then
                break
            fi
            __sp_cd ${selection}
            if [[ -f $d ]]; then;
                $EDITOR $d
            fi

            if [[ -d $d ]]; then;
                __sp_cd $d
            fi

            break
        else
            break
        fi
        ((i++));
    done
}

# Shortcuts
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
    __sp_find_other_project "$@"
}
