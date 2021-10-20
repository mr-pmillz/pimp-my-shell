package extra

import (
	"embed"
	"fmt"
	"gopkg.in/ini.v1"
	"pimp-my-shell/githubapi"
	"pimp-my-shell/localio"
)

//go:embed templates/*
var extraConfigs embed.FS

// InstallExtraPackages ...
func InstallExtraPackages(osType string, dirs *localio.Directories, packages *localio.InstalledPackages) error {
	switch osType {
	case "darwin":
		if !localio.CorrectOS("darwin") {
			break
		}
		// install lsd
		if err := localio.BrewInstallProgram("lsd", "lsd", packages); err != nil {
			return err
		}
		// install cowsay because it's Pimpin'
		if err := localio.BrewInstallProgram("cowsay", "cowsay", packages); err != nil {
			return err
		}
		// install gnu-sed because mac BSD sed doesn't work very good
		if err := localio.BrewInstallProgram("gnu-sed", "gsed", packages); err != nil {
			return err
		}
		// install gotop
		if err := localio.BrewInstallProgram("gotop", "gotop", packages); err != nil {
			return err
		}
		// install yamllint
		if err := localio.BrewInstallProgram("yamllint", "yamllint", packages); err != nil {
			return err
		}
		// install git-delta
		if err := localio.BrewInstallProgram("git-delta", "delta", packages); err != nil {
			return err
		}
		// install bat
		if err := localio.BrewInstallProgram("bat", "bat", packages); err != nil {
			return err
		}
		// install fzf
		if _, exists := localio.CommandExists("fzf"); !exists {
			if err := localio.BrewInstallProgram("fzf", "fzf", packages); err != nil {
				return err
			}
			//To install useful key bindings and fuzzy completion:
			// /usr/local/opt/fzf/install --all
			if exists, err := localio.Exists("/usr/local/opt/fzf/install"); err == nil && exists {
				if err = localio.RunCommandPipeOutput("/usr/local/opt/fzf/install --all"); err != nil {
					return err
				}
			}
		}
		// install pip for python3 https://bootstrap.pypa.io/get-pip.py and python requests module
		if _, exists := localio.CommandExists("python3"); !exists {
			if err := localio.BrewInstallProgram("python@3.9", "python3", packages); err != nil {
				return err
			}
		}
		// download get-pip
		if _, exists := localio.CommandExists("python3"); exists {
			if err := localio.DownloadFile(fmt.Sprintf("%s/get-pip.py", dirs.HomeDir), "https://bootstrap.pypa.io/get-pip.py"); err != nil {
				return err
			}
			if err := localio.RunCommandPipeOutput(fmt.Sprintf("cd %s && python3 get-pip.py || true", dirs.HomeDir)); err != nil {
				return err
			}
			if err := localio.RunCommandPipeOutput("python3 -m pip install requests --user || true"); err != nil {
				return err
			}
		}

	case "linux":
		if !localio.CorrectOS("linux") {
			break
		}
		// install lsd
		if _, exists := localio.CommandExists("lsd"); !exists {
			lsdDebPackage, err := githubapi.DownloadLatestRelease(osType, dirs, "Peltoche", "lsd")
			if err != nil {
				return err
			}
			if err = localio.RunCommandPipeOutput(fmt.Sprintf("sudo dpkg --no-pager -i %s", lsdDebPackage)); err != nil {
				return err
			}
		}
		// install gotop
		if _, exists := localio.CommandExists("go"); exists {
			if _, exists = localio.CommandExists("gotop"); !exists {
				if err := localio.RunCommandPipeOutput("go install github.com/xxxserxxx/gotop/v4/cmd/gotop@latest"); err != nil {
					return err
				}
			}
		}

		// install cowsay
		if err := localio.AptInstall(packages, "cowsay", "bat"); err != nil {
			return err
		}

		// add batcat to path as bat
		if err := localio.RunCommandPipeOutput(fmt.Sprintf("mkdir -p %s/.local/bin", dirs.HomeDir)); err != nil {
			return err
		}

		if exists, err := localio.Exists(fmt.Sprintf("%s/.local/bin/bat", dirs.HomeDir)); err == nil && !exists {
			if err = localio.RunCommandPipeOutput(fmt.Sprintf("ln -sf /usr/bin/batcat %s/.local/bin/bat", dirs.HomeDir)); err != nil {
				return err
			}
		}

		// install fzf configuration
		if exists, err := localio.Exists(fmt.Sprintf("%s/.fzf", dirs.HomeDir)); err == nil && !exists {
			installString := fmt.Sprintf("git clone --depth 1 https://github.com/junegunn/fzf.git %s/.fzf && %s/.fzf/install --all", dirs.HomeDir, dirs.HomeDir)
			if err = localio.RunCommandPipeOutput(installString); err != nil {
				return err
			}
		}

		// install git-delta from latest github release
		if _, exists := localio.CommandExists("delta"); !exists {
			debPackage, err := githubapi.DownloadLatestRelease(osType, dirs, "dandavison", "delta")
			if err != nil {
				return err
			}
			fmt.Println("[+] Installing git-delta latest release")
			if err = localio.RunCommandPipeOutput(fmt.Sprintf("sudo dpkg --no-pager -i %s", debPackage)); err != nil {
				return err
			}
		}
	}
	if err := updateGitConfig(); err != nil {
		return err
	}

	return nil
}

// updateGitConfig ...
func updateGitConfig() error {
	gitConfig, err := extraConfigs.ReadFile("templates/.gitconfig")
	if err != nil {
		return err
	}

	exists, err := localio.Exists("~/.gitconfig")
	if err == nil && !exists {
		if err = localio.EmbedFileStringAppendToDest(gitConfig, "~/.gitconfig"); err != nil {
			return err
		}
	} else if err == nil && exists {
		gitConfigPath, err := localio.ResolveAbsPath("~/.gitconfig")
		if err != nil {
			return err
		}

		opts := ini.LoadOptions{PreserveSurroundedQuote: true}
		embeddedConfig, err := ini.LoadSources(opts, gitConfig)
		if err != nil {
			return err
		}

		localGitConfig, err := ini.LoadSources(opts, gitConfigPath)
		if err != nil {
			return err
		}

		sections := embeddedConfig.Sections()
		for _, section := range sections {
			keys := section.Keys()
			for _, key := range keys {
				localGitConfig.Section(section.Name()).Key(key.Name()).SetValue(key.Value())
			}
		}

		if err := localGitConfig.SaveToIndent(gitConfigPath, "    "); err != nil {
			return err
		}
	}

	return nil
}
