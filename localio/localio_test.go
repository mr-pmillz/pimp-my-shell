package localio

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"reflect"
	"runtime"
	"testing"
)

//go:embed test/*
var testEmbedFiles embed.FS

func TestEmbedFileCopy(t *testing.T) {
	myConfigFS, err := testEmbedFiles.Open("test/zshrc-test-append-to-file.zshrc")
	if err != nil {
		t.Errorf("EmbedFileStringAppendToDest() error = %v", err)
	}
	type args struct {
		dst string
		src fs.File
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestEmbedFileCopy 1", args: args{
			dst: "test/zshrc-test-copy-to-file.zshrc",
			src: myConfigFS,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EmbedFileCopy(tt.args.dst, tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("EmbedFileCopy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	t.Cleanup(func() {
		if exists, err := Exists("test/zshrc-test-copy-to-file.zshrc"); err == nil && exists {
			if err = os.Remove("test/zshrc-test-copy-to-file.zshrc"); err != nil {
				t.Errorf("couldnt remove test file: %v", err)
			}
		}
	})
}

func TestEmbedFileStringAppendToDest(t *testing.T) {
	testConfig, err := testEmbedFiles.ReadFile("test/zshrc_test_extra.zsh")
	if err != nil {
		t.Errorf("EmbedFileStringAppendToDest() error = %v", err)
	}
	type args struct {
		data []byte
		dest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestEmbedFileStringAppendToDest 1", args: args{
			data: testConfig,
			dest: "test/zshrc-test-append-to-file.zshrc",
		}},
		{name: "TestEmbedFileStringAppendToDest 1", args: args{
			data: testConfig,
			dest: "test/copy-embed-file-test.zshrc",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EmbedFileStringAppendToDest(tt.args.data, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("EmbedFileStringAppendToDest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	t.Cleanup(func() {
		if exists, err := Exists("test/copy-embed-file-test.zshrc"); err == nil && exists {
			if err = os.Remove("test/copy-embed-file-test.zshrc"); err != nil {
				t.Errorf("couldnt remove test file: %v", err)
			}
		}
	})
}

func TestDownloadAndInstallLatestVersionOfGolang(t *testing.T) {
	dirs, err := NewDirectories()
	if err != nil {
		t.Errorf("failed to create Directories type: %v", err)
	}
	type args struct {
		homeDir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestDownloadAndInstallLatestVersionOfGolang 1", args: args{homeDir: dirs.HomeDir}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DownloadAndInstallLatestVersionOfGolang(tt.args.homeDir); (err != nil) != tt.wantErr {
				t.Errorf("DownloadAndInstallLatestVersionOfGolang() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetVariableValue(t *testing.T) {
	type args struct {
		varName    string
		val        string
		osType     string
		configPath string
	}
	packages := &InstalledPackages{}
	packages.BrewInstalledPackages = &BrewInstalled{
		Names: []string{"bat"}, CaskFullNames: []string{"bat"}, Taps: []string{"bat"},
	}
	packages.AptInstalledPackages = &AptInstalled{Name: []string{"bat"}}

	osType := runtime.GOOS
	switch osType {
	case "linux":
		if err := AptInstall(packages, "zsh", "apt-transport-https"); err != nil {
			t.Errorf("couldn't install zsh via apt-get: %v", err)
		}
	case "darwin":
		// ensure zsh is installed for unit tests
		if err := BrewInstallProgram("zsh", "zsh", packages); err != nil {
			t.Errorf("couldn't install zsh via homebrew: %v", err)
		}
		// ensure gnu-sed is installed for unit tests
		if err := BrewInstallProgram("gnu-sed", "gsed", packages); err != nil {
			t.Errorf("couldn't install gnu-sed via homebrew: %v", err)
		}
	}
	zshrcTestTemplate, err := ResolveAbsPath("test/zshrc-test-template.zshrc")
	if err != nil {
		t.Errorf("Couldnt resolve zshrcTestTemplate path: %v", err)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: fmt.Sprintf("Test SetVariableValue %s 1", osType), args: args{
			varName:    "ZSH_THEME",
			val:        "powerlevel10k\\/powerlevel10k",
			osType:     osType,
			configPath: zshrcTestTemplate,
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetVariableValue(tt.args.varName, tt.args.val, tt.args.osType, tt.args.configPath); (err != nil) != tt.wantErr {
				t.Errorf("SetVariableValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResolveAbsPath(t *testing.T) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("cant get homedir: %v", err)
	}
	type args struct {
		path string
	}
	tests := []struct {
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
	tests := []struct {
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

func TestBrewInstallProgram(t *testing.T) {
	type args struct {
		brewName   string
		binaryName string
		packages   *InstalledPackages
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test BrewInstallProgram 1", args{brewName: "fzf", binaryName: "fzf", packages: &InstalledPackages{
			AptInstalledPackages: nil,
			BrewInstalledPackages: &BrewInstalled{
				Names: []string{"bat"}, CaskFullNames: []string{"bat"}, Taps: []string{"bat"},
			},
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BrewInstallProgram(tt.args.brewName, tt.args.binaryName, tt.args.packages); (err != nil) != tt.wantErr {
				t.Errorf("BrewInstallProgram() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAptInstall(t *testing.T) {
	type args struct {
		packages *InstalledPackages
		aptName  []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test AptInstall 1", args{packages: &InstalledPackages{
			AptInstalledPackages:  &AptInstalled{Name: []string{"bat"}},
			BrewInstalledPackages: nil,
		}, aptName: []string{"xclip"},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AptInstall(tt.args.packages, tt.args.aptName...); (err != nil) != tt.wantErr {
				t.Errorf("AptInstall() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewBrewInstalled(t *testing.T) {
	tests := []struct {
		name    string
		want    *BrewInstalled
		wantErr bool
	}{
		{"Test NewBrewInstalled 1", &BrewInstalled{
			Names:         []string{"cowsay", "bat", "watch"},
			CaskFullNames: []string{"font-meslo-lg-nerd-font"},
			Taps:          []string{"homebrew/core", "homebrew/cask-fonts"},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBrewInstalled()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBrewInstalled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(&BrewInstalled{}) {
				t.Errorf("NewBrewInstalled() = %v, want %v", got, &BrewInstalled{})
			}
		})
	}
}

func TestNewAptInstalled(t *testing.T) {
	tests := []struct {
		name    string
		want    *AptInstalled
		wantErr bool
	}{
		{"Test NewAptInstalled 1", &AptInstalled{[]string{"python3-dev", "cowsay"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAptInstalled()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAptInstalled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(&AptInstalled{}) {
				t.Errorf("NewAptInstalled() = %v, want %v", got, &AptInstalled{})
			}
		})
	}
}
func TestBrewInstallCaskProgram(t *testing.T) {
	type args struct {
		brewName     string
		brewFullName string
		packages     *InstalledPackages
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test BrewInstallCaskProgram 1 that already exists", args{packages: &InstalledPackages{
			AptInstalledPackages: nil,
			BrewInstalledPackages: &BrewInstalled{
				Names: []string{"font-meslo-lg-nerd-font"}, CaskFullNames: []string{"font-meslo-lg-nerd-font"}, Taps: []string{"homebrew/core"},
			},
		},
			brewFullName: "font-meslo-lg-nerd-font",
			brewName:     "font-meslo-lg-nerd-font",
		}, false},
		{"Test BrewInstallCaskProgram 2 that doesn't already exist", args{packages: &InstalledPackages{
			AptInstalledPackages: nil,
			BrewInstalledPackages: &BrewInstalled{
				Names: []string{"bat"}, CaskFullNames: []string{"bat"}, Taps: []string{"homebrew/core"},
			},
		},
			brewFullName: "font-meslo-lg-nerd-font",
			brewName:     "font-meslo-lg-nerd-font",
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BrewTap("homebrew/cask-fonts", tt.args.packages); err != nil {
				t.Errorf("BrewTap() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := BrewInstallCaskProgram(tt.args.brewName, tt.args.brewFullName, tt.args.packages); (err != nil) != tt.wantErr {
				t.Errorf("BrewInstallCaskProgram() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBrewTap(t *testing.T) {
	type args struct {
		brewTap  string
		packages *InstalledPackages
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test BrewInstallCaskProgram 1 that doesn't already exist", args{packages: &InstalledPackages{
			AptInstalledPackages: nil,
			BrewInstalledPackages: &BrewInstalled{
				Names: []string{"bat"}, CaskFullNames: []string{"bat"}, Taps: []string{"homebrew/core"},
			},
		},
			brewTap: "homebrew/cask-fonts",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BrewTap(tt.args.brewTap, tt.args.packages); (err != nil) != tt.wantErr {
				t.Errorf("BrewTap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		s   []string
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Test Contains 1", args: args{
			s:   []string{"hello", "world"},
			str: "foobar",
		}, want: false},
		{name: "Test Contains 2", args: args{
			s:   []string{"hello", "world"},
			str: "world",
		}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.str); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandExists(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{"Test CommandExists 1 Nonexistent command", args{cmd: "asdfasdfsdfsadf"}, "", false},
		{"Test CommandExists 2 find command", args{cmd: "find"}, "/usr/bin/find", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CommandExists(tt.args.cmd)
			if got != tt.want {
				t.Errorf("CommandExists() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CommandExists() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRunCommandPipeOutput(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test RunCommandPipeOutput 1", args{command: "ls -la"}, false},
		{"Test RunCommandPipeOutput 2", args{command: "asdfsdfsadfsdf"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunCommandPipeOutput(tt.args.command); (err != nil) != tt.wantErr {
				t.Errorf("RunCommandPipeOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecCMD(t *testing.T) {
	username, err := user.Current()
	if err != nil {
		t.Errorf("can't get username: %v", err)
	}
	currentUsername := username.Username
	type args struct {
		command string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Test ExecCMD 1", args{command: "whoami"}, currentUsername, false},
		{"Test ExecCMD 2", args{command: "asdfsadfasdf"}, currentUsername, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecCMD(tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecCMD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf("") {
				t.Errorf("ExecCMD() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDownloadFile(t *testing.T) {
	type args struct {
		dest string
		url  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test DownloadFile 1", args{
			dest: "pimp-my-shell-linux-amd64.gz",
			url:  "https://github.com/mr-pmillz/pimp-my-shell/releases/download/v1.5.6/pimp-my-shell-linux-amd64.gz",
		}, false},
		{"Test DownloadFile 2 non-existent", args{
			dest: "pimp-my-shell-linux-amd64.gz",
			url:  "https://notarealwebsite.notreal.fake.com/fakefile.gz",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DownloadFile(tt.args.dest, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCopyStringToFile(t *testing.T) {
	type args struct {
		data string
		dest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test CopyStringToFile 1", args{
			data: "this is a test string to write to file",
			dest: "test/testWriteStringFile.txt",
		}, false},
		{"Test CopyStringToFile 2 non-existent directory path", args{
			data: "this is a test string to write to file",
			dest: "/fake/path/that/definitely/doesnt/exists/i/hope/testWriteStringFile.txt",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyStringToFile(tt.args.data, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("CopyStringToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	t.Cleanup(func() {
		if exists, err := Exists("test/testWriteStringFile.txt"); err == nil && exists {
			if err = os.Remove("test/testWriteStringFile.txt"); err != nil {
				t.Errorf("couldnt remove test file: %v", err)
			}
		}
	})
}

func TestCmdExec(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Test CmdExec 1", args{args: []string{"ls", "-la"}}, "", false},
		{"Test CmdExec 2", args{args: []string{"fakeCommandThatDoesntExist", "--foo", "--bar"}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CmdExec(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CmdExec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf("") {
				t.Errorf("CmdExec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunCommands(t *testing.T) {
	type args struct {
		cmds []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test RunCommands 1", args{cmds: []string{"whoami", "ls -la"}}, false},
		{"Test RunCommands 2", args{cmds: []string{"pwd", "uname -a"}}, false},
		{"Test RunCommands 3", args{cmds: []string{"id", "python3 -c 'print(\"hello world\")'"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunCommands(tt.args.cmds); (err != nil) != tt.wantErr {
				t.Errorf("RunCommands() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
