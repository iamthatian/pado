package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
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
	normalizedPath, err := CanonicalizePath(path)
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
	normalizedPath, err := CanonicalizePath(projectPath)
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
	normalizedPath, err := CanonicalizePath(projectPath)
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
	normalizedPath, err := CanonicalizePath(path)
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

// Get projects that contain term
// does not work properly
func (ps *ProjectState) FilterProject(searchTerm string) []Project {
	var result []Project
	for _, value := range ps.Projects {
		v := reflect.ValueOf(value)

		values := make([]interface{}, v.NumField())

		for i := 0; i < v.NumField(); i++ {
			// values[i] =
			val := v.Field(i).Interface()
			// fmt.Println(v.Field(i).Interface())
			// figure out what this type is
			if reflect.TypeOf(val) == reflect.TypeOf("") {
				if strings.Contains(val.(string), searchTerm) {
					result = append(result, value)
				}
			} else if reflect.TypeOf(val) == reflect.TypeOf(1) {
				if strings.Contains(strconv.Itoa(val.(int)), searchTerm) {
					result = append(result, value)
				}
			}
			// fmt.Println(reflect.TypeOf(val))
			// if strings.Contains(v.Field(i).Interface(), searchTerm) {
			// result = append(result, value)
			// }
		}

		fmt.Println(values)
	}
	return result
}

// TODO: Define UX for this
func (ps *ProjectState) ExecProject(path string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	// Create a new project to find the root
	project := Project{}

	// If path is provided, use it, otherwise use current directory
	searchPath := path
	if searchPath == "" {
		var err error
		searchPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Find the project root
	if err := project.FindProject(searchPath); err != nil {
		return fmt.Errorf("failed to find project: %w", err)
	}

	// If no project found (root is "/"), use the original path
	projectPath := project.Path
	if projectPath == "/" {
		projectPath = searchPath
	}

	// Create and configure command
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Execute command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed in %s: %w", projectPath, err)
	}

	return nil
}

func (ps *ProjectState) RunProject(path string) error {
	normalizedPath, err := CanonicalizePath(path)
	if err != nil {
		return err
	}
	project := ps.Projects[normalizedPath]
	if len(project.BuildCommand) == 0 {
		fmt.Println(project.BuildCommand)

		fmt.Println("No command for project!")
		fmt.Print("Enter project command: ")

		fmt.Print("Enter space-separated values: ")

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input := scanner.Text()

			fields := strings.Fields(input)
			project.BuildCommand = fields

			ps.Projects[normalizedPath] = project
			ps.SaveState()
		}

		// Check for scanner errors
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
		}
	} else {
		command := project.BuildCommand[0] // First element is the command
		args := project.BuildCommand[1:]   // Remaining elements are the arguments
		runner := exec.Command(command, args...)
		runner.Stdout = os.Stdout
		runner.Stderr = os.Stderr
		runner.Stdin = os.Stdin
		runner.Dir = project.Path
		err := runner.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func (ps *ProjectState) UpdateProject(projectPath, key, value string) error {
	normalizedPath, err := CanonicalizePath(projectPath)
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
	// case "BuildCommand":
	// 	project.BuildCommand = value
	case "Priority":
		priority, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid priority value: %s", value)
		}
		project.Priority = priority
	default:
		return fmt.Errorf("unknown key: %s", key)
	}

	// Friggin rust it too difficult for this
	// here it saves
	ps.Projects[normalizedPath] = project
	return ps.SaveState()
}
