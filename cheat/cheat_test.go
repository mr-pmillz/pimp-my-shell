package cheat

import (
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/mr-pmillz/pimp-my-shell/localio"
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
	if val, _ := os.LookupEnv("GITHUB_ACTIONS"); val == "true" {
		if err = localio.RunCommandPipeOutput(fmt.Sprintf("rm -rf %s/.config/cheat 2>/dev/null", dirs.HomeDir)); err != nil {
			t.Errorf("couldn't remove ~/.config/cheat dir. error = %v", err)
		}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test InstallCheat darwin 1", args{
			osType: "darwin",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages: nil,
				BrewInstalledPackages: &localio.BrewInstalled{
					Names: []string{"aom"}, CaskFullNames: []string{"aom"}, Taps: []string{"homebrew/core"},
				},
			}}, false},
		{"Test InstallCheat Linux 2", args{
			osType: "linux",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages:  &localio.AptInstalled{Name: []string{"aom"}},
				BrewInstalledPackages: nil,
			}}, false},
		{"Test InstallCheat darwin 3", args{
			osType: "darwin",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages: nil,
				BrewInstalledPackages: &localio.BrewInstalled{
					Names: []string{"cheat"}, CaskFullNames: []string{"aom"}, Taps: []string{"homebrew/core"},
				},
			}}, false},
		{"Test InstallCheat Linux 4", args{
			osType: "linux",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages:  nil,
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
