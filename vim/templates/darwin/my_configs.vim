try
    " vim-devicons set encoding utf-8
    set encoding=UTF-8
    syntax enable
    filetype plugin indent on
    "Set Dracula to the default Theme
    color dracula
    let g:dracula_colorterm = 0
    let g:dracula_italic = 0
    " Show Line Numbers
    set nu
    " Set Background to dark mode
    set background=dark

    " disable autopairs
    let g:AutoPairs = {}
    "
    " Disable MarkDown Folding Because it's slow af.
    let g:vim_markdown_folding_disabled = 1

    " YCM Settings
	"YouCompleteMe Fix Tabs and disable linter / diagnostics
	set nosmarttab
    let g:ycm_show_diagnostics_ui = 0
    let g:ycm_enable_diagnostic_signs = 0
    let g:ycm_enable_diagnostic_highlighting = 0
    " Disable Docs Annoying Split window in vim from YouCompleteMe Plugin
    set completeopt-=preview

    " ALE Settings, change message format, error symbols and only lint when file is saved
    let g:ale_echo_msg_format = '[%linter%] %s [%severity%]'
    let g:ale_sign_error = '✘'
    let g:ale_sign_warning = '⚠'
    let g:ale_lint_on_text_changed = 'never'

	"Bracket Pair Colorizer enable
	let g:rainbow_active = 1

	" Better Python Syntax Highlighting
	let g:python_highlight_all = 1

    " Vim-go Golang customizations
    let g:go_highlight_extra_types = 1
    let g:go_highlight_functions = 1
    let g:go_highlight_function_parameters = 1
    let g:go_highlight_function_calls = 1
    let g:go_highlight_types = 1
    let g:go_highlight_fields = 1
    let g:go_highlight_build_constraints = 1
    let g:go_highlight_format_strings = 1

    " vim-airline settings
    let g:airline_theme = 'badwolf'
    let g:airline#extensions#tabline#enabled = 1
    let g:airline#extensions#syntastic#enabled = 1
    let g:airline#extensions#branch#enabled = 1
    let g:airline#extensions#tagbar#enabled = 1
    let g:airline_skip_empty_sections = 1

    " NERDTree
    " map Nerdtree to F2 or can use leader key which is a comma ,
    " ,nn
    map <F2> :NERDTreeToggle<CR>
    let g:NERDTreeWinPos = "left"
    let NERDTreeQuitOnOpen = 1
    " let NERDTreeMinimalUI = 1
    let NERDTreeDirArrows = 1
    let NERDTreeShowHidden = 1

    " Set ssh config syntax highlighting
    au BufNewFile,BufRead ssh_config,*/.ssh/*  setf sshconfig
    " YAML Settings
    " Set syntax highlighting for all yaml,yml files
    au! BufNewFile,BufReadPost *.{yaml,yml} set filetype=yaml foldmethod=indent
    au BufNewFile,BufReadPost config,*/.kube/config setf yaml
    " Set ncie 2-space YAML as default
    autocmd FileType yaml setlocal ts=2 sts=2 sw=2 expandtab
    " Show indentation line
    let g:indentLine_char = '⦙'
    autocmd FileType yaml execute
      \'syn match yamlBlockMappingKey /^\s*\zs.*\ze\s*:\%(\s\|$\)/'
    " END YAML Settings
    
    " prevent vim from hiding quotes in json files
    let g:vim_json_conceal=0 

    " disable folding
    set nofoldenable

     " Fix indentLine plugin for markdown backticks and invisible chars
     let g:indentLine_fileTypeExclude = ['markdown']

    " Tab hotkeys
    map <leader>tb :tabprevious<cr>
    map <leader>tn :tabnext<cr>

    " lightline-bufferline config
    let g:lightline#bufferline#enable_devicons = 1
    let g:lightline#bufferline#show_number  = 1
    let g:lightline = {
      \ 'colorscheme': 'one',
      \ 'active': {
      \   'left': [ [ 'mode', 'paste' ], [ 'readonly', 'filename', 'modified' ] ]
      \ },
      \ 'tabline': {
      \   'left': [ ['buffers'] ],
      \   'right': [ ['close'] ]
      \ },
      \ 'component_expand': {
      \   'buffers': 'lightline#bufferline#buffers'
      \ },
      \ 'component_type': {
      \   'buffers': 'tabsel'
      \ }
      \ }
    if
        has('gui_running') set guioptions-=e
    endif
    let g:lightline#bufferline#clickable = 1
    " in order for tabs to be clickable, need to set mouse to a mode
    "FZF settings
    set rtp+=/usr/local/opt/fzf

    " fix python formatting
	au BufNewFile,BufRead *.py
				\ setlocal tabstop=4 |
				\ setlocal softtabstop=4 |
				\ setlocal shiftwidth=4 |
				\ setlocal textwidth=120 |
				\ setlocal expandtab |
				\ setlocal autoindent |
				\ setlocal fileformat=unix

    " fix paste settings
	let g:yankstack_map_keys = 0

    " Disbale visual mode warnings from vim-visual-multi
    let g:VM_show_warnings = 0

catch
endtry
