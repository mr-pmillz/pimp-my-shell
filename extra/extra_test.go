package extra

import (
	"runtime"
	"testing"
	"time"

	"github.com/mr-pmillz/pimp-my-shell/localio"
)

func TestInstallExtraPackages(t *testing.T) {
	dirs, err := localio.NewDirectories()
	if err != nil {
		t.Errorf("failed to create Directories type: %v", err)
	}
	osType := runtime.GOOS
	switch osType {
	case "linux":
		if err = localio.DownloadAndInstallLatestVersionOfGolang(dirs.HomeDir); err != nil {
			t.Errorf("couldn't download and install golang: %v", err)
		}
		if err = localio.RunCommandPipeOutput("go version"); err != nil {
			t.Errorf("couldn't get go version: %v", err)
		}
	default:
		//DoNothing
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
		{"Test InstallExtraPackages darwin 1", args{
			osType: "darwin",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages: nil,
				BrewInstalledPackages: &localio.BrewInstalled{
					Names: []string{"ca-certificates"}, CaskFullNames: []string{"ca-certificates"}, Taps: []string{"homebrew/core", "cjbassi/gotop"},
				},
			}}, false},
		{"Test InstallExtraPackages darwin lots of packages already installed 2", args{
			osType: "darwin",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages: nil,
				BrewInstalledPackages: &localio.BrewInstalled{
					Names: []string{"aom", "apr", "apr-util", "argon2", "aspell", "assimp", "autoconf", "bdw-gc", "binwalk", "boost", "brotli", "c-ares", "ca-certificates",
						"cairo", "cheat", "cmake", "cointop", "coreutils", "cscope", "curl", "dbus", "deployer", "docbook", "docbook-xsl", "double-conversion", "exiftool"},
					CaskFullNames: []string{"font-meslo-lg-nerd-font", "wireshark"},
					Taps:          []string{"hashicorp/tap", "homebrew/core", "microsoft/mssql-release"},
				},
			}}, false},
		{"Test InstallExtraPackages Linux 3", args{
			osType: "linux",
			dirs:   dirs,
			packages: &localio.InstalledPackages{
				AptInstalledPackages:  &localio.AptInstalled{Name: []string{"lsd"}},
				BrewInstalledPackages: nil,
			}}, false},
	}
	timeout := time.After(20 * time.Minute)
	done := make(chan bool)
	go func() {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := InstallExtraPackages(tt.args.osType, tt.args.dirs, tt.args.packages); (err != nil) != tt.wantErr {
					t.Errorf("InstallExtraPackages() error = %v, wantErr %v", err, tt.wantErr)
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
