package localio

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/klauspost/cpuid/v2"
	"github.com/schollz/progressbar/v3"
	"github.com/tidwall/gjson"
)

const (
	darwin = "darwin"
	linux  = "linux"
)

// GitClone clones a public git repo url to directory
func GitClone(url, directory string) error {
	if exists, err := Exists(directory); err == nil && !exists {
		Info("git clone %s %s", url, directory)
		_, err := git.PlainClone(directory, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		CheckIfError(err)
	} else {
		fmt.Printf("[+] Repo: %s already exists at %s, skipping... \n", url, directory)
	}

	return nil
}

func sedReplaceKeysValue(sedName, varName, val, configPath string) error {
	sedPath, exists := CommandExists(sedName)
	if exists && sedPath != "" {
		if err := RunCommandPipeOutput(fmt.Sprintf("%s -i 's/%s=.*/%s=\"%s\"/' %s", sedPath, varName, varName, val, configPath)); err != nil {
			return err
		}
	}
	return nil
}

// SetVariableValue ...
func SetVariableValue(varName, val, osType, configPath string) error {
	cfgPath, err := ResolveAbsPath(configPath)
	if err != nil {
		return err
	}
	switch osType {
	case darwin:
		if err := sedReplaceKeysValue("gsed", varName, val, cfgPath); err != nil {
			return err
		}
	case linux:
		if err := sedReplaceKeysValue("sed", varName, val, cfgPath); err != nil {
			return err
		}
	default:
		if err := sedReplaceKeysValue("sed", varName, val, cfgPath); err != nil {
			return err
		}
	}
	return nil
}

// ResolveAbsPath ...
func ResolveAbsPath(path string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return path, err
	}

	dir := usr.HomeDir
	if path == "~" {
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return path, err
	}

	return path, nil
}

