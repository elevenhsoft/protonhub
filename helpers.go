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
	final_path := filepath.Join(phStorePath(), s)
	writer := new(bytes.Buffer)

	err := toml.NewEncoder(writer).Encode(data)

	if err != nil {
		log.Fatal(err)
	}

	final_toml := fmt.Sprintf("[umu]\n%s", writer.String())

	initStore()
	file, err := os.Create(final_path)

	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(final_toml)
}

func homePath() string {
	return os.Getenv("HOME")
}

func phStorePath() string {
	return filepath.Join(homePath(), ".local/share/protonhub")
}

func initStore() {
	_, err := os.Stat(phStorePath())

	if os.IsNotExist(err) {
		os.MkdirAll(phStorePath(), 0755)
	}
}
