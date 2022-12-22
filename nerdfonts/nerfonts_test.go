package nerdfonts

import (
	"testing"

	"github.com/mr-pmillz/pimp-my-shell/v2/localio"
)

func TestInstallNerdFontsLSD(t *testing.T) {
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
		{"Test InstallNerdFontsLSD Darwin", args{
			osType: "darwin",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages: nil,
				BrewInstalledPackages: &localio.BrewInstalled{
					Names: []string{"bat", "lsd", "gnu-sed", "yamllint", "git-delta"}, CaskFullNames: []string{"bat"}, Taps: []string{"homebrew/core", "cjbassi/gotop"},
				},
			}}, false},
		{"Test InstallNerdFontsLSD Linux", args{
			osType: "linux",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages:  &localio.AptInstalled{Name: []string{"bat", "lsd", "delta"}},
				BrewInstalledPackages: nil,
			}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InstallNerdFontsLSD(tt.args.osType, tt.args.dirs, tt.args.packages); (err != nil) != tt.wantErr {
				t.Errorf("InstallNerdFontsLSD() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
