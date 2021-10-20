package zsh

import (
	"fmt"
	"os"
	"pimp-my-shell/localio"
	"runtime"
	"testing"
)

// Test_updateZSHPlugins is a basic test checking for errors, it doesn't validate the correct plugins
// were added yet...
func Test_updateZSHPlugins(t *testing.T) {
	zshrcTestTemplatePlugins, err := localio.ResolveAbsPath("test/zshrc-test-template-plugins.zshrc")
	if err != nil {
		t.Errorf("Couldnt resolve zshrcTestTemplatePlugins path: %v", err)
	}
	zshrcTestTemplateSingleLinePlugin, err := localio.ResolveAbsPath("test/zshrc-test-single-line-plugin.zshrc")
	if err != nil {
		t.Errorf("Couldnt resolve zshrcTestTemplateSingleLinePlugin path: %v", err)
	}
	type args struct {
		zshrcPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test_updateZSHPlugins Multi-line 1", args: args{zshrcTestTemplatePlugins}, wantErr: false},
		{name: "Test_updateZSHPlugins Single-line 2", args: args{zshrcTestTemplateSingleLinePlugin}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = updateZSHPlugins(tt.args.zshrcPath); (err != nil) != tt.wantErr {
				t.Errorf("updateZSHPlugins() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if err = localio.RunCommandPipeOutput(fmt.Sprintf("cat %s %s", zshrcTestTemplatePlugins, zshrcTestTemplateSingleLinePlugin)); err != nil {
		t.Errorf("couldn't cat the files: %v", err)
	}
}

func TestInstallOhMyZsh(t *testing.T) {
	dirs, err := localio.NewDirectories()
	if err != nil {
		t.Errorf("failed to create Directories type: %v", err)
	}
	type args struct {
		osType string
		dirs   *localio.Directories
	}
	packages := &localio.InstalledPackages{}
	packages.BrewInstalledPackages = &localio.BrewInstalled{
		Names: []string{"bat"}, CaskFullNames: []string{"bat"}, Taps: []string{"bat"},
	}
	packages.AptInstalledPackages = &localio.AptInstalled{Name: []string{"bat"}}

	if val, _ := os.LookupEnv("GITHUB_ACTIONS"); val == "true" {
		if err = localio.RunCommandPipeOutput(fmt.Sprintf("rm -rf %s/.oh-my-zsh", dirs.HomeDir)); err != nil {
			t.Errorf("couldn't remove ~/.vim_runtime dir. error = %v", err)
		}
	}
	osType := runtime.GOOS
	switch osType {
	case "linux":
		if err := localio.AptInstall(packages, "zsh"); err != nil {
			t.Errorf("couldn't install zsh via apt-get: %v", err)
		}
	case "darwin":
		// ensure zsh is installed for unit tests
		if err := localio.BrewInstallProgram("zsh", "zsh", packages); err != nil {
			t.Errorf("couldn't install zsh via homebrew: %v", err)
		}
		// ensure gnu-sed is installed for unit tests
		if err := localio.BrewInstallProgram("gnu-sed", "gsed", packages); err != nil {
			t.Errorf("couldn't install gnu-sed via homebrew: %v", err)
		}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test InstallOhMyZsh darwin 1", args{
			osType: "darwin",
			dirs:   dirs,
		}, false},
		{"Test InstallOhMyZsh Linux 2", args{
			osType: "linux",
			dirs:   dirs,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InstallOhMyZsh(tt.args.osType, tt.args.dirs); (err != nil) != tt.wantErr {
				t.Errorf("InstallOhMyZsh() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
