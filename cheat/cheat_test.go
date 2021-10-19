package cheat

import (
	_ "embed"
	"pimp-my-shell/localio"
	"testing"
)

func TestInstallCheat(t *testing.T) {
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
		{"Test InstallCheat darwin 1", args{
			osType:   "darwin",
			dirs:     dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages:  nil,
				BrewInstalledPackages: &localio.BrewInstalled{
					Names: []string{"bat"}, CaskFullNames: []string{"bat"}, Taps: []string{"homebrew/core"},
				},
		}}, false},
		{"Test InstallCheat Linux 2", args{
			osType:   "linux",
			dirs:     dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages:  &localio.AptInstalled{Name: []string{"bat"}},
				BrewInstalledPackages: nil,
			}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InstallCheat(tt.args.osType, tt.args.dirs, tt.args.packages); (err != nil) != tt.wantErr {
				t.Errorf("InstallCheat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
