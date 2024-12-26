// this should define every path operation
// walking up/normalizing etc

package main

import (
	"errors"
	"os"
	"path/filepath"
)

type ProjectPath interface {
	ChangePath()
	Get()
}

func normalizePath(path string) (string, error) {
	var err error

	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	var fullPath string
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		fullPath = absPath
	} else {
		fullPath = filepath.Clean(path)
	}

	// TODO: Catch all error for now go deeper
	if stat, err := os.Stat(fullPath); err == nil {
		if !stat.IsDir() {
			fullPath = filepath.Dir(fullPath)
		}
		return fullPath, nil
	} else {
		return "", errors.New("wrong file path")
	}
}

func getParent(path string) string {
	return filepath.Dir(path)
}

func getBase(path string) string {
	return filepath.Base(path)
}
