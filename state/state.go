package state

import (
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/duckonomy/parkour/project"
	"github.com/duckonomy/parkour/utils"
)

func StateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library/Application Support/pk/pk.db"), nil
	case "linux":
		return filepath.Join(home, ".local/state/pk/pk.db"), nil
	default:
		return "", errors.New("unsupported OS")
	}
}

func (ps *ProjectState) LoadState() error {
	stateFilePath, err := StateFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(stateFilePath); errors.Is(err, os.ErrNotExist) {
		ps.Projects = make(map[string]project.Project)
		ps.Blacklist = make(map[string]bool)
		return nil
	} else if err != nil {
		return err
	}

	file, err := os.Open(stateFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(ps)
}

func (ps *ProjectState) SaveState() error {
	stateFilePath, err := StateFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(stateFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(ps)
}

func (ps *ProjectState) ShowBlacklist() ([]string, error) {
	if ps.Blacklist == nil {
		return nil, nil
	}

	var blacklist []string
	for path, isBlacklisted := range ps.Blacklist {
		if isBlacklisted {
			blacklist = append(blacklist, path)
		}
	}
	return blacklist, nil
}

func (ps *ProjectState) ManageBlacklist(path string, add bool) error {
	normalizedPath, err := utils.CanonicalizePath(path)
	if err != nil {
		return err
	}

	if add {
		ps.Blacklist[normalizedPath] = true
	} else {
		delete(ps.Blacklist, normalizedPath)
	}

	return ps.SaveState()
}
