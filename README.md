# Pimp-My-Shell

[![Donate](https://img.shields.io/badge/Donate-PayPal-yellow.svg)](https://www.paypal.com/donate?business=YR6C4WB5CDZZL&no_recurring=0&item_name=contribute+to+open+source&currency_code=USD)
[![Go Report Card](https://goreportcard.com/badge/github.com/mr-pmillz/pimp-my-shell)](https://goreportcard.com/report/github.com/mr-pmillz/pimp-my-shell)
![GitHub all releases](https://img.shields.io/github/downloads/mr-pmillz/pimp-my-shell/total?style=social)
![GitHub repo size](https://img.shields.io/github/repo-size/mr-pmillz/pimp-my-shell?style=plastic)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mr-pmillz/pimp-my-shell?style=plastic)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/mr-pmillz/pimp-my-shell?style=plastic)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/mr-pmillz/pimp-my-shell?style=plastic)
[![Twitter](https://img.shields.io/twitter/url?style=social&url=https%3A%2F%2Fgithub.com%2Fmr-pmillz%2Fpimp-my-shell)](https://twitter.com/intent/tweet?text=Wow:&url=https%3A%2F%2Fgithub.com%2Fmr-pmillz%2Fpimp-my-shell)
[![CI](https://github.com/mr-pmillz/pimp-my-shell/actions/workflows/go.yml/badge.svg)](https://github.com/mr-pmillz/pimp-my-shell/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/mr-pmillz/pimp-my-shell/branch/master/graph/badge.svg?token=8NOCNMK0OE)](https://codecov.io/gh/mr-pmillz/pimp-my-shell)

Table of Contents
=================

* [Pimp\-My\-Shell](#pimp-my-shell)
  * [Install](#install)
  * [Usage](#usage)
  * [About](#about)
  * [Resources](#resources)
  * [Tmux Hotkeys](#tmux-hotkeys)
  * [VIM Hotkeys](#vim-hotkeys)
  * [Adjusting](#adjusting)
  * [Custom Aliases](#custom-aliases)
  * [Mac Fix Terminal bind keys](#mac-fix-terminal-bind-keys)
  * [Enjoy](#enjoy)


![pimp-my-shell.png](https://github.com/mr-pmillz/pimp-my-shell/blob/master/imgs/pimp-my-shell.png?raw=true)

## Install

```shell
go install github.com/mr-pmillz/pimp-my-shell/v2@latest
```

### Manual installation

Download the latest release for your system from [Releases](https://github.com/mr-pmillz/pimp-my-shell/releases)
- or clone the repo and run `go build` to build the binary.
  - If you're going to build from source, this project requires >= go v1.17.X
  - This project only works on MacOSX and Linux Ubuntu/Debian systems currently

### MacOS Users Ensure that you have x-code CommandLineTools installed

Vim YouCompleteMe plugin requires this for C lang completion

```shell
xcode-select --install
```

## Usage

```shell
./pimp-my-shell
```

If you already have oh-my-zsh installed, don't worry! Your ~/.zshrc file will not be overridden by the pimp-my-shell.
The only thing that will change is your zsh theme and the following plugins will be merged into your existing plugins=() object

- `git zsh-syntax-highlighting tmux zsh-autosuggestions virtualenv ansible docker docker-compose terraform helm kubectl fzf`

After Installation, if you want to Customize Powerlevel10k zsh theme differently, run

```shell
p10k configure
```

If you want your custom vim plugins to automatically update, simply create this cronjob

```shell
crontab -e
0 12 * * * cd ~/.vim_runtime/my_plugins && ./update.sh > gitPullUpdates.txt 2>&1
```

## About

This project was designed to automate all the configurations that I typically set up for my terminal on Macos and Debian/Ubuntu Linux.

Currently, this will (if not already installed and setup)

- install oh-my-zsh + awesome plugins
- install tmux + awesome mac config + plugins
- install vim + awesome vim setup + plugins
- install cheat + configure + community cheatsheets
- fzf + bat for finding files fast + file preview CTRL+r search history stupendously
- and various other dependencies

## Resources

Please see the following repos for more information about these configurations and plugins
All these configurations can be modified to your needs

- **Terminal**
  - **Fonts**
    - [NerdFonts](https://github.com/ryanoasis/nerd-fonts)
      - In Iterm2 Preferences, Profiles -> Text -> Font
      - ![enable-nerd-fonts-iterm2.png](https://github.com/mr-pmillz/pimp-my-shell/blob/master/imgs/enable-nerd-fonts-iterm2.png?raw=true)
  - **CLI Tools**
    - [lsd](https://github.com/Peltoche/lsd)
    - [fzf](https://github.com/junegunn/fzf)
      - [fzf Video](https://www.youtube.com/watch?v=qgG5Jhi_Els)
    - [bat](https://github.com/sharkdp/bat)
      - [bat Video](https://egghead.io/lessons/egghead-interactively-preview-files-with-fzf-and-bat-in-the-terminal)
    - [cheat](https://github.com/cheat/cheat)
      - useful command syntax cheatsheets Ex. `cheat tar`
    - [git-delta](https://github.com/dandavison/delta)
      - Beautiful less pager for `git diff`
    - [bpytop](https://github.com/aristocratos/bpytop)
      - Nice process monitor for the cli
- **Oh-My-ZSH**
  - **Theme**
  - [Powerlevel10k](https://github.com/romkatv/powerlevel10k)
    - This is an awesome theme for zsh
  - **Oh-My-ZSH Custom Plugins**
    - [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting)
    - [zsh-autosuggestions](https://github.com/zsh-users/zsh-autosuggestions)
- **TMUX**
  - [Oh-My-Tmux](https://github.com/gpakosz/.tmux)
  - **Tmux Plugins**
    - [Tmux Plugin Manager TPM](https://github.com/tmux-plugins/tpm)
    - [tmux-better-mouse-mode](https://github.com/NHDaly/tmux-better-mouse-mode)
    - [tmux-sensible](https://github.com/tmux-plugins/tmux-sensible)
    - [tmux-resurrect](https://github.com/tmux-plugins/tmux-resurrect)
    - [tmux-logging](https://github.com/tmux-plugins/tmux-logging)
- **VIM**
  - [The Ultimate vimrc](https://github.com/amix/vimrc)
    - This comes with a NERDTree Plugin by default
    - [NERDTree](https://github.com/preservim/nerdtree)
    - [ale](https://github.com/dense-analysis/ale)
    - [vim-dracula theme](https://github.com/dracula/vim)
  - **Custom VIM Plugins not included by default with Ultimate vimrc**
    - [nerdtree-git-plugin](https://github.com/Xuyuanp/nerdtree-git-plugin)
    - [YouCompleteMe](https://github.com/ycm-core/YouCompleteMe)
    - [vim-devicons](https://github.com/ryanoasis/vim-devicons)
    - [vim-visual-multi](https://github.com/mg979/vim-visual-multi)
    - [vim-yaml](https://github.com/stephpy/vim-yaml)
    - [vim-go](https://github.com/fatih/vim-go)
    - [vim rainbow highlighting](https://github.com/luochen1990/rainbow)
    - [fzf-vim](https://github.com/junegunn/fzf.vim)
    - [vim-helm](https://github.com/towolf/vim-helm)
    - [indentLine](https://github.com/Yggdroot/indentLine)
    - [lightline-bufferline](https://github.com/mengelbrecht/lightline-bufferline)
    - [vim-airline](https://github.com/vim-airline/vim-airline)
    - [vim-airline-themes](https://github.com/vim-airline/vim-airline-themes)

## Tmux Hotkeys

See [Tmux-Cheat-Sheet](https://tmuxcheatsheet.com/)

```shell
CTRL^b %     = split vertical
CTRL^b "     = split horizontal
CTRL^b h     = jump to left window
CTRL^b k     = jump to up window
CTRL^b c     = create new pane
CTRL^b ,     = rename pane
CTRL^b 1     = jump to 1 pane
CTRL^b I     = source tmux and install plugins
CTRL^b !     = open current window to new pane`
```

## VIM Hotkeys

```shell
,          = leader key <leader>
,nn        = toggle nerdtree
F12        = toggle nerdtree
i          = Nerdtree open pane horizontal
s          = Nerdtree open pane vertical
CTRL+ww    = cycle selected vim pane
,te        = open new tab after selecting file
,tb        = previous tab
,tn        = next tab
:bd        = buffer delete (similar to :q except it removes the tab buffer as well as closing the pane but will not quit)
,j         = jump to file with fzf fuzzy finder
```

## Adjusting

If you want to customize these configs further,
The main files you'll want to look at are the following

- ~/.tmux.conf.local
- ~/.vim_runtime/my_config.vim
  - This is where all further customization for vim can be done
    - It's the same as default .vimrc you would normally edit
  - Do not edit ~/.vimrc
- ~/.zshrc
  - You can edit this file with env var and aliases however, it is best a practice to put customizations
  - such as aliases in `~/.oh-my-zsh/custom/aliases.zsh`

## Custom Aliases

```shell
alias .z='source ~/.zshrc'
alias l='lsd -al --group-dirs first'
alias lld='ls -d -alh $PWD/*'
alias hg='history | grep'
alias myip='dig +short myip.opendns.com @resolver1.opendns.com'
alias pyup='python3 -m http.server'
alias a2='curl wttr.in/Ann_Arbor'
alias gs='git status'
alias gcmsg='git commit -m '
alias gl="git log --all --graph --pretty=tformat:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --date=short"
alias gll="git log --all --stat --pretty=tformat:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit --date=relative"
alias gln="git --no-pager log --all --stat --pretty=tformat:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%ad) %C(bold blue)<%an>%Creset' --abbrev-commit --date=relative -n 10"
alias clp="pbcopy < $1"
alias fzfbat="fzf --preview 'bat --style numbers,changes --color=always {}' | head -500"
```

## Mac Fix Terminal bind keys

- because of a shortcut conflict with Mission Control/Spaces on MacOSX
- make sure to uncheck these 2 options in
- System Preferences -> Keyboard -> Shortcuts -> Mission Control, Move left/right a space
  ![mac-bind-keys.png](https://github.com/mr-pmillz/pimp-my-shell/blob/master/imgs/mac-bind-keys.png?raw=true)


## Enjoy
