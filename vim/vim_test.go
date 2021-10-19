package vim

import (
	"pimp-my-shell/localio"
	"testing"
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InstallVimPlugins(tt.args.osType, tt.args.dirs); (err != nil) != tt.wantErr {
				t.Errorf("InstallVimPlugins() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
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
				AptInstalledPackages:  &localio.AptInstalled{Name: []string{"bat", "lsd", "gotop", "delta"}},
				BrewInstalledPackages: nil,
			}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InstallVimAwesome(tt.args.osType, tt.args.dirs, tt.args.packages); (err != nil) != tt.wantErr {
				t.Errorf("InstallVimAwesome() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
