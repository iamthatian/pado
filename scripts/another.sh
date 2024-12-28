#!/bin/sh

fools() {
   # c="fd -L --color=always --full-path ."
   c="ls"
   # [[ ! -f "$data/deep" ]] && c=" $c --max-depth 1" || c=" $c --max-depth 3"
   # [[ -f "$data/hidden" ]] && c=" $c -H"
   $c
}

export -f fools

gothere() {
  cd $1
}

export -f gothere
d=""
# s=`ls | fzf --bind="start:$fools"`
s=`ls | fzf`
while true
do
  # SELECTED=`ls -d */* | fzf | awk '{print "\"" $0 "\""}'`
  # SELECTED=`ls | fzf | awk '{print "\"" $0 "\""}'`
  # s=`ls | fzf -- `
  # s=`ls | fzf --bind="start:$fools"`
  command -v $s

   # c="fd -L --color=always --full-path ."
  # if [ -z "$SELECTED" ]
  # then
  #     read action
  #     $action
  #     continue
  # fi
  # d="$(echo "$s" | tail -n1)"
  d=$(echo "$s")
  break
  if [[ -d $d ]]; then
      echo $d
      cd $d
      # || continue;
      echo "Fuck"
      echo "Fuck"
      # break
  elif [[ -f $d ]]; then
      s=`rg . | fzf`
      # $EDITOR $d
      echo "FUCK"
      echo "ME"
      # break
  else
      break
  fi

  # clear
  # echo $SELECTED
  # printf "$ "
  # read action
  # if [ "$action" = "vi" ]; then
  #     # session="Code"
  #     # tmux new-session -d -s $session
  #     # tmux new-window -t $session:1 -n 'NVim'
  #     # tmux send-keys -t 'NVim' "nvim $SELECTED" C-m
  #     # tmux attach-session -t $SESSION:1
  #     break
  # else
  #     eval $action $SELECTED
  #     break
  # fi
done

if [[ -d $d ]]; then
    # echo $d
    # cd $d
    gothere $d
    # || continue;
    # echo "Fuck"
    # echo "Fuck"
    # break
elif [[ -f $d ]]; then
    # s=`rg . | fzf`
    $EDITOR $d
    # echo "FUCK"
    # echo "ME"
    # break
# else
#     break
fi

