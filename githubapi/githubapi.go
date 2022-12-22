package githubapi

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/mr-pmillz/pimp-my-shell/localio"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

const (
	darwin = "darwin"
	linux  = "linux"
	amd64  = "amd64"
	arm64  = "arm64"
	musl   = "musl"
)

// ReleaseAssets ...
type ReleaseAssets struct {
	LinuxARMURL       string
	LinuxARMFileName  string
	LinuxAMDURL       string
	LinuxAMDFileName  string
	DarwinARMURL      string
	DarwinARMFileName string
	DarwinAMDURL      string
	DarwinAMDFileName string
}

func getLatestReleasesFromGithubRepo(owner, repo string) (*ReleaseAssets, error) {
	var client = &github.Client{}
	ctx := context.Background()
	githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if ok {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	latestRelease, resp, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	// Rate.Limit should most likely be 5000 when authorized.
	log.Printf("Github API Rate: %#v\n", resp.Rate)

	// If a Token Expiration has been set, it will be displayed.
	if !resp.TokenExpiration.IsZero() {
		log.Printf("Github Token Expiration: %v\n", resp.TokenExpiration)
	}
	var releaseAssetsDownloadURLS []string
	r := ReleaseAssets{}

	for _, release := range latestRelease.Assets {
		releaseAssetsDownloadURLS = append(releaseAssetsDownloadURLS, *release.BrowserDownloadURL)
	}
	for _, releaseTypeURL := range releaseAssetsDownloadURLS {
		if strings.Contains(releaseTypeURL, amd64) && strings.HasSuffix(releaseTypeURL, ".deb") && !strings.Contains(releaseTypeURL, musl) {
			r.LinuxAMDURL = releaseTypeURL
			r.LinuxAMDFileName = path.Base(releaseTypeURL)
		}
		if strings.Contains(releaseTypeURL, amd64) && strings.HasSuffix(releaseTypeURL, ".gz") && !strings.Contains(releaseTypeURL, musl) && strings.Contains(releaseTypeURL, linux) {
			r.LinuxAMDURL = releaseTypeURL
			r.LinuxAMDFileName = path.Base(releaseTypeURL)
		}
		if strings.Contains(releaseTypeURL, arm64) && strings.HasSuffix(releaseTypeURL, ".deb") && !strings.Contains(releaseTypeURL, musl) {
			r.LinuxARMURL = releaseTypeURL
			r.LinuxARMFileName = path.Base(releaseTypeURL)
		}
		if strings.Contains(releaseTypeURL, arm64) && strings.HasSuffix(releaseTypeURL, ".gz") && !strings.Contains(releaseTypeURL, musl) && strings.Contains(releaseTypeURL, linux) {
			r.LinuxARMURL = releaseTypeURL
			r.LinuxARMFileName = path.Base(releaseTypeURL)
		}
		if strings.Contains(releaseTypeURL, "x86_64") && strings.HasSuffix(releaseTypeURL, ".gz") && strings.Contains(releaseTypeURL, darwin) && !strings.Contains(releaseTypeURL, arm64) {
			r.DarwinAMDURL = releaseTypeURL
			r.DarwinAMDFileName = path.Base(releaseTypeURL)
		}
		if strings.Contains(releaseTypeURL, arm64) && strings.HasSuffix(releaseTypeURL, ".gz") && strings.Contains(releaseTypeURL, darwin) && !strings.Contains(releaseTypeURL, amd64) {
			r.DarwinARMURL = releaseTypeURL
			r.DarwinARMFileName = path.Base(releaseTypeURL)
		}
	}
	return &r, nil
}

// DownloadLatestRelease ...
//
//nolint:gocognit
func DownloadLatestRelease(osType string, dirs *localio.Directories, owner, repo string) (string, error) {
	releaseAssets, err := getLatestReleasesFromGithubRepo(owner, repo)
	if err != nil {
		return "", err
	}

	cpuType := localio.GetCPUType()
	switch osType {
	case "darwin":
		switch cpuType {
		case "AMD64":
			if releaseAssets.DarwinAMDURL != "" {
				dest := fmt.Sprintf("%s/%s", dirs.HomeDir, releaseAssets.DarwinAMDFileName)
				if err = localio.DownloadFile(dest, releaseAssets.DarwinAMDURL); err != nil {
					return "", err
				}
				return dest, nil
			}
		case "ARM64":
			if releaseAssets.DarwinARMURL != "" {
				dest := fmt.Sprintf("%s/%s", dirs.HomeDir, releaseAssets.DarwinARMFileName)
				if err = localio.DownloadFile(dest, releaseAssets.DarwinARMURL); err != nil {
					return "", err
				}
				return dest, nil
			}
		default:
			fmt.Println("[-] Unsupported CPU")
		}
	case "linux":
		switch cpuType {
		case "AMD64":
			if releaseAssets.LinuxAMDURL != "" {
				dest := fmt.Sprintf("%s/%s", dirs.HomeDir, releaseAssets.LinuxAMDFileName)
				if err = localio.DownloadFile(dest, releaseAssets.LinuxAMDURL); err != nil {
					return "", err
				}
				return dest, nil
			}
		case "ARM64":
			if releaseAssets.LinuxARMURL != "" {
				dest := fmt.Sprintf("%s/%s", dirs.HomeDir, releaseAssets.LinuxARMFileName)
				if err = localio.DownloadFile(dest, releaseAssets.LinuxARMURL); err != nil {
					return "", err
				}
				return dest, nil
			}
		default:
			fmt.Println("[-] Unsupported CPU")
		}
	default:
		fmt.Println("[-] Unsupported OS or release doesn't exist")
	}

	return "", nil
}
