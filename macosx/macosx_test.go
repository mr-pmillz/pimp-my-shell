package macosx

import (
	"pimp-my-shell/localio"
	"testing"
)

func TestInstallHomebrew(t *testing.T) {
	dirs, err := localio.NewDirectories()
	if err != nil {
		t.Errorf("failed to create Directories type: %v", err)
	}
	fakeDirs, err := localio.NewDirectories()
	if err != nil {
		t.Errorf("failed to create Directories type: %v", err)
	}
	fakeDirs.HomeDir = "/asdfasdf/sadfasdf/asdfsadf"
	type args struct {
		dirs *localio.Directories
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test InstallHomebrew 1", args{dirs: dirs}, false},
		{"Test InstallHomebrew 2 Should Fail", args{dirs: fakeDirs}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InstallHomebrew(tt.args.dirs); (err != nil) != tt.wantErr {
				t.Errorf("InstallHomebrew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
