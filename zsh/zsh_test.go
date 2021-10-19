package zsh

import (
	"pimp-my-shell/localio"
	"testing"
)

func TestInstallOhMyZsh(t *testing.T) {
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
