# zsh-syntax-highlighting colors
ZSH_HIGHLIGHT_HIGHLIGHTERS=(
    main
    brackets
)
ZSH_HIGHLIGHT_STYLES[single-hyphen-option]=fg=magenta
ZSH_HIGHLIGHT_STYLES[double-hyphen-option]=fg=magenta

# For linux add ~/.local/bin to path
[[ ":$PATH:" != *":${HOME}/.locaal/bin:"* ]] && export PATH="${PATH}:${HOME}/.local/bin"

# Add ~/go/bin to path
[[ ":$PATH:" != *":${HOME}/go/bin:"* ]] && export PATH="${PATH}:${HOME}/go/bin"
# Set GOPATH
if [[ -z "${GOPATH}" ]]; then export GOPATH="${HOME}/go"; fi
[[ ":$PATH:" != *":/usr/local/go/bin:"* ]] && export PATH="${PATH}:/usr/local/go/bin"

# Set CTRL+U to only delete backwards to the left of the cursor.
bindkey \^U backward-kill-line

# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.
[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh
POWERLEVEL9K_INSTANT_PROMPT=quiet
