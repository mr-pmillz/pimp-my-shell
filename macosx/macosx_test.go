package macosx

import (
	"testing"

	"github.com/mr-pmillz/pimp-my-shell/v2/localio"
)

func TestInstallHomebrew(t *testing.T) {
	dirs, err := localio.NewDirectories()
	if err != nil {
		t.Errorf("failed to create Directories type: %v", err)
	}
	type args struct {
		dirs *localio.Directories
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test InstallHomebrew 1", args{dirs: dirs}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InstallHomebrew(tt.args.dirs); (err != nil) != tt.wantErr {
				t.Errorf("InstallHomebrew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
