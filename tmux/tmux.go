package tmux

import (
	"embed"
	"fmt"
	"os"

	gotmux "github.com/jubnzv/go-tmux"
	"github.com/mr-pmillz/pimp-my-shell/localio"
)

//go:embed templates/*
var tmuxConfig embed.FS

// runTmux ...
func runTmux() error {
	server := new(gotmux.Server)

	// Check if "PimpMyShell" session already exists.
	exists, err := server.HasSession("PimpMyShell")
	if err != nil {
		msg := fmt.Errorf("Can't check 'PimpMyShell' session: %s", err)
		fmt.Println(msg)
		return err
	}
	if exists {
		fmt.Println("Session 'PimpMyShell' already exists!")
		return err
	}

	// Prepare configuration for a new session with some windows.
	session := gotmux.Session{Name: "PimpMyShell"}
	w1 := gotmux.Window{Name: "Human", Id: 0}
	w2 := gotmux.Window{Name: "Element", Id: 1}
	session.AddWindow(w1)
	session.AddWindow(w2)
	server.AddSession(session)
	var sessions []*gotmux.Session
	sessions = append(sessions, &session)
	conf := gotmux.Configuration{
		Server:        server,
		Sessions:      sessions,
		ActiveSession: nil}

	// Setup this configuration.
	err = conf.Apply()
	if err != nil {
		msg := fmt.Errorf("Can't apply prepared configuration: %s", err)
		fmt.Println(msg)
		return err
	}

	// Attach to created session
	err = session.AttachSession()
	if err != nil {
		msg := fmt.Errorf("Can't attached to created session: %s", err)
		fmt.Println(msg)
		return err
	}
	return nil
}

// StartTMUX is a helper func
func StartTMUX() error {
	if val, _ := os.LookupEnv("GITHUB_ACTIONS"); val != "true" {
		if err := localio.StartTmuxSession(); err != nil {
			return err
		}
		if err := runTmux(); err != nil {
			return err
		}
	}
	return nil
}

// InstallOhMyTmux takes 1 of 3 possible strings, linux darwin windows
func InstallOhMyTmux(osType string, dirs *localio.Directories, packages *localio.InstalledPackages) error {
	switch osType {
	case "darwin":
		if !localio.CorrectOS("darwin") {
			break
		}
		if err := localio.BrewInstallProgram("tmux", "tmux", packages); err != nil {
			return err
		}
		if err := localio.BrewInstallProgram("iproute2mac", "ip", packages); err != nil {
			return err
		}
		if err := localio.BrewInstallProgram("jq", "jq", packages); err != nil {
			return err
		}
		if err := localio.BrewInstallProgram("reattach-to-user-namespace", "reattach-to-user-namespace", packages); err != nil {
			return err
		}
	case "linux":
		if !localio.CorrectOS("linux") {
			break
		}
		if err := localio.AptInstall(packages, "tmux", "xclip"); err != nil {
			return err
		}
	}

	if exists, err := localio.Exists(fmt.Sprintf("%s/.tmux.conf.local", dirs.HomeDir)); err == nil && exists == true {
		return nil
	}

	cmdSlice := []string{
		fmt.Sprintf("cd %s && git clone https://github.com/gpakosz/.tmux.git", dirs.HomeDir),
		fmt.Sprintf("cd %s && ln -s -f .tmux/.tmux.conf", dirs.HomeDir),
	}

	if err := localio.RunCommands(cmdSlice); err != nil {
		return err
	}

	tmuxConfPath := fmt.Sprintf("templates/%s/tmux.conf.local", osType)
	tmuxConfFS, err := tmuxConfig.Open(tmuxConfPath)
	if err != nil {
		return err
	}

	if err := localio.EmbedFileCopy("~/.tmux.conf.local", tmuxConfFS); err != nil {
		return err
	}

	if err = localio.RunCommandPipeOutput(fmt.Sprintf("cd %s && git clone https://github.com/tmux-plugins/tpm .tmux/plugins/tpm", dirs.HomeDir)); err != nil {
		return err
	}

	return nil
}
