#!/usr/bin/env bash
# rg --column --line-number --no-heading --color=always --smart-case' \\
#     fzf --bind 'start:reload:$rg_prefix ""' \\
#         --bind 'change:reload:$rg_prefix {q} || true' \\
#         --bind 'enter:become(vim {1} +{2})' \\
#         --ansi --disabled \
#

# export rg_prefix='rg'
# fzf --bind 'start:reload:$rg_prefix ""' --bind 'change:reload:$rg_prefix {q} || true' --bind 'enter:become(vim {1} +{2})' --height=50% --layout=reverse

    # fzf --bind 'start:reload:$rg_prefix ""' \\
    #     --bind 'change:reload:$rg_prefix {q} || true' \\
    #     --bind 'enter:become(vim {1} +{2})' \\
    #     --ansi --disabled \\
    #     --height=50% --layout=reverse\

RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case "
INITIAL_QUERY="${*:-}"
fzf --ansi --disabled --query "$INITIAL_QUERY" --bind "start:reload:$RG_PREFIX {q}" --bind "change:reload:sleep 0.1; $RG_PREFIX {q} || true" --delimiter : --preview 'bat --color=always {1} --highlight-line {2}' --preview-window 'up,60%,border-bottom,+{2}+3/3,~3' --bind 'enter:become($EDITOR {1} +{2})'
