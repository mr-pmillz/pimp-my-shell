package localio

import (
	"fmt"
	"os"
	"testing"
)

func TestResolveAbsPath(t *testing.T) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("cant get homedir: %v", err)
	}
	type args struct {
		path string
	}
	var tests = []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Tilda Test", args{path: "~/.bash_history"}, fmt.Sprintf("%s/.bash_history", homedir), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveAbsPath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveAbsPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResolveAbsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	type args struct {
		path string
	}
	var tests = []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"Dir test", args{path: "/etc"}, true, false},
		{"File test", args{path: "/bin/bash"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}
