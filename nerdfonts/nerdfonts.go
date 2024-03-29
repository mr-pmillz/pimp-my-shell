package nerdfonts

import (
	"fmt"
	"os"

	"github.com/mr-pmillz/pimp-my-shell/v2/localio"
)

// InstallNerdFontsLSD ...
//
//nolint:gocognit
func InstallNerdFontsLSD(osType string, dirs *localio.Directories, packages *localio.InstalledPackages) error {
	switch osType {
	case "darwin":
		if !localio.CorrectOS("darwin") {
			break
		}
		// brew tap homebrew/cask-fonts
		if err := localio.BrewTap("homebrew/cask-fonts", packages); err != nil {
			return err
		}
		// brew install --cask font-meslo-lg-nerd-font
		if err := localio.BrewInstallCaskProgram("font-meslo-lg-nerd-font", "font-meslo-lg-nerd-font", packages); err != nil {
			return err
		}
	case "linux":
		if !localio.CorrectOS("linux") {
			break
		}
		// install meslo nerd fonts
		fontsDir := fmt.Sprintf("%s/.local/share/fonts", dirs.HomeDir)
		if exists, err := localio.Exists(fmt.Sprintf("%s/%s", fontsDir, "MesloLGS NF Regular.ttf")); err == nil && !exists {
			mesloLGSNFRegularURL := "https://github.com/romkatv/powerlevel10k-media/raw/master/MesloLGS%20NF%20Regular.ttf"
			mesloLGSNFBoldURL := "https://github.com/romkatv/powerlevel10k-media/raw/master/MesloLGS%20NF%20Bold.ttf"
			mesloLGSNFItalicURL := "https://github.com/romkatv/powerlevel10k-media/raw/master/MesloLGS%20NF%20Italic.ttf"
			mesloLGSNFBoldItalicURL := "https://github.com/romkatv/powerlevel10k-media/raw/master/MesloLGS%20NF%20Bold%20Italic.ttf"
			if err = os.MkdirAll(fontsDir, 0750); err != nil {
				return err
			}
			if err = localio.DownloadFile(fmt.Sprintf("%s/MesloLGS NF Regular.ttf", fontsDir), mesloLGSNFRegularURL); err != nil {
				return err
			}
			if err = localio.DownloadFile(fmt.Sprintf("%s/MesloLGS NF Bold.ttf", fontsDir), mesloLGSNFBoldURL); err != nil {
				return err
			}
			if err = localio.DownloadFile(fmt.Sprintf("%s/MesloLGS NF Italic.ttf", fontsDir), mesloLGSNFItalicURL); err != nil {
				return err
			}
			if err = localio.DownloadFile(fmt.Sprintf("%s/MesloLGS NF Bold Italic.ttf", fontsDir), mesloLGSNFBoldItalicURL); err != nil {
				return err
			}
			if err = localio.RunCommandPipeOutput("fc-cache -f -v"); err != nil {
				return err
			}
		}
	}

	return nil
}
