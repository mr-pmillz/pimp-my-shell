alias .z='source ~/.zshrc'
alias l='lsd -al --group-dirs first'
alias lld='ls -d -alh $PWD/*'
alias hg='history | grep'
alias myip='dig +short myip.opendns.com @resolver1.opendns.com'
alias pyup='python3 -m http.server'
alias a2='curl wttr.in/Ann_Arbor'
alias gs='git status'
alias nstat='netstat -p tcp -van | grep LISTEN'
alias gl="git log --all --graph --pretty=tformat:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --date=short"
alias gll="git log --all --stat --pretty=tformat:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit --date=relative"
alias gln="git --no-pager log --all --stat --pretty=tformat:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%ad) %C(bold blue)<%an>%Creset' --abbrev-commit --date=relative -n 10"
alias clp="pbcopy < $1"
alias fzfbat="fzf --preview 'bat --style numbers,changes --color=always {}' | head -500"
alias ip='ip --color=auto'
