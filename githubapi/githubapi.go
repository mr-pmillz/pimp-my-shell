package githubapi

import (
	"context"
	"fmt"
	"path"
	"pimp-my-shell/localio"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/klauspost/cpuid/v2"
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
	client := github.NewClient(nil)
	ctx := context.Background()
	latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	var releaseAssetsDownloadURLS []string
	r := ReleaseAssets{}

	for _, release := range latestRelease.Assets {
		releaseAssetsDownloadURLS = append(releaseAssetsDownloadURLS, *release.BrowserDownloadURL)
	}
	for _, releaseTypeURL := range releaseAssetsDownloadURLS {
		if strings.Contains(releaseTypeURL, "amd64") && strings.HasSuffix(releaseTypeURL, ".deb") && !strings.Contains(releaseTypeURL, "musl") {
			r.LinuxAMDURL = releaseTypeURL
			r.LinuxAMDFileName = path.Base(releaseTypeURL)
		}
		if strings.Contains(releaseTypeURL, "arm64") && strings.HasSuffix(releaseTypeURL, ".deb") && !strings.Contains(releaseTypeURL, "musl") {
			r.LinuxARMURL = releaseTypeURL
			r.LinuxARMFileName = path.Base(releaseTypeURL)
		}
		if strings.Contains(releaseTypeURL, "x86_64") && strings.HasSuffix(releaseTypeURL, ".gz") && strings.Contains(releaseTypeURL, "darwin") {
			r.DarwinAMDURL = releaseTypeURL
			r.DarwinAMDFileName = path.Base(releaseTypeURL)
		}
		if strings.Contains(releaseTypeURL, "arm64") && strings.HasSuffix(releaseTypeURL, ".gz") && strings.Contains(releaseTypeURL, "darwin") {
			r.DarwinARMURL = releaseTypeURL
			r.DarwinARMFileName = path.Base(releaseTypeURL)
		}
	}
	return &r, nil
}

// DownloadLatestRelease ...
func DownloadLatestRelease(osType string, dirs *localio.Directories, owner, repo string) (string, error) {
	releaseAssets, err := getLatestReleasesFromGithubRepo(owner, repo)
	if err != nil {
		return "", err
	}

	var cpuType string
	cpuid.Detect()
	if cpuid.CPU.VendorID.String() == "AMD" || cpuid.CPU.VendorID.String() == "Intel" && cpuid.CPU.CacheLine == 64 {
		cpuType = "AMD64"
	} else if cpuid.CPU.VendorID.String() == "ARM" && cpuid.CPU.CacheLine == 64 {
		cpuType = "ARM64"
	} else {
		cpuType = ""
	}

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
