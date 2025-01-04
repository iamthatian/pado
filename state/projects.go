// Everything related to go
package state

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/duckonomy/sp/project"
	"github.com/duckonomy/sp/utils"
)

type ProjectState struct {
	Projects  map[string]project.Project
	Blacklist map[string]bool
}

func (ps *ProjectState) GetProject(path string) (project.Project, error) {
	normalizedPath, err := utils.CanonicalizePath(path)
	if err != nil {
		return project.Project{}, fmt.Errorf("failed to normalize path: %w", err)
	}

	p, exists := ps.Projects[normalizedPath]
	if !exists {
		return project.Project{}, nil
	}

	// Increment priority since the project is being accessed
	if err := ps.incrementProjectPriority(normalizedPath); err != nil {
		return project.Project{}, fmt.Errorf("failed to increment project priority: %w", err)
	}

	return p, nil
}

func (ps *ProjectState) incrementProjectPriority(path string) error {
	p, exists := ps.Projects[path]
	if !exists {
		return fmt.Errorf("project does not exist: %s", path)
	}

	p.Priority++
	ps.Projects[path] = p

	return ps.SaveState()
}

func (ps *ProjectState) ListProjects() []project.Project {
	projects := make([]project.Project, 0, len(ps.Projects))
	for _, p := range ps.Projects {
		if !ps.Blacklist[p.Path] {
			projects = append(projects, p)
		}
	}

	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].Priority > projects[j].Priority
	})

	return projects
}

func (ps *ProjectState) AddProject(projectPath string) error {
	normalizedPath, err := utils.CanonicalizePath(projectPath)
	if err != nil {
		return err
	}

	if ps.Blacklist[normalizedPath] {
		return fmt.Errorf("path %s is blacklisted", normalizedPath)
	}

	if _, exists := ps.Projects[normalizedPath]; exists {
		return fmt.Errorf("project already exists: %s", normalizedPath)
	}

	ps.Projects[normalizedPath] = project.Project{
		Name: utils.GetBase(normalizedPath), //
		Path: normalizedPath,
	}

	return ps.SaveState()
}

func (ps *ProjectState) RemoveProject(projectPath string) error {
	normalizedPath, err := utils.CanonicalizePath(projectPath)
	if err != nil {
		return err
	}

	delete(ps.Projects, normalizedPath)
	return ps.SaveState()
}

// Get projects that contain term
// does not work properly
func (ps *ProjectState) FilterProject(searchTerm string) []project.Project {
	var result []project.Project
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
	p := project.Project{}

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
	if err := p.FindProjectRoot(searchPath); err != nil {
		return fmt.Errorf("failed to find project: %w", err)
	}

	// If no project found (root is "/"), use the original path
	projectPath := p.Path
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

// TODO: refactor
func (ps *ProjectState) RunProject(path string) error {
	normalizedPath, err := utils.CanonicalizePath(path)
	if err != nil {
		return err
	}
	p := ps.Projects[normalizedPath]
	if len(p.BuildCommand) == 0 {
		fmt.Println("No command for project!")
		fmt.Print("Enter project command: ")

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input := scanner.Text()

			fields := strings.Fields(input)
			p.BuildCommand = fields

			ps.Projects[normalizedPath] = p
			ps.SaveState()
		}

		// Check for scanner errors
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
		}
	} else {
		command := p.BuildCommand[0] // First element is the command
		fmt.Println("Running:", strings.Join(p.BuildCommand, " "))
		fmt.Print("Is this the right command? [y/n/(c)hange] ")
		var confirm string
		fmt.Scanln(&confirm)

		switch confirm {
		case "n", "N":
			fmt.Println("Canceled run")
		case "y", "Y":
			args := p.BuildCommand[1:] // Remaining elements are the arguments
			runner := exec.Command(command, args...)
			runner.Stdout = os.Stdout
			runner.Stderr = os.Stderr
			runner.Stdin = os.Stdin
			runner.Dir = p.Path
			err := runner.Run()
			if err != nil {
				return err
			}

		case "c", "C":
			fmt.Print("Enter project command: ")

			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				input := scanner.Text()

				fields := strings.Fields(input)
				p.BuildCommand = fields

				ps.Projects[normalizedPath] = p
				ps.SaveState()
			}

			// Check for scanner errors
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading input:", err)
			}
		default:
			fmt.Println("doing nothing...")
		}
	}

	return nil
}

func (ps *ProjectState) UpdateProject(projectPath, key, value string) error {
	normalizedPath, err := utils.CanonicalizePath(projectPath)
	if err != nil {
		return err
	}

	p, exists := ps.Projects[normalizedPath]
	if !exists {
		return fmt.Errorf("project does not exist: %s", normalizedPath)
	}

	switch key {
	case "Path":
		p.Path = value
	case "Name":
		p.Name = value
	case "Kind":
		p.SimpleType = value
	case "Priority":
		priority, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid priority value: %s", value)
		}
		p.Priority = priority
	default:
		return fmt.Errorf("unknown key: %s", key)
	}

	// Friggin rust it too difficult for this
	// here it saves
	ps.Projects[normalizedPath] = p
	return ps.SaveState()
}
