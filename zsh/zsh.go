package zsh

import (
	"embed"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/mr-pmillz/pimp-my-shell/v2/localio"
	"github.com/mr-pmillz/pimp-my-shell/v2/osrelease"
)

//go:embed templates/*
var zshConfigs embed.FS

func updateZSHPlugins(zshrcPath string) error {
	var re = regexp.MustCompile(`(?m)^plugins=\(\s*((?:[a-z][a-z0-9_-]+\s*)+)\)$`)
	input, err := os.ReadFile(zshrcPath)
	if err != nil {
		return err
	}
	newPlugins := []string{"git", "zsh-syntax-highlighting", "tmux", "golang", "zsh-autosuggestions", "virtualenv", "ansible", "docker", "docker-compose", "terraform", "kubectl", "helm", "fzf", "fd"}
	currentlyInstalledPlugins := re.FindStringSubmatch(string(input))
	// fmt.Printf("Plugins: %+v\n", currentlyInstalledPlugins[1])
	// Fixes https://github.com/mr-pmillz/pimp-my-shell/issues/53
	var installedPlugins []string
	switch {
	case len(currentlyInstalledPlugins) == 0 || len(currentlyInstalledPlugins) == 1:
		installedPlugins = []string{}
	default:
		installedPlugins = strings.Fields(currentlyInstalledPlugins[1])
	}

	for _, plugin := range newPlugins {
		if !localio.Contains(installedPlugins, plugin) {
			installedPlugins = append(installedPlugins, plugin)
		}
	}
	sort.Strings(installedPlugins)

	// fmt.Printf("InstalledPlugins: %+v\n", installedPlugins)
	pluginsToAdd := strings.Join(installedPlugins, "\n\t")
	updatedZshrcFile := re.ReplaceAllString(string(input), fmt.Sprintf("plugins=(\n\t%s\n)", pluginsToAdd))

	if err = os.WriteFile(zshrcPath, []byte(updatedZshrcFile), 0644); err != nil { //nolint:gosec
		return err
	}

	fmt.Println("[+] Plugins ACTIVATED !!!")
	fmt.Printf("plugins=(\n\t%s\n)\n", pluginsToAdd)
	return nil
}

const zshExtraConfigPrependTemplate = "templates/zshrc_extra_prepend.zsh"

// InstallOhMyZsh ...
//
//nolint:gocognit
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
			if !localio.CorrectOS("linux") {
				break
			}
			if exists, err = localio.Exists(fmt.Sprintf("%s/.zshrc", dirs.HomeDir)); err == nil && exists {
				// Kali linux weird zshrc constraint
				if exists, err = localio.Exists("/etc/os-release"); err == nil && !exists {
					break
				}
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

	if err = localio.EmbedFileCopy("~/.oh-my-zsh/tools/upgrade.zsh", upgradeScript); err != nil {
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

		if err = localio.EmbedFileCopy("~/.p10k.zsh", p10kConfig); err != nil {
			return err
		}

		fmt.Println("[+] Setting Powerlevel10k Theme")
		if err = localio.SetVariableValue("ZSH_THEME", "powerlevel10k\\/powerlevel10k", osType, "~/.zshrc"); err != nil {
			return err
		}

		zshExtraPrependConfig, err := zshConfigs.ReadFile(zshExtraConfigPrependTemplate)
		if err != nil {
			return err
		}

		if err := localio.EmbedFileStringPrependToDest(zshExtraPrependConfig, "~/.zshrc"); err != nil {
			return err
		}
		zshExtraConfigTemplate := fmt.Sprintf("templates/%s/zshrc_extra.zsh", osType)
		zshExtraConfig, err := zshConfigs.ReadFile(zshExtraConfigTemplate)
		if err != nil {
			return err
		}
		if err := localio.EmbedFileStringAppendToDest(zshExtraConfig, "~/.zshrc"); err != nil {
			return err
		}
	}
	// install zsh-autosuggestions
	if exists, err := localio.Exists(fmt.Sprintf("%s/.oh-my-zsh/custom/plugins/zsh-autosuggestions", dirs.HomeDir)); err == nil && !exists {
		if err = localio.GitClone("https://github.com/zsh-users/zsh-autosuggestions", fmt.Sprintf("%s/.oh-my-zsh/custom/plugins/zsh-autosuggestions", dirs.HomeDir)); err != nil {
			return err
		}
	}

	// install zsh-syntax-highlighting
	if exists, err := localio.Exists(fmt.Sprintf("%s/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting", dirs.HomeDir)); err == nil && !exists {
		if err = localio.GitClone("https://github.com/zsh-users/zsh-syntax-highlighting.git", fmt.Sprintf("%s/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting", dirs.HomeDir)); err != nil {
			return err
		}
	}

	return updateZSHPlugins(fmt.Sprintf("%s/.zshrc", dirs.HomeDir))
}
