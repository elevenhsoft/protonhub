package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

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

}
