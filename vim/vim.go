package vim

import (
	"embed"
	"fmt"
	"github.com/google/periph/host/distro"
	"pimp-my-shell/localio"
)

//go:embed templates/*
var myConfigs embed.FS

// InstallVimPlugins ...
func InstallVimPlugins(osType string, dirs *localio.Directories) error {
	if exists, err := localio.Exists(fmt.Sprintf("%s/.vim_runtime", dirs.HomeDir)); err == nil && exists == true {
		return nil
	}
	// install awesome vim
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone --depth=1 https://github.com/amix/vimrc.git %s/.vim_runtime", dirs.HomeDir)); err != nil {
		return err
	}
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("cd %s && bash .vim_runtime/install_awesome_vimrc.sh", dirs.HomeDir)); err != nil {
		return err
	}

	// install YCM
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/ycm-core/YouCompleteMe.git %s/.vim_runtime/my_plugins/YouCompleteMe", dirs.HomeDir)); err != nil {
		return err
	}

	pythonPath, _ := localio.CommandExists("python3")
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("cd %s/.vim_runtime/my_plugins/YouCompleteMe && git submodule update --init --recursive && %s install.py --all || true", dirs.HomeDir, pythonPath)); err != nil {
		return err
	}
	// vim-yaml plugin
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/stephpy/vim-yaml.git %s/.vim_runtime/my_plugins/vim-yaml", dirs.HomeDir)); err != nil {
		return err
	}
	// vim-go plugin
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/fatih/vim-go.git %s/.vim_runtime/my_plugins/vim-go", dirs.HomeDir)); err != nil {
		return err
	}
	// rainbow brackets vim plugin
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/luochen1990/rainbow.git %s/.vim_runtime/my_plugins/rainbow", dirs.HomeDir)); err != nil {
		return err
	}
	// fzf.vim plugin
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/junegunn/fzf.vim.git %s/.vim_runtime/my_plugins/fzf.vim", dirs.HomeDir)); err != nil {
		return err
	}
	// nerdtree-git-plugin
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/Xuyuanp/nerdtree-git-plugin.git %s/.vim_runtime/my_plugins/nerdtree-git-plugin", dirs.HomeDir)); err != nil {
		return err
	}
	// vim-devicons plugin
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/ryanoasis/vim-devicons.git %s/.vim_runtime/my_plugins/vim-devicons", dirs.HomeDir)); err != nil {
		return err
	}
	// lightline-bufferline plugin
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/mengelbrecht/lightline-bufferline.git %s/.vim_runtime/my_plugins/lightline-bufferline", dirs.HomeDir)); err != nil {
		return err
	}
	// vim-visual-multi
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/mg979/vim-visual-multi.git %s/.vim_runtime/my_plugins/vim-visual-multi", dirs.HomeDir)); err != nil {
		return err
	}
	// vim-helm
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/towolf/vim-helm.git %s/.vim_runtime/my_plugins/vim-helm", dirs.HomeDir)); err != nil {
		return err
	}
	// vim-airline plugin
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/vim-airline/vim-airline.git %s/.vim_runtime/my_plugins/vim-airline", dirs.HomeDir)); err != nil {
		return err
	}
	// vim-airline-themes plugin
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/vim-airline/vim-airline-themes.git %s/.vim_runtime/my_plugins/vim-airline-themes", dirs.HomeDir)); err != nil {
		return err
	}
	// indentLine vim plugin
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("git clone https://github.com/Yggdroot/indentLine.git %s/.vim_runtime/my_plugins/indentLine", dirs.HomeDir)); err != nil {
		return err
	}

	vimCustomConfig := fmt.Sprintf("templates/%s/my_configs.vim", osType)
	myConfigFS, err := myConfigs.Open(vimCustomConfig)
	if err != nil {
		return err
	}

	vimrcFS, err := myConfigs.Open("templates/.vimrc")
	if err != nil {
		return err
	}

	if err := localio.EmbedFileCopy("~/.vim_runtime/my_configs.vim", myConfigFS); err != nil {
		return err
	}
	if err := localio.EmbedFileCopy("~/.vimrc", vimrcFS); err != nil {
		return err
	}

	fmt.Println("[+] Installing vim binaries via +GoInstallBinaries")
	if err := localio.RunCommandPipeOutput("vim -es -c \"GoInstallBinaries\" -c \"qa\" || true"); err != nil {
		return err
	}

	return nil
}

