package zsh

import (
	"embed"
	"fmt"
	"github.com/Masterminds/semver"
	"io/ioutil"
	"os"
	"pimp-my-shell/localio"
	"pimp-my-shell/osrelease"
	"regexp"
	"strings"
)

//go:embed templates/*
var zshConfigs embed.FS

func updateZSHPlugins(dirs *localio.Directories) error {
	var re = regexp.MustCompile(`(?m)^plugins=\(\s*((?:[a-z][a-z0-9_-]+\s*)+)\)$`)
	input, err := ioutil.ReadFile(fmt.Sprintf("%s/.zshrc", dirs.HomeDir))
	if err != nil {
		return err
	}
	newPlugins := []string{"git", "zsh-syntax-highlighting", "tmux", "zsh-autosuggestions", "virtualenv", "ansible", "docker", "docker-compose", "terraform", "kubectl", "helm", "fzf"}
	currentlyInstalledPlugins := re.FindStringSubmatch(string(input))
	//fmt.Printf("Plugins: %+v\n", currentlyInstalledPlugins[1])
	installedPlugins := strings.Fields(currentlyInstalledPlugins[1])

	for _, plugin := range newPlugins {
		if !localio.Contains(installedPlugins, plugin) {
			installedPlugins = append(installedPlugins, plugin)
		}
	}

	//fmt.Printf("InstalledPlugins: %+v\n", installedPlugins)
	pluginsToAdd := strings.Join(installedPlugins, "\n\t")
	updatedZshrcFile := re.ReplaceAllString(string(input), fmt.Sprintf("plugins=(\n\t%s\n)", pluginsToAdd))

	err = ioutil.WriteFile(fmt.Sprintf("%s/.zshrc", dirs.HomeDir), []byte(updatedZshrcFile), 0644)
	if err != nil {
		return err
	}

	fmt.Println("[+] Plugins ACTIVATED !!!")
	fmt.Printf("plugins=(\n\t%s\n)\n", pluginsToAdd)
	return nil
}

// InstallOhMyZsh ...
func InstallOhMyZsh(osType string, dirs *localio.Directories) error {
	// install ohmyzsh if not already installed
	if exists, err := localio.Exists(fmt.Sprintf("%s/.oh-my-zsh", dirs.HomeDir)); err == nil && !exists {
		ohMyZshInstallScriptURL := "https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh"
		dest := fmt.Sprintf("%s/install_ohmyzsh.sh", dirs.HomeDir)
		if err = localio.DownloadFile(dest, ohMyZshInstallScriptURL); err != nil {
			return err
		}

		switch osType {
		case "linux":
			if exists, err = localio.Exists(fmt.Sprintf("%s/.zshrc", dirs.HomeDir)); err == nil && exists {
				// Kali linux weird zshrc constraint
				osINFO, err := osrelease.Read()
				if err != nil {
					return err
				}
				if osINFO["ID"] == "kali" {
					kaliZshConstraint, err := semver.NewConstraint(">= 2020.4")
					if err != nil {
						return err
					}
					currentOSReleaseID, err := semver.NewVersion(osINFO["VERSION"])
					if err != nil {
						return err
					}
					isKaliLaterThan20204 := kaliZshConstraint.Check(currentOSReleaseID)
					if isKaliLaterThan20204 {
						fmt.Println("Your Kali version >= 2020.4 has highly custom .zshrc. Moving to ~/.zshrc_pre_pimpmyshell.bak")
						if err = os.Rename(fmt.Sprintf("%s/.zshrc", dirs.HomeDir), fmt.Sprintf("%s/.zshrc_pre_pimpmyshell.bak", dirs.HomeDir)); err != nil {
							return err
						}
					}
				}
			}
		default:
			// Do Nothing
		}

		if err = localio.RunCommandPipeOutput(fmt.Sprintf("cd %s && sh %s --keep-zshrc --unattended || true", dirs.HomeDir, dest)); err != nil {
			return err
		}
		if err = os.Remove(dest); err != nil {
			return err
		}
	}

	if exists, err := localio.Exists("~/.oh-my-zsh/custom/aliases.zsh"); err == nil && !exists {
		zshCustomAliasesTemplate := fmt.Sprintf("templates/%s/aliases.zsh", osType)
		aliasesConfig, err := zshConfigs.Open(zshCustomAliasesTemplate)
		if err != nil {
			return err
		}

		if err := localio.EmbedFileCopy("~/.oh-my-zsh/custom/aliases.zsh", aliasesConfig); err != nil {
			return err
		}

	}

	if exists, err := localio.Exists("~/.oh-my-zsh/custom/man-pages.zsh"); err == nil && !exists {
		manPagesConfig, err := zshConfigs.Open("templates/man-pages.zsh")
		if err != nil {
			return err
		}

		if err := localio.EmbedFileCopy("~/.oh-my-zsh/custom/man-pages.zsh", manPagesConfig); err != nil {
			return err
		}
	}

	upgradeScript, err := zshConfigs.Open("templates/upgrade.sh")
	if err != nil {
		return err
	}

	if err := localio.EmbedFileCopy("~/.oh-my-zsh/tools/upgrade.zsh", upgradeScript); err != nil {
		return err
	}

	// install powerlevel 10k
	if exists, err := localio.Exists("~/.oh-my-zsh/custom/themes/powerlevel10k"); err == nil && !exists {
		installString := fmt.Sprintf("git clone --depth=1 https://github.com/romkatv/powerlevel10k.git %s/.oh-my-zsh/custom/themes/powerlevel10k", dirs.HomeDir)
		if err = localio.RunCommandPipeOutput(fmt.Sprintf("cd %s && %s", dirs.HomeDir, installString)); err != nil {
			return err
		}

		p10kConfig, err := zshConfigs.Open("templates/.p10k.zsh")
		if err != nil {
			return err
		}

		if err := localio.EmbedFileCopy("~/.p10k.zsh", p10kConfig); err != nil {
			return err
		}

		fmt.Println("[+] Setting Powerlevel10k Theme")
		if err := localio.SetVariableValue("ZSH_THEME", "powerlevel10k\\/powerlevel10k", osType, "~/.zshrc"); err != nil {
			return err
		}

		zshExtraConfig, err := zshConfigs.ReadFile("templates/.zshrc_extra.zsh")
		if err != nil {
			return err
		}
		if err := localio.EmbedFileStringAppendToDest(zshExtraConfig, "~/.zshrc"); err != nil {
			return err
		}

	}

	// install zsh-autosuggestions
	if exists, err := localio.Exists(fmt.Sprintf("%s/.oh-my-zsh/custom/plugins/zsh-autosuggestions", dirs.HomeDir)); err == nil && !exists {
		installString := fmt.Sprintf("git clone https://github.com/zsh-users/zsh-autosuggestions %s/.oh-my-zsh/custom/plugins/zsh-autosuggestions", dirs.HomeDir)
		if err = localio.RunCommandPipeOutput(fmt.Sprintf("cd %s && %s", dirs.HomeDir, installString)); err != nil {
			return err
		}
	}

	// install zsh-syntax-highlighting
	if exists, err := localio.Exists(fmt.Sprintf("%s/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting", dirs.HomeDir)); err == nil && !exists {
		installString := fmt.Sprintf("git clone https://github.com/zsh-users/zsh-syntax-highlighting.git %s/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting", dirs.HomeDir)
		if err = localio.RunCommandPipeOutput(fmt.Sprintf("cd %s && %s", dirs.HomeDir, installString)); err != nil {
			return err
		}
	}

	if err = updateZSHPlugins(dirs); err != nil {
		return err
	}

	return nil
}
