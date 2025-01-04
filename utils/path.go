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

// fileExists checks if a file or directory exists at the given path
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// isDirectory checks if the path points to a directory
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// isFile checks if the path points to a regular file
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// isSymlink checks if the path points to a symbolic link
func IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

// followSymlink follows a symbolic link and returns the destination path
func FollowSymlink(path string) (string, error) {
	if !IsSymlink(path) {
		return path, nil
	}
	return os.Readlink(path)
}

// filePermissions returns the Unix-style permission bits of a file
func FilePermissions(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Mode().Perm(), nil
}

// isExecutable checks if a file is executable by the current user
func IsExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}

// isReadable checks if a file is readable by the current user
func IsReadable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0444 != 0
}

// isWritable checks if a file is writable by the current user
func IsWritable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0222 != 0
}

func IsRootDir(path string) bool {
	return path == "/"
}

func GetBase(path string) string {
	return filepath.Base(path)
}

func CheckPath(path string, isDir bool) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		if isDir {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		} else {
			parentDir := filepath.Dir(path)
			if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create parent directories for file: %w", err)
			}

			file, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			file.Close()
		}
	} else if err != nil {
		// Handle other errors (e.g., permission issues)
		return fmt.Errorf("error checking path: %w", err)
	} else {
		if isDir && !info.IsDir() {
			return fmt.Errorf("path exists but is not a directory: %s", path)
		} else if !isDir && info.IsDir() {
			return fmt.Errorf("path exists but is a directory: %s", path)
		}
	}
	return nil
}
