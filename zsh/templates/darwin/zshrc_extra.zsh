# Add ~/go/bin to path
[[ ":$PATH:" != *":$HOME/go/bin:"* ]] && export PATH="${PATH}:$HOME/go/bin"

# Set CTRL+U to only delete backwards to the left of the cursor.
bindkey \^U backward-kill-line

# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.
[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh
POWERLEVEL9K_INSTANT_PROMPT=quiet