// Exists returns whether the given file or directory exists
func Exists(path string) (bool, error) {
	absPath, err := ResolveAbsPath(path)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(absPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CopyStringToFile ...
func CopyStringToFile(data, dest string) error {
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = destFile.WriteString(data)
	return err
}

// DownloadFile ...
func DownloadFile(dest, url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fmt.Sprintf("Downloading %s", filepath.Base(url)),
	)
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	return err
}

// GetCPUType Returns the CPU type for the current runtime environment
func GetCPUType() string {
	cpuid.Detect()
	if cpuid.CPU.VendorID.String() == "AMD" || cpuid.CPU.VendorID.String() == "Intel" && cpuid.CPU.CacheLine == 64 {
		return "AMD64"
	} else if cpuid.CPU.VendorID.String() == "ARM" && cpuid.CPU.CacheLine == 64 {
		return "ARM64"
	} else {
		return ""
	}
}

// DownloadAndInstallLatestVersionOfGolang Only for linux x86_64. Mac uses homebrew
func DownloadAndInstallLatestVersionOfGolang(homeDir string) error {
	if !CorrectOS(linux) {
		return nil
	}
	req, err := http.NewRequest("GET", "https://golang.org/VERSION?m=text", nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	goversion, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	latestGoVersion := string(goversion)

	switch GetCPUType() {
	case "ARM64":
		armGoURL := fmt.Sprintf("https://dl.google.com/go/%s.linux-arm64.tar.gz", latestGoVersion)
		dest := fmt.Sprintf("%s/%s", homeDir, path.Base(armGoURL))
		if err = DownloadFile(dest, armGoURL); err != nil {
			return err
		}
		// // Now extract the go binary. pimp-my-shell ensures that ~/.zshrc will already have the path setup for you
		if err = RunCommandPipeOutput(fmt.Sprintf("sudo rm -rf /usr/local/go 2>/dev/null && sudo tar -C /usr/local -xzf %s || true", dest)); err != nil {
			return err
		}

	case "AMD64":
		amdGoURL := fmt.Sprintf("https://dl.google.com/go/%s.linux-amd64.tar.gz", latestGoVersion)
		dest := fmt.Sprintf("%s/%s", homeDir, path.Base(amdGoURL))
		if err = DownloadFile(dest, amdGoURL); err != nil {
			return err
		}

		if err = RunCommandPipeOutput(fmt.Sprintf("sudo rm -rf /usr/local/go 2>/dev/null && sudo tar -C /usr/local -xzf %s || true", dest)); err != nil {
			return err
		}
	default:
		fmt.Println("[-] Couldn't download and install golang Unsupported CPU... Pimp-My-Shell only supports 64 bit")
	}

	return nil
}

// RunCommands ...
func RunCommands(cmds []string) error {
	for _, c := range cmds {
		fmt.Printf("[+] %s\n", c)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		bashPath, err := exec.LookPath("bash")
		if err != nil {
			return err
		}
		cmd := exec.Command(bashPath, "-c", c)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Env = os.Environ()
		err = cmd.Run()
		if err != nil {
			return err
		}
		fmt.Println(stdout.String(), stderr.String())
	}
	return nil
}

// CmdExec Execute a command and return stdout
func CmdExec(args ...string) (string, error) {
	baseCmd, err := exec.LookPath(args[0])
	if err != nil {
		return "", err
	}
	cmdArgs := args[1:]

	cmd := exec.Command(baseCmd, cmdArgs...)
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// EmbedFileCopy ...
func EmbedFileCopy(dst string, src fs.File) error {
	destFilePath, err := ResolveAbsPath(dst)
	if err != nil {
		return err
	}

	if exists, err := Exists(filepath.Dir(destFilePath)); err == nil && !exists {
		if err = os.MkdirAll(filepath.Dir(destFilePath), 0750); err != nil {
			return err
		}
	}

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(destFile, src); err != nil {
		return err
	}

	return nil
}

// EmbedFileStringAppendToDest takes a slice of bytes and writes it as a string to dest file path
func EmbedFileStringAppendToDest(data []byte, dest string) error {
	fmt.Printf("[+] Appending: \n%s\n-> %s\n", string(data), dest)
	fileDest, err := ResolveAbsPath(dest)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(fileDest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(string(data))
	if err != nil {
		return err
	}

	return nil
}

// EmbedFileStringPrependToDest takes a slice of bytes and prepend writes it as a string
// to the beginning of the dest file path
func EmbedFileStringPrependToDest(data []byte, dest string) error {
	fmt.Printf("[+] Prepending: \n%s\n-> %s\n", string(data), dest)
	fileDest, err := ResolveAbsPath(dest)
	if err != nil {
		return err
	}

	if err = NewRecord(fileDest).PrependStringToFile(string(data)); err != nil {
		return err
	}

	return nil
}

// Record is a type for prepending string text to a file
type Record struct {
	Filename string
	Contents []string
}

// NewRecord returns the Record type
func NewRecord(filename string) *Record {
	return &Record{
		Filename: filename,
		Contents: make([]string, 0),
	}
}

// readLines reads the lines of a file and appends them to Record.Contents
func (r *Record) readLines() error {
	if _, err := os.Stat(r.Filename); err != nil {
		return err
	}

	f, err := os.OpenFile(r.Filename, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if len(strings.TrimSpace(line)) != 0 {
			r.Contents = append(r.Contents, line)
			continue
		}

		if len(r.Contents) != 0 {
			r.Contents = append(r.Contents, line)
		}
	}

	return nil
}

// PrependStringToFile prepends a given string to an existing file while preserving the original formatting
func (r *Record) PrependStringToFile(content string) error {
	err := r.readLines()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(r.Filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	_, err = writer.WriteString(fmt.Sprintf("%s\n", content))
	if err != nil {
		return err
	}
	for _, line := range r.Contents {
		_, err = writer.WriteString(fmt.Sprintf("%s\n", line))
		if err != nil {
			return err
		}
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}

// ExecCMD Execute a command
func ExecCMD(command string) (string, error) {
	fmt.Printf("[+] %s\n", command)
	bashPath, err := exec.LookPath("bash")
	if err != nil {
		return "", err
	}

	cmd := exec.Command(bashPath, "-c", command)
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// RunCommandPipeOutput ...
func RunCommandPipeOutput(command string) error {
	fmt.Printf("[+] %s\n", command)
	bashPath, err := exec.LookPath("bash")
	if err != nil {
		return err
	}

	timeout := 60

	var cancel context.CancelFunc
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, bashPath, "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	errScanner := bufio.NewScanner(stderr)
	go func() {
		for errScanner.Scan() {
			fmt.Printf("%s\n", errScanner.Text())
		}
	}()

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "MONO_GAC_PREFIX=\"/usr/local\"")
	cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
	if err = cmd.Start(); err != nil {
		return err
	}

	if err = cmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Error waiting for Cmd %s\n", err)
		return err
	}

	return nil
}

// StartTmuxSession ...
func StartTmuxSession() error {
	cmd := exec.Command("tmux", "new-session", "-d")
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Directories ...
type Directories struct {
	HomeDir string
}

// CommandExists ...
func CommandExists(cmd string) (string, bool) {
	cmdPath, err := exec.LookPath(cmd)
	if err != nil {
		return "", false
	}
	return cmdPath, true
}

// NewDirectories ...
func NewDirectories() (*Directories, error) {
	d := &Directories{}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	d.HomeDir = homedir
	return d, nil
}

// Contains checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// CorrectOS ... Useful for go tests
func CorrectOS(osType string) bool {
	operatingSystem := runtime.GOOS
	return operatingSystem == osType
}

// BrewInstallProgram ...
func BrewInstallProgram(brewName, binaryName string, packages *InstalledPackages) error {
	if !CorrectOS(darwin) {
		return nil
	}
	if _, exists := CommandExists(binaryName); exists && Contains(packages.BrewInstalledPackages.Names, brewName) {
		return nil
	}
	log.Printf("[+] Installing %s\nbrew install %s\n", binaryName, brewName)
	command := fmt.Sprintf("brew install %s || true", brewName)
	if err := RunCommandPipeOutput(command); err != nil {
		return err
	}
	return nil
}

// BrewTap ...
func BrewTap(brewTap string, packages *InstalledPackages) error {
	if !CorrectOS(darwin) {
		return nil
	}
	if Contains(packages.BrewInstalledPackages.Taps, brewTap) {
		return nil
	}
	log.Printf("[+] Tapping %s\nbrew tap %s\n", brewTap, brewTap)
	command := fmt.Sprintf("brew tap %s", brewTap)
	if err := RunCommandPipeOutput(command); err != nil {
		return err
	}
	return nil
}

// BrewInstallCaskProgram ...
func BrewInstallCaskProgram(brewName, brewFullName string, packages *InstalledPackages) error {
	if !CorrectOS(darwin) {
		return nil
	}
	if Contains(packages.BrewInstalledPackages.CaskFullNames, brewFullName) {
		return nil
	}
	log.Printf("[+] Installing %s\nbrew install --cask %s\n", brewName, brewName)
	command := fmt.Sprintf("brew install --cask %s", brewName)
	if err := RunCommandPipeOutput(command); err != nil {
		return err
	}
	return nil
}

// InstalledPackages ...
type InstalledPackages struct {
	AptInstalledPackages  *AptInstalled
	BrewInstalledPackages *BrewInstalled
}

// AptInstalled ...
type AptInstalled struct {
	Name []string
}

// NewAptInstalled ...
func NewAptInstalled() (*AptInstalled, error) {
	if !CorrectOS(linux) {
		return nil, nil
	}
	var ai = &AptInstalled{}
	// installed cli to array ["blah", "foot"]
	cmd := "apt list --installed 2>/dev/null | cut -d / -f1"
	aptInstalledLines, err := ExecCMD(cmd)
	if err != nil {
		return nil, err
	}
	installedList := strings.Split(aptInstalledLines, "\n")
	if err != nil {
		return nil, err
	}

	ai.Name = append(ai.Name, installedList...)

	return ai, nil
}

// AptInstall ...
func AptInstall(packages *InstalledPackages, aptName ...string) error {
	if !CorrectOS(linux) {
		return nil
	}
	var notInstalled []string
	for _, name := range aptName {
		if Contains(packages.AptInstalledPackages.Name, name) {
			continue
		} else {
			notInstalled = append(notInstalled, name)
		}
	}
	if len(notInstalled) >= 1 {
		packagesToInstall := strings.Join(notInstalled, " ")
		command := fmt.Sprintf("sudo apt-get update -y && sudo apt-get -q install -y %s", packagesToInstall)
		if err := RunCommandPipeOutput(command); err != nil {
			return err
		}
	}

	return nil
}

// BrewInstalled ...
type BrewInstalled struct {
	Names         []string
	CaskFullNames []string
	Taps          []string
}

// NewBrewInstalled ...
func NewBrewInstalled() (*BrewInstalled, error) {
	if !CorrectOS(darwin) {
		return nil, nil
	}
	var brewInfo = &BrewInstalled{}
	brewJSON, err := CmdExec("brew", "info", "--json=v2", "--installed")
	if err != nil {
		return nil, err
	}
	result := gjson.Get(brewJSON, "formulae.#.name")
	for _, name := range result.Array() {
		brewInfo.Names = append(brewInfo.Names, name.String())
	}
	tapResult := gjson.Get(brewJSON, "casks.#.tap")
	for _, tap := range tapResult.Array() {
		if !Contains(brewInfo.Taps, tap.String()) {
			brewInfo.Taps = append(brewInfo.Taps, tap.String())
		}
	}
	fullNameResult := gjson.Get(brewJSON, "casks.#.token")
	for _, fullName := range fullNameResult.Array() {
		brewInfo.CaskFullNames = append(brewInfo.CaskFullNames, fullName.String())
	}

	return brewInfo, nil
}
