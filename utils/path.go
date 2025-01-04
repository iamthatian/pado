package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func CanonicalizePath(input string) (string, error) {
	absPath, err := filepath.Abs(input)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	cleanPath := filepath.Clean(absPath)

	resolvedPath, err := filepath.EvalSymlinks(cleanPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("path does not exist: %s", cleanPath)
		}
		return "", fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	info, err := os.Stat(resolvedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("path does not exist: %s", resolvedPath)
		}
		return "", fmt.Errorf("error checking path: %w", err)
	}

	if !info.IsDir() && info.Mode().IsRegular() {
		// return parent
		return filepath.Dir(resolvedPath), nil
	} else if info.IsDir() {
		return resolvedPath, nil
	}

	return "", fmt.Errorf("unknown path type for: %s", resolvedPath)
}

func IsHomeDir(path string) bool {
	return path == os.Getenv("HOME")
}

func IsRootDir(path string) bool {
	return path == "/"
}

func GetBase(path string) string {
	return filepath.Base(path)
}
