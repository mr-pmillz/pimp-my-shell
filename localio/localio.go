package localio

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	gitex "github.com/go-git/go-git/v5/_examples"
	"github.com/klauspost/cpuid/v2"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
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
		gitex.Info("git clone %s %s", url, directory)
		_, err := git.PlainClone(directory, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		gitex.CheckIfError(err)
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
	switch {
	case cpuid.CPU.VendorID.String() == "AMD" || cpuid.CPU.VendorID.String() == "Intel" && cpuid.CPU.CacheLine == 64:
		return "AMD64"
	case cpuid.CPU.VendorID.String() == "ARM" && cpuid.CPU.CacheLine == 64:
		return "ARM64"
	default:
		return ""
	}
}

// DownloadAndInstallLatestVersionOfGolang Only for linux x86_64. Mac uses homebrew
func DownloadAndInstallLatestVersionOfGolang(homeDir string, packages *InstalledPackages) error {
	if CorrectOS(linux) {
		return AptInstall(packages, "golang")
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

	return NewRecord(fileDest).PrependStringToFile(string(data))
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

	return writer.Flush()
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

// TimeTrack ...
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	Infof("%s \ntook: %s\n", name, elapsed)
}

// RunCommandPipeOutput runs a bash command and pipes the output to stdout and stderr in realtime
//
//nolint:gocognit
func RunCommandPipeOutput(command string) error {
	defer TimeTrack(time.Now(), command)
	timeout := 120

	bashPath, err := exec.LookPath("bash")
	if err != nil {
		return LogError(err)
	}

	var cancel context.CancelFunc
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, bashPath, "-c", command)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return LogError(err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return LogError(err)
	}

	// Start goroutines for reading command output and colorizing it.
	stdoutReader := bufio.NewReader(stdoutPipe)
	stderrReader := bufio.NewReader(stderrPipe)

	go func() {
		for {
			line, _, err := stdoutReader.ReadLine()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					LogWarningf("Error reading stdout: %v", err)
				}
				break
			}
			if _, err = os.Stdout.Write(line); err != nil {
				return
			}
			if _, err = os.Stdout.Write([]byte{'\n'}); err != nil {
				return
			}
		}
	}()

	go func() {
		for {
			line, _, err := stderrReader.ReadLine()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					LogWarningf("Error reading stdout: %v", err)
				}
				break
			}
			if _, err = os.Stderr.Write(line); err != nil {
				return
			}
			if _, err = os.Stderr.Write([]byte{'\n'}); err != nil {
				return
			}
		}
	}()

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "MONO_GAC_PREFIX=\"/usr/local\"")
	cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
	if err = cmd.Start(); err != nil {
		return LogError(err)
	}

	// Create a channel to receive the signal interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Use a goroutine to wait for the signal interrupt
	go func() {
		<-interrupt
		LogWarningf("[!] CTRL^C Detected. Stopping the current running command and exiting...\n%s", command)
		cancel() // Cancel the context
		// Wait for the command to finish before exiting
		if err = cmd.Wait(); err != nil {
			os.Exit(0)
		}
		os.Exit(1)
	}()

	// Wait for the command to finish
	if err = cmd.Wait(); err != nil {
		// If timeout exceeded, log a warning and return a timeout error
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			timeoutMsg := fmt.Sprintf("Timeout exceeded (%d minutes) for command: %s", timeout, command)
			LogWarningf(timeoutMsg)
			return errors.New(timeoutMsg)
		}
		_, _ = fmt.Fprintf(os.Stderr, "Error waiting for Cmd %s\n", err)
		return LogError(err)
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
	return cmd.Run()
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
	return RunCommandPipeOutput(command)
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
	return RunCommandPipeOutput(command)
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
	return RunCommandPipeOutput(command)
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
		}
		notInstalled = append(notInstalled, name)
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

// LogWarningf logs a warning to stdout
func LogWarningf(format string, args ...interface{}) {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelWarning)
	gologger.Warning().Label("WARN").Msgf(format, args...)
}

// Infof ...
func Infof(format string, args ...interface{}) {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelInfo)
	gologger.Info().Label("INFO").Msgf(format, args...)
}

// LogError ...
func LogError(err error) error {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		LogWarningf("Failed to retrieve Caller information")
	}
	fn := runtime.FuncForPC(pc).Name()
	gologger.DefaultLogger.SetMaxLevel(levels.LevelError)
	gologger.Error().Msgf("Error in function %s, called from %s:%d:\n %v", fn, file, line, err)
	return err
}

type PipInstalled struct {
	Name     []string
	Versions PythonVersions
}

type PythonVersions struct {
	Python3Version string
	PipVersion     string
}

