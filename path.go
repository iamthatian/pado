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

func Wd(path string) (string, error) {
	var err error
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

func NormalizePath(path string) (string, error) {
	wd, err := Wd(path)
	if err != nil {
		return "", err
	}

	var fullPath string
	if !filepath.IsAbs(wd) {
		absPath, err := filepath.Abs(wd)
		if err != nil {
			return "", err
		}
		fullPath = absPath
	} else {
		fullPath = filepath.Clean(wd)
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
