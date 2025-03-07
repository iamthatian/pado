// TODO:
package state

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/duckonomy/parkour/project"
	"github.com/duckonomy/parkour/utils"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Blacklist []string          `koanf:"blacklist"`
	Ignores   []string          `koanf:"ignores"`
	Projects  []project.Project `koanf:"projects"`
}

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

// TODO: Handle nil & handle value in a type
func checkBlacklist(conf *koanf.Koanf) []string {
	blacklist := conf.Strings("blacklist")
	if blacklist == nil {
		return make([]string, 0)
	}
	return blacklist
}

func getProjects(conf *koanf.Koanf) ([]project.Project, error) {
	var projects []project.Project

	// Unmarshal the projects array into our slice
	if err := conf.Unmarshal("projects", &projects); err != nil {
		return nil, fmt.Errorf("failed to unmarshal projects: %w", err)
	}

	return projects, nil
}

func checkIgnore(conf *koanf.Koanf) {
	ignore := conf.Get("ignore")
	if ignore == nil {
		return
	}

	ignore = ignore.([]interface{})
}

// TODO: projects should mirror Project struct
func checkProjects(conf *koanf.Koanf) {
	projects := conf.Get("projects")
	if projects == nil {
		return
	}

	projects = projects.([]interface{})
}

// TODO: change this to per OS settings
func GetConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	checkBlacklist(k)
	checkIgnore(k)
	checkProjects(k)
	fmt.Println("awesome")
	return nil
}
