package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-ps"
)

func StripANSI(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)

	return re.ReplaceAllString(str, "")
}

func ParseLauncherArgs(args string) []string {
	var launcherArgs []string

	for _, split := range strings.Split(args, ",") {
		arg := strings.TrimSpace(split)
		launcherArgs = append(launcherArgs, arg)
	}

	return launcherArgs
}

func UnParseLauncherArgs(args []string) string {
	var result string

	for _, arg := range args {
		if len(args) > 1 {
			result += arg + ", "
		} else {
			result += arg
		}
	}

	return result
}

type Launcher struct {
	Config     string
	Name       string
	Prefix     string
	Proton     string
	GameID     string
	Exe        string
	LaunchArgs []string
	Store      string
}

type umu struct {
	Prefix     string   `toml:"prefix"`
	Proton     string   `toml:"proton"`
	GameID     string   `toml:"game_id"`
	Exe        string   `toml:"exe"`
	LaunchArgs []string `toml:"launch_args"`
	Store      string   `toml:"store"`
}

func toTomlFileName(s string) string {
	prefix := rand.Intn(1000)

	var filename string

	filename = strings.ReplaceAll(s, " ", "_")
	filename = strings.TrimSpace(filename)
	filename = strings.ToLower(filename)

	return fmt.Sprintf("%d_%s.toml", prefix, filename)
}

func createTomlConfig(s string, data umu) {
	final_path := filepath.Join(phStorePath(), "configs", s)
	writer := new(bytes.Buffer)

	err := toml.NewEncoder(writer).Encode(data)

	if err != nil {
		log.Fatal(err)
	}

	final_toml := fmt.Sprintf("[umu]\n%s", writer.String())

	file, err := os.Create(final_path)

	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(final_toml)
}

func updateTomlFile(s string, data umu) {
	final_path := filepath.Join(phStorePath(), "configs", s)

	err := os.Remove(final_path)

	if err != nil {
		log.Fatal(err)
	}

	writer := new(bytes.Buffer)

	err = toml.NewEncoder(writer).Encode(data)

	if err != nil {
		log.Fatal(err)
	}

	final_toml := fmt.Sprintf("[umu]\n%s", writer.String())

	file, err := os.Create(final_path)

	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(final_toml)
}

func DeleteDataForLauncher(launcher Launcher) {
	config := GetConfigPath(launcher.Config)

	err := os.Remove(config)

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	winePrefix := translatePath(launcher.Prefix)

	err = os.RemoveAll(winePrefix)

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	conn := DbConnection()
	RemoveLauncherFromDb(conn, launcher.Config)
}

func GetConfigPath(cfg string) string {
	return filepath.Join(phStorePath(), "configs", cfg)
}

func GetLockfilePath(gameId string) string {
	filename := fmt.Sprintf("%s.lock", gameId)
	return filepath.Join(phStorePath(), "locks", filename)
}

func ListRunningGameIds() []string {
	var locks []string
	target := filepath.Join(phStorePath(), "locks")

	entries, err := os.ReadDir(target)

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	for _, entry := range entries {
		locks = append(locks, strings.Split(entry.Name(), ".lock")[0])
	}

	return locks
}

func CheckLockfileForProcess(gameId string) bool {
	path := GetLockfilePath(gameId)

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return true
	}

	return false
}

func CreateLockfileForProcess(gameId string, pid int) bool {
	path := GetLockfilePath(gameId)

	file, err := os.Create(path)

	if err != nil {
		return false
	}

	file.WriteString(fmt.Sprintf("%d", pid))

	return err == nil
}

func RemoveLockfileForProcess(gameId string) bool {
	path := GetLockfilePath(gameId)

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	os.Remove(path)
	return err == nil
}

func KillProcessByGameId(gameId string) bool {
	var pid_folder []int
	pids, err := getNewPids()

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	path := GetLockfilePath(gameId)
	_, err = os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	content, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	pid, err := strconv.Atoi(string(content))
	pid_folder = append(pid_folder, pid)

	conn := DbConnection()
	launcher := GetLauncherByIdFromDb(conn, gameId)

	target_exe := strings.Split(launcher.Exe, "/")

	for _, exe := range pids {
		if exe.Executable() == target_exe[len(target_exe)-1] {
			pid_folder = append(pid_folder, exe.Pid())
		}
	}

	for _, pid := range pid_folder {
		syscall.Kill(-pid, 15)
	}

	RemoveLockfileForProcess(gameId)

	return err == nil
}

func getNewPids() ([]ps.Process, error) {
	pids, err := ps.Processes()

	if err != nil {
		fmt.Printf("cannot get list of currently running processes")
		return nil, err
	}

	return pids, nil
}

func homePath() string {
	return os.Getenv("HOME")
}

func translatePath(path string) string {
	if strings.HasPrefix(path, "~") {
		return strings.ReplaceAll(path, "~", homePath())
	}

	return path
}

func phStorePath() string {
	return filepath.Join(homePath(), ".local/share/protonhub")
}

func initStore() {
	path := filepath.Join(phStorePath(), "configs")
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}

	path = filepath.Join(phStorePath(), "locks")
	_, err = os.Stat(path)

	if os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}

func CmdToResponse(gameId string, cmd *exec.Cmd, w http.ResponseWriter) {
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err != nil {
		log.Fatal(err)
	}

	var pgid int

	// Start the command and check for errors
	if CheckLockfileForProcess(gameId) {
		err = cmd.Start()

		if err != nil {
			log.Fatal(err)
		}

		pgid, err = syscall.Getpgid(cmd.Process.Pid)
		if err != nil {
			log.Fatal(err)
		}
	}

	dataCh := make(chan string)

	go func() {
		for data := range dataCh {
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		}
	}()
	// Create a scanner which scans stdout in a line-by-line fashion
	scanner := bufio.NewScanner(stdout)

	go func() {
		if CreateLockfileForProcess(gameId, pgid) {
			dataCh <- fmt.Sprintf("[%s] create lockfile for pgid: %d\n", gameId, pgid)
		}
	}()

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	go func() {
		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()

			fmt.Println(line)
			dataCh <- StripANSI(line)
		}
	}()

	_ = cmd.Wait()

	if RemoveLockfileForProcess(gameId) {
		fmt.Fprintf(w, "data: [exit: 0] removing lockfile for pgid: %d\n\n", pgid)
		w.(http.Flusher).Flush()
	}
}
