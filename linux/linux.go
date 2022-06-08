package linux

import (
	"embed"
	"fmt"
	"os"

	"github.com/mr-pmillz/pimp-my-shell/localio"
)

//go:embed templates/*
var configFiles embed.FS

// CustomTerminalBookmarks copies the upgrade your pty commands for quick usage when upgrading a netcat shell
func CustomTerminalBookmarks() error {
	if !localio.CorrectOS("linux") {
		return nil
	}
	desktopSession := os.Getenv("DESKTOP_SESSION")
	switch desktopSession {
	case "lightdm-xsession":
		if exists, err := localio.Exists("~/.config/qterminal.org/qterminal_bookmarks.xml"); err == nil && !exists {
			fmt.Println("[+] Setting up Custom Terminal Bookmarks")
			bookmarksJSON, err := configFiles.Open("templates/bookmarks.xml")
			if err != nil {
				return err
			}
			defer bookmarksJSON.Close()
			if err = localio.EmbedFileCopy("~/.config/qterminal.org/qterminal_bookmarks.xml", bookmarksJSON); err != nil {
				return err
			}
		}
	case "gnome":
		// Do nothing for now
	}

	return nil
}
