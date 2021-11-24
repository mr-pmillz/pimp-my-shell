package cheat

import (
	"bytes"
	_ "embed" // single file embed
	"fmt"
	"os"
	"text/template"

	"github.com/mr-pmillz/pimp-my-shell/localio"
)

type cheatConfigOptions struct {
	CommunityPath string
	PersonalPath  string
}

//go:embed templates/conf.yml.tmpl
var cheatConfigTemplate string

func generateCheatConfig(dirs *localio.Directories) (string, error) {
	cheatConfig, err := template.New("cheatConfig").Parse(cheatConfigTemplate)
	if err != nil {
		return "", err
	}
	var cheatConfigBuf bytes.Buffer
	err = cheatConfig.Execute(&cheatConfigBuf, cheatConfigOptions{
		CommunityPath: fmt.Sprintf("%s/.config/cheat/cheatsheets/community", dirs.HomeDir),
		PersonalPath:  fmt.Sprintf("%s/.config/cheat/cheatsheets/personal", dirs.HomeDir),
	})
	if err != nil {
		return "", err
	}
	return cheatConfigBuf.String(), nil
}

// InstallCheat installs cheat command and clones the community cheat sheets and sets your cheat config paths
func InstallCheat(osType string, dirs *localio.Directories, packages *localio.InstalledPackages) error {
	switch osType {
	case "darwin":
		if !localio.CorrectOS("darwin") {
			break
		}
		if err := localio.BrewInstallProgram("cheat", "cheat", packages); err != nil {
			return err
		}
	case "linux":
		if !localio.CorrectOS("linux") {
			break
		}
		if _, exists := localio.CommandExists("cheat"); !exists {
			fmt.Println("[+] Installing cheat")
			if _, exists = localio.CommandExists("go"); exists {
				if err := localio.RunCommandPipeOutput("go get -u github.com/cheat/cheat/cmd/cheat"); err != nil {
					return err
				}
			}
		}
	}
	if exists, err := localio.Exists(fmt.Sprintf("%s/.config/cheat", dirs.HomeDir)); err == nil && !exists {
		if err = os.MkdirAll(fmt.Sprintf("%s/.config/cheat/cheatsheets/personal", dirs.HomeDir), 0750); err != nil {
			return err
		}
		if err = localio.GitClone("https://github.com/cheat/cheatsheets.git", fmt.Sprintf("%s/.config/cheat/cheatsheets/community", dirs.HomeDir)); err != nil {
			return err
		}
		generatedCheatConfig, err := generateCheatConfig(dirs)
		if err != nil {
			return err
		}
		if err = localio.CopyStringToFile(generatedCheatConfig, fmt.Sprintf("%s/.config/cheat/conf.yml", dirs.HomeDir)); err != nil {
			return err
		}
	}
	return nil
}
