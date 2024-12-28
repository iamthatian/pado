package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
)

type ProjectState struct {
	Projects  map[string]Project
	Blacklist map[string]bool
}

func getStateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library/Application Support/sp/sp.db"), nil
	case "linux":
		return filepath.Join(home, ".local/state/sp/sp.db"), nil
	default:
		return "", errors.New("unsupported OS")
	}
}

func (ps *ProjectState) LoadState() error {
	stateFilePath, err := getStateFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(stateFilePath); errors.Is(err, os.ErrNotExist) {
		ps.Projects = make(map[string]Project)
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
	stateFilePath, err := getStateFilePath()
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

func (ps *ProjectState) GetProject(path string) (Project, error) {
	normalizedPath, err := NormalizePath(path)
	if err != nil {
		return Project{}, fmt.Errorf("failed to normalize path: %w", err)
	}

	project, exists := ps.Projects[normalizedPath]
	if !exists {
		return Project{}, nil
	}

	// Increment priority since the project is being accessed
	if err := ps.incrementProjectPriority(normalizedPath); err != nil {
		return Project{}, fmt.Errorf("failed to increment project priority: %w", err)
	}

	return project, nil
}

func (ps *ProjectState) incrementProjectPriority(path string) error {
	project, exists := ps.Projects[path]
	if !exists {
		return fmt.Errorf("project does not exist: %s", path)
	}

	project.Priority++
	ps.Projects[path] = project

	return ps.SaveState()
}

func (ps *ProjectState) ListProjects() []Project {
	projects := make([]Project, 0, len(ps.Projects))
	for _, project := range ps.Projects {
		if !ps.Blacklist[project.Path] {
			projects = append(projects, project)
		}
	}

	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].Priority > projects[j].Priority
	})

	return projects
}

func (ps *ProjectState) AddProject(projectPath string) error {
	normalizedPath, err := NormalizePath(projectPath)
	if err != nil {
		return err
	}

	if ps.Blacklist[normalizedPath] {
		return fmt.Errorf("path %s is blacklisted", normalizedPath)
	}

	if _, exists := ps.Projects[normalizedPath]; exists {
		return fmt.Errorf("project already exists: %s", normalizedPath)
	}

	ps.Projects[normalizedPath] = Project{
		Name: getBase(normalizedPath),
		Path: normalizedPath,
	}

	return ps.SaveState()
}

func (ps *ProjectState) RemoveProject(projectPath string) error {
	normalizedPath, err := NormalizePath(projectPath)
	if err != nil {
		return err
	}

	delete(ps.Projects, normalizedPath)
	return ps.SaveState()
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
	normalizedPath, err := NormalizePath(path)
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

func (ps *ProjectState) UpdateProject(projectPath, key, value string) error {
	normalizedPath, err := NormalizePath(projectPath)
	if err != nil {
		return err
	}

	project, exists := ps.Projects[normalizedPath]
	if !exists {
		return fmt.Errorf("project does not exist: %s", normalizedPath)
	}

	switch key {
	case "Path":
		project.Path = value
	case "Name":
		project.Name = value
	case "Kind":
		project.Kind = value
	case "Description":
		project.Description = value
	case "Priority":
		priority, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid priority value: %s", value)
		}
		project.Priority = priority
	default:
		return fmt.Errorf("unknown key: %s", key)
	}

	ps.Projects[normalizedPath] = project
	return ps.SaveState()
}
