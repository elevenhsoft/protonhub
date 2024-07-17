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
	"strings"

	"github.com/BurntSushi/toml"
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

func GetConfigPath(cfg string) string {
	return filepath.Join(phStorePath(), "configs", cfg)
}

func GetLockfilePath(gameId string) string {
	filename := fmt.Sprintf("%s.lock", gameId)
	return filepath.Join(phStorePath(), "locks", filename)
}

func CreateLockfileForProcess(gameId, pid string) {
	path := GetLockfilePath(gameId)

	file, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(pid)
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

func homePath() string {
	return os.Getenv("HOME")
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

func CmdToResponse(cmd *exec.Cmd, w http.ResponseWriter) {
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err != nil {
		log.Fatal(err)
	}

	// Start the command and check for errors
	err = cmd.Start()

	if err != nil {
		log.Fatal(err)
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

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	go func() {
		// Read line by line and process it
		dataCh <- fmt.Sprintf("pid: %d", cmd.Process.Pid)
		for scanner.Scan() {
			line := scanner.Text()

			fmt.Println(line)

			if line != "" {
				dataCh <- StripANSI(line)
			}
		}
		dataCh <- "0"
	}()

	_ = cmd.Wait()
}