// InstallVimAwesome ...
func InstallVimAwesome(osType string, dirs *localio.Directories, packages *localio.InstalledPackages) error {
	if exists, err := localio.Exists(fmt.Sprintf("%s/.vim_runtime", dirs.HomeDir)); err == nil && exists == true {
		return nil
	}

	switch osType {
	case "darwin":
		if !localio.CorrectOS("darwin") {
			break
		}
		// brew install macvim
		if err := localio.BrewInstallProgram("macvim", "vim", packages); err != nil {
			return err
		}
		if err := localio.BrewInstallProgram("cmake", "cmake", packages); err != nil {
			return err
		}
		if err := localio.BrewInstallProgram("pkg-config", "pkg-config", packages); err != nil {
			return err
		}
		if err := localio.BrewInstallProgram("python@3.9", "python3", packages); err != nil {
			return err
		}
		if err := localio.BrewInstallProgram("mono", "mono", packages); err != nil {
			return err
		}
		if err := localio.BrewInstallProgram("go", "go", packages); err != nil {
			return err
		}
		if err := localio.BrewInstallProgram("nodejs", "nodejs", packages); err != nil {
			return err
		}
	case "linux":
		if !localio.CorrectOS("linux") {
			break
		}
		// add mono to apt repos. Different for debian / ubuntu
		var packagesToInstall []string
		switch {
		case distro.IsUbuntu():
			packagesToInstall = []string{
				"gnupg",
				"ca-certificates"}
			if err := localio.AptInstall(packages, packagesToInstall...); err != nil {
				return err
			}
			if err := localio.RunCommandPipeOutput("sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys \"3FA7E0328081BFF6A14DA29AA6A19B38D3D831EF\" || true"); err != nil {
				return err
			}
			if err := localio.RunCommandPipeOutput("echo \"deb https://download.mono-project.com/repo/ubuntu stable-focal main\" | sudo tee /etc/apt/sources.list.d/mono-official-stable.list || true"); err != nil {
				return err
			}
			if err := localio.RunCommandPipeOutput("sudo apt-get update -y"); err != nil {
				return err
			}
			if err := localio.AptInstall(packages, "mono-complete"); err != nil {
				return err
			}
		case distro.IsDebian():
			packagesToInstall = []string{
				"gnupg",
				"ca-certificates",
				"apt-transport-https",
				"dirmngr"}
			if err := localio.AptInstall(packages, packagesToInstall...); err != nil {
				return err
			}
			if err := localio.RunCommandPipeOutput("sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys \"3FA7E0328081BFF6A14DA29AA6A19B38D3D831EF\" || true"); err != nil {
				return err
			}
			if err := localio.RunCommandPipeOutput("echo \"deb https://download.mono-project.com/repo/debian stable-buster main\" | sudo tee /etc/apt/sources.list.d/mono-official-stable.list || true"); err != nil {
				return err
			}
			if err := localio.RunCommandPipeOutput("sudo apt-get update -y"); err != nil {
				return err
			}
			if err := localio.AptInstall(packages, "mono-complete"); err != nil {
				return err
			}
		default:
			packagesToInstall = []string{
				"gnupg",
				"ca-certificates",
				"mono-complete"}
			if err := localio.AptInstall(packages, packagesToInstall...); err != nil {
				return err
			}
		}
		universalPackagesToInstall := []string{
			"build-essential",
			"cmake",
			"vim-nox",
			"python3-dev",
			"nodejs",
			"default-jdk",
			"npm",
			"jq"}

		if err := localio.AptInstall(packages, universalPackagesToInstall...); err != nil {
			return err
		}
	}
	if err := InstallVimPlugins(osType, dirs); err != nil {
		return err
	}

	return nil
}
