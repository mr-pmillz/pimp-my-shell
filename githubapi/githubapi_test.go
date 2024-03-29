package githubapi

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/mr-pmillz/pimp-my-shell/v2/localio"
	"github.com/stretchr/testify/assert"
)

const (
	owner = "mr-pmillz"
	repo  = "pimp-my-shell"
)

func TestDownloadLatestRelease(t *testing.T) {
	dirs, err := localio.NewDirectories()
	if err != nil {
		log.Panic(err)
	}
	releaseAssets, err := getLatestReleasesFromGithubRepo(owner, repo)
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
				owner:  owner,
				repo:   repo,
			},
			fmt.Sprintf("%s/%s", dirs.HomeDir, releaseAssets.LinuxAMDFileName),
			false,
		},
		{
			"Test DownloadLatestRelease From This Repo darwin 2",
			args{
				osType: "darwin",
				dirs:   dirs,
				owner:  owner,
				repo:   repo,
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

func Test_getLatestReleasesFromGithubRepo(t *testing.T) {
	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name    string
		args    args
		want    *ReleaseAssets
		wantErr bool
	}{
		{name: "PimpMyShell Releases", args: args{
			owner: owner,
			repo:  repo,
		}, want: &ReleaseAssets{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getLatestReleasesFromGithubRepo(tt.args.owner, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLatestReleasesFromGithubRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ok := assert.IsType(t, &ReleaseAssets{}, got); !ok {
				t.Errorf("getLatestReleasesFromGithubRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}
