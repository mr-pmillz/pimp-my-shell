package vim

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mr-pmillz/pimp-my-shell/localio"
)

func TestInstallVimPlugins(t *testing.T) {
	dirs, err := localio.NewDirectories()
	if err != nil {
		t.Errorf("failed to create Directories type: %v", err)
	}
	type args struct {
		osType string
		dirs   *localio.Directories
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test VimPlugins darwin 1", args{
			osType: "darwin",
			dirs:   dirs,
		}, false},
		{"Test VimPlugins Linux 2", args{
			osType: "linux",
			dirs:   dirs,
		}, false},
	}
	timeout := time.After(20 * time.Minute)
	done := make(chan bool)
	go func() {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := InstallVimPlugins(tt.args.osType, tt.args.dirs); (err != nil) != tt.wantErr {
					t.Errorf("InstallVimPlugins() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}
}

func TestInstallVimAwesome(t *testing.T) {
	dirs, err := localio.NewDirectories()
	if err != nil {
		t.Errorf("failed to create Directories type: %v", err)
	}
	type args struct {
		osType   string
		dirs     *localio.Directories
		packages *localio.InstalledPackages
	}
	if val, _ := os.LookupEnv("GITHUB_ACTIONS"); val == "true" {
		if err = localio.RunCommandPipeOutput(fmt.Sprintf("rm -rf %s/.vim_runtime", dirs.HomeDir)); err != nil {
			t.Errorf("couldn't remove ~/.vim_runtime dir. error = %v", err)
		}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Darwin 1", args{
			osType: "darwin",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages: nil,
				BrewInstalledPackages: &localio.BrewInstalled{
					Names: []string{"bat", "lsd", "gnu-sed", "gotop", "yamllint", "git-delta"}, CaskFullNames: []string{"bat"}, Taps: []string{"homebrew/core", "cjbassi/gotop"},
				},
			}}, false},
		{"Linux 1", args{
			osType: "linux",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages:  &localio.AptInstalled{Name: []string{"bat", "lsd", "gotop", "delta"}},
				BrewInstalledPackages: nil,
			}}, false},
		{"Linux 2", args{
			osType: "linux",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages:  nil,
				BrewInstalledPackages: nil,
			}}, false},
		{"Darwin 2", args{
			osType: "darwin",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages:  nil,
				BrewInstalledPackages: nil,
			}}, false},
	}
	timeout := time.After(20 * time.Minute)
	done := make(chan bool)
	go func() {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := InstallVimAwesome(tt.args.osType, tt.args.dirs, tt.args.packages); (err != nil) != tt.wantErr {
					t.Errorf("InstallVimAwesome() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
		done <- true
	}()
	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}
}
