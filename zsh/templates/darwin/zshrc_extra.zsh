# Add ~/go/bin to path
[[ ":$PATH:" != *":${HOME}/go/bin:"* ]] && export PATH="${PATH}:$HOME/go/bin"
# Add gsed to path
[[ ":$PATH:" != *":/usr/local/opt/gnu-sed/libexec/gnubin:" ]] && export PATH="${PATH}:/usr/local/opt/gnu-sed/libexec/gnubin"

# Set GOPATH, GOROOT and add GOROOT to path
if [[ -z "${GOPATH}" ]]; then export GOPATH="${HOME}/go"; fi
if [[ -z "${GOROOT}" ]]; then export GOROOT="$(brew --prefix golang)/libexec"; fi
[[ ":$PATH:" != *":${GOROOT}/bin:"* ]] && export PATH="${PATH}:${GOROOT}/bin"

# Set CTRL+U to only delete backwards to the left of the cursor.
bindkey \^U backward-kill-line

# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.
[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh
POWERLEVEL9K_INSTANT_PROMPT=quiet
