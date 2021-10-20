# export PATH=$HOME/bin:/usr/local/bin:$PATH

# Path to your oh-my-zsh installation.
export ZSH=$HOME/.oh-my-zsh

# See https://github.com/ohmyzsh/ohmyzsh/wiki/Themes
ZSH_THEME="robbyrussell"

# Uncomment the following line to use case-sensitive completion.
# CASE_SENSITIVE="true"

plugins=(
        git
        tmux
        fakename
        cool-test
        zsh-syntax-highlighting
        zsh-autosuggestions
        virtualenv
        ansible
        docker
        docker-compose
        terraform
        fzf
        another-plugin
        someplugin
        coolplugin
        testing
)

source $ZSH/oh-my-zsh.sh

# User configuration

# export MANPATH="/usr/local/man:$MANPATH"

# Example aliases
# alias ohmyzsh="mate ~/.oh-my-zsh"
