// TODO
package state

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/duckonomy/parkour/utils"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

func getConfigPath() (string, error) {
	var configPath string

	switch runtime.GOOS {
	case "darwin":
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configPath = filepath.Join(userHomeDir, ".config", "parkour", "parkour.toml")
	case "linux":
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		configPath = filepath.Join(userConfigDir, "parkour", "parkour.toml")
	default:
		return "", errors.New("unsupported OS")
	}

	if err := utils.CheckPath(configPath, false); err != nil {
		return "", err
	}

	return configPath, nil
}

func checkBlacklist() {
}

func checkIgnore() {
}

func checkProjects() {
}

// TODO: change this to per OS settings
func GetConfig() error {
	// change behavior for macos
	// Should be $HOME/.config/
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	fmt.Println(configPath)
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// TODO: Check if config stuff exists
	// blacklist := k.Get("blacklist").([]interface{})
	// for _, j := range blacklist {
	// 	fmt.Println(j == "awesome")
	// }
	//
	return nil
}
