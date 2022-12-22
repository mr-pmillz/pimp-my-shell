package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/mr-pmillz/pimp-my-shell/v2/localio"
)

func Test_pimpMyShell(t *testing.T) {
	dirs, err := localio.NewDirectories()
	if err != nil {
		t.Errorf("failed to create Directories type: %v", err)
	}
	type args struct {
		osType            string
		dirs              *localio.Directories
		installedPackages *localio.InstalledPackages
	}
	packages := &localio.InstalledPackages{}
	packages.BrewInstalledPackages = &localio.BrewInstalled{
		Names: []string{"cointop"}, CaskFullNames: []string{"cointop"}, Taps: []string{"homebrew/core"},
	}
	packages.AptInstalledPackages = &localio.AptInstalled{Name: []string{"sslscan"}}
	osType := runtime.GOOS
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: fmt.Sprintf("Test Test_pimpMyShell %s 1", osType), args: args{
			osType:            osType,
			dirs:              dirs,
			installedPackages: packages,
		}},
	}
	timeout := time.After(30 * time.Minute)
	done := make(chan bool)
	go func() {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := pimpMyShell(tt.args.osType, tt.args.dirs, tt.args.installedPackages); (err != nil) != tt.wantErr {
					t.Errorf("pimpMyShell() error = %v, wantErr %v", err, tt.wantErr)
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

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Test_main 1"},
	}
	timeout := time.After(30 * time.Minute)
	done := make(chan bool)
	go func() {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				main()
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