// PipInstall ...
func PipInstall(packagesToInstall []string) error {
	installedPipPackages, err := NewPipInstalled()
	if err != nil {
		return LogError(err)
	}

	if err = InstallPipPackages(installedPipPackages, packagesToInstall...); err != nil {
		return LogError(err)
	}
	return nil
}

// InstallPipPackages ...
func InstallPipPackages(installedPackages *PipInstalled, pkgName ...string) error {
	var notInstalled []string
	for _, name := range pkgName {
		if !Contains(installedPackages.Name, name) {
			notInstalled = append(notInstalled, name)
		}
	}
	if len(notInstalled) >= 1 {
		packagesToInstall := strings.Join(notInstalled, " ")
		if IsRoot() {
			command := fmt.Sprintf("python3 -m pip install %s || true", packagesToInstall)
			if err := RunCommandPipeOutput(command); err != nil {
				return LogError(err)
			}
		} else {
			var cmd string
			if VersionGreaterOrEqual(installedPackages.Versions.Python3Version, "3.11.0") || VersionGreaterOrEqual(installedPackages.Versions.PipVersion, "22.2.3") {
				cmd = fmt.Sprintf("python3 -m pip install %s --break-system-packages || true", packagesToInstall)
			} else {
				cmd = fmt.Sprintf("python3 -m pip install %s --user || true", packagesToInstall)
			}
			if err := RunCommandPipeOutput(cmd); err != nil {
				return LogError(err)
			}
		}
	}

	return nil
}

// IsRoot checks if the current user is root or not
func IsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

// NewPipInstalled returns a slice of all the installed python3 pip packages
func NewPipInstalled() (*PipInstalled, error) {
	pip := &PipInstalled{}
	cmd := "python3 -m pip list | awk '{print $1}'"
	pipPackages, err := ExecCMD(cmd)
	if err != nil {
		return nil, LogError(err)
	}
	installedList := strings.Split(pipPackages, "\n")
	pip.Name = append(pip.Name, installedList...)
	v, err := GetPythonAndPipVersion()
	if err != nil {
		return nil, LogError(err)
	}
	pip.Versions.PipVersion = v.PipVersion
	pip.Versions.Python3Version = v.Python3Version

	return pip, nil
}

// GetPythonAndPipVersion ...
func GetPythonAndPipVersion() (*PythonVersions, error) {
	v := &PythonVersions{}

	pythonVersionCMD := "python3 --version"
	pythonVersion, err := ExecCMD(pythonVersionCMD)
	if err != nil {
		return nil, LogError(err)
	}

	pythonVersionStr := parseVersionWithRegex(pythonVersion, `Python (\d+\.\d+\.\d+)`)
	v.Python3Version = pythonVersionStr

	pipVersionCMD := "python3 -m pip -V"
	pipVersion, err := ExecCMD(pipVersionCMD)
	if err != nil {
		return nil, LogError(err)
	}

	pipVersionStr := parseVersionWithRegex(pipVersion, `pip (\d+\.\d+)`)
	v.PipVersion = pipVersionStr

	return v, nil
}

// parseVersionWithRegex ...
func parseVersionWithRegex(versionString, regx string) string {
	re := regexp.MustCompile(regx)
	match := re.FindStringSubmatch(versionString)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

type Version struct {
	Major int
	Minor int
	Patch int
}

// VersionGreaterOrEqual ...
func VersionGreaterOrEqual(versionString string, minVersionString string) bool {
	version, err := parseVersion(versionString)
	if err != nil {
		return false
	}
	minVersion, err := parseVersion(minVersionString)
	if err != nil {
		return false
	}

	// Compare the major version numbers
	if version.Major > minVersion.Major {
		return true
	} else if version.Major < minVersion.Major {
		return false
	}

	// Compare the minor version numbers
	if version.Minor > minVersion.Minor {
		return true
	} else if version.Minor < minVersion.Minor {
		return false
	}

	// Compare the patch version numbers
	if version.Patch >= minVersion.Patch {
		return true
	}
	return false
}

// parseVersion ...
func parseVersion(versionString string) (*Version, error) {
	re := regexp.MustCompile(`(\d+)\.(\d+)(?:\.(\d+))?`)
	match := re.FindStringSubmatch(versionString)
	if len(match) < 3 {
		return nil, fmt.Errorf("invalid version string")
	}
	major, err := strconv.Atoi(match[1])
	if err != nil {
		return nil, err
	}
	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return nil, err
	}
	var patch int
	if len(match) == 4 && match[3] != "" {
		patch, err = strconv.Atoi(match[3])
		if err != nil {
			return nil, err
		}
	}
	return &Version{Major: major, Minor: minor, Patch: patch}, nil
}
