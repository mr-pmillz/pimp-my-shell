package main

import (
	"fmt"
	"github.com/mr-pmillz/pimp-my-shell/cheat"
	"github.com/mr-pmillz/pimp-my-shell/extra"
	"github.com/mr-pmillz/pimp-my-shell/linux"
	"github.com/mr-pmillz/pimp-my-shell/localio"
	"github.com/mr-pmillz/pimp-my-shell/macosx"
	"github.com/mr-pmillz/pimp-my-shell/nerdfonts"
	"github.com/mr-pmillz/pimp-my-shell/tmux"
	"github.com/mr-pmillz/pimp-my-shell/vim"
	"github.com/mr-pmillz/pimp-my-shell/zsh"
	"log"
	"runtime"
)

// pimpMyShell runs all the installation setup tasks
func pimpMyShell(osType string, dirs *localio.Directories, installedPackages *localio.InstalledPackages) error {
	fmt.Println(`
 _______________
< Pimp My Shell >
 ---------------
          \
           \
            \          __---__
                    _-       /--______
               __--( /     \ )XXXXXXXXXXX\v.
             .-XXX(   O   O  )XXXXXXXXXXXXXXX-
            /XXX(       U     )        XXXXXXX\
          /XXXXX(              )--_  XXXXXXXXXXX\
         /XXXXX/ (      O     )   XXXXXX   \XXXXX\
         XXXXX/   /            XXXXXX   \__ \XXXXX
         XXXXXX__/          XXXXXX         \__---->
 ---___  XXX__/          XXXXXX      \__         /
   \-  --__/   ___/\  XXXXXX            /  ___--/=
    \-\    ___/    XXXXXX              '--- XXXXXX
       \-\/XXX\ XXXXXX                      /XXXXX
         \XXXXXXXXX   \                    /XXXXX/
          \XXXXXX      >                 _/XXXXX/
            \XXXXX--__/              __-- XXXX/
             -XXXXXXXX---------------  XXXXXX-
                \XXXXXXXXXXXXXXXXXXXXXXXXXX/
                  ""VXXXXXXXXXXXXXXXXXXV""`)
	switch osType {
	case "darwin":
		if err := localio.BrewInstallProgram("ca-certificates", "ca-certificates", installedPackages); err != nil {
			return err
		}
		if err := localio.BrewInstallProgram("zsh", "zsh", installedPackages); err != nil {
			return err
		}
		// install gnu-sed because mac BSD sed doesn't work very good
		if err := localio.BrewInstallProgram("gnu-sed", "gsed", installedPackages); err != nil {
			return err
		}
	case "linux":
		if err := localio.AptInstall(installedPackages, "zsh", "tilix", "apt-transport-https"); err != nil {
			return err
		}
		if err := linux.CustomTilixBookmarks(); err != nil {
			return err
		}
		// Install the latest version of golang
		if err := localio.DownloadAndInstallLatestVersionOfGolang(dirs.HomeDir); err != nil {
			return err
		}
	}

	// install oh-my-zsh if not already installed
	if err := zsh.InstallOhMyZsh(osType, dirs); err != nil {
		return err
	}
	// install oh-my-tmux and tpm
	if err := tmux.InstallOhMyTmux(osType, dirs, installedPackages); err != nil {
		return err
	}
	// install extra packages
	if err := extra.InstallExtraPackages(osType, dirs, installedPackages); err != nil {
		return err
	}
	// install awesome vim configuration
	if err := vim.InstallVimAwesome(osType, dirs, installedPackages); err != nil {
		return err
	}
	// install cheat
	if err := cheat.InstallCheat(osType, dirs, installedPackages); err != nil {
		return err
	}
	// install nerdfonts
	if err := nerdfonts.InstallNerdFontsLSD(osType, dirs, installedPackages); err != nil {
		return err
	}

	fmt.Println("To Customize Powerlevel10k Theme, Run")
	fmt.Println("p10k configure")
	if err := localio.RunCommandPipeOutput("cowsay -f eyes \"See you Space Cowboy\""); err != nil {
		return err
	}

	if err := tmux.StartTMUX(); err != nil {
		return err
	}

	return nil
}

func main() {
	dirs, err := localio.NewDirectories()
	if err != nil {
		log.Panic(err)
	}
	var packages = &localio.InstalledPackages{}
	os := runtime.GOOS
	switch os {
	case "windows":
		fmt.Println("[-] Doesn't work for Windows")
	case "darwin":
		fmt.Println("[+] Pimping your Mac Terminal")
		if err = macosx.InstallHomebrew(dirs); err != nil {
			log.Panic(err)
		}
		// generate the brew installed cli cache
		installedBrewPackages, err := localio.NewBrewInstalled()
		if err != nil {
			log.Panic(err)
		}
		packages.BrewInstalledPackages = installedBrewPackages
		if err = pimpMyShell(os, dirs, packages); err != nil {
			log.Panic(err)
		}
	case "linux":
		fmt.Println("[+] Pimping your Linux Terminal")
		// generate apt installed cli cache
		installedAptPackages, err := localio.NewAptInstalled()
		if err != nil {
			log.Panic(err)
		}
		packages.AptInstalledPackages = installedAptPackages
		// do linux pimp my shell work
		if err = pimpMyShell(os, dirs, packages); err != nil {
			log.Panic(err)
		}
	default:
		fmt.Printf("Unknown OS: %s .\n", os)
	}
}
