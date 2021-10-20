package tmux

import (
	"pimp-my-shell/localio"
	"testing"
	"time"
)

func TestInstallOhMyTmux(t *testing.T) {
	dirs, err := localio.NewDirectories()
	if err != nil {
		t.Errorf("failed to create Directories type: %v", err)
	}
	type args struct {
		osType   string
		dirs     *localio.Directories
		packages *localio.InstalledPackages
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test InstallOhMyTmux darwin 1", args{
			osType: "darwin",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages: nil,
				BrewInstalledPackages: &localio.BrewInstalled{
					Names: []string{"bat", "lsd", "gnu-sed", "gotop", "yamllint", "git-delta"}, CaskFullNames: []string{"bat"}, Taps: []string{"homebrew/core", "cjbassi/gotop"},
				},
			}}, false},
		{"Test InstallOhMyTmux Linux 2", args{
			osType: "linux",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages:  &localio.AptInstalled{Name: []string{"xclip"}},
				BrewInstalledPackages: nil,
			}}, false},
	}
	timeout := time.After(20 * time.Minute)
	done := make(chan bool)
	go func() {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := InstallOhMyTmux(tt.args.osType, tt.args.dirs, tt.args.packages); (err != nil) != tt.wantErr {
					t.Errorf("InstallOhMyTmux() error = %v, wantErr %v", err, tt.wantErr)
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
