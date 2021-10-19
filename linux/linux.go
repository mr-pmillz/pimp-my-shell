package linux

import (
	"embed"
	"fmt"
	"pimp-my-shell/localio"
)

//go:embed templates/*
var configFiles embed.FS

// CustomTilixBookmarks copies the upgrade your pty commands for quick usage when upgrading a netcat shell
func CustomTilixBookmarks() error {
	if !localio.CorrectOS("linux") {
		return nil
	}
	if exists, err := localio.Exists("~/.config/tilix/bookmarks.json"); err == nil && !exists {
		fmt.Println("[+] Setting up Custom Tilix Bookmarks")
		bookmarksJSON, err := configFiles.Open("templates/bookmarks.json")
		if err != nil {
			return err
		}
		if err = localio.EmbedFileCopy("~/.config/tilix/bookmarks.json", bookmarksJSON); err != nil {
			return err
		}
	}

	return nil
}
