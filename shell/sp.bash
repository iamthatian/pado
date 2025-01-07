function __parkour_pwd() {
    \builtin pwd -L
}

function __parkour_cd() {
    # shellcheck disable=SC2164
    \builtin cd -- "$@"
}

function __parkour_compile() {
    command parkour compile
}

function __parkour_cd_root() {
    __parkour_cd $(pk)
}

function __parkour_grep_edit() {
    pk grep-file
}

function __parkour_find_edit() {
    pk find-file
}

# if not currently in a project (give option to add current directory/project or to go to a project, in the later case this is ran)
# TODO show selected project on selection
function __parkour_find_project() {
    pk find
    local exit_code=$?

    # If it exited successfully and we need to cd (go to project was selected)
    if [ -f "/tmp/pk_last_dir" ]; then
        local dir
        dir=$(cat "/tmp/pk_last_dir")
        if [ -d "$dir" ]; then
            cd "$dir"
        fi
        rm -f "/tmp/pk_last_dir"
    fi

    return $exit_code
}

function pkc() {
    __parkour_compile "$@"
}

function pkf() {
    __parkour_find_file "$@"
}

function pkg() {
    __parkour_grep_edit "$@"
}

# Go to project root
function pkr() {
    __parkour_cd_root "$@"
}

function pkp() {
    __parkour_find_project "$@"
}
