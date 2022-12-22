package macosx

import (
	"fmt"

	"github.com/mr-pmillz/pimp-my-shell/v2/localio"
)

// InstallHomebrew if Not already installed
func InstallHomebrew(dirs *localio.Directories) error {
	if !localio.CorrectOS("darwin") {
		return nil
	}
	if _, exists := localio.CommandExists("brew"); exists {
		return nil
	}
	installString := "/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
	err := localio.RunCommandPipeOutput(fmt.Sprintf("cd %s && %s", dirs.HomeDir, installString))
	return err
}
