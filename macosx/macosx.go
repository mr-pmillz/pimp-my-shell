package macosx

import (
	"fmt"
	"pimp-my-shell/localio"
)

// InstallHomebrew if Not already installed
func InstallHomebrew(dirs *localio.Directories) error {
	if _, exists := localio.CommandExists("brew"); exists {
		return nil
	}
	installString := "/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
	err := localio.RunCommandPipeOutput(fmt.Sprintf("cd %s && %s", dirs.HomeDir, installString))
	return err
}
