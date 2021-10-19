package githubapi

import (
	"fmt"
	"log"
	"pimp-my-shell/localio"
	"reflect"
	"testing"
)

func TestDownloadLatestRelease(t *testing.T) {
	dirs, err := localio.NewDirectories()
	if err != nil {
		log.Panic(err)
	}
	releaseAssets, err := getLatestReleasesFromGithubRepo("mr-pmillz", "pimp-my-shell")
	if err != nil {
		t.Errorf("couldn't get release assets: %v\n", err)
	}
	type args struct {
		osType string
		dirs   *localio.Directories
		owner  string
		repo   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"Test DownloadLatestRelease From This Repo linux 1",
			args{
				osType: "linux",
				dirs:   dirs,
				owner:  "mr-pmillz",
				repo:   "pimp-my-shell",
			},
			fmt.Sprintf("%s/%s", dirs.HomeDir, releaseAssets.LinuxAMDFileName),
			false,
		},
		{
			"Test DownloadLatestRelease From This Repo darwin 2",
			args{
				osType: "darwin",
				dirs:   dirs,
				owner:  "mr-pmillz",
				repo:   "pimp-my-shell",
			},
			fmt.Sprintf("%s/%s", dirs.HomeDir, releaseAssets.DarwinAMDFileName),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DownloadLatestRelease(tt.args.osType, tt.args.dirs, tt.args.owner, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadLatestRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf("") {
				t.Errorf("DownloadLatestRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}
