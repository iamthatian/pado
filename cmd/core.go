package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/duckonomy/parkour/project"
	"github.com/duckonomy/parkour/state"
)

var ps state.ProjectState

func Init() error {
	if err := ps.LoadState(); err != nil {
		return err
	}
	return nil
}

func tryGetProject(path string) (project.Project, error) {
	p, err := ps.GetProject(path)
	if err != nil {
		return p, fmt.Errorf("%v", err)
	}

	if p.IsEmpty() {
		if err := p.FindProjectRoot(path); err != nil {
			return p, fmt.Errorf("%v", err)
		}
	}

	return p, nil
}

func List() {
	for _, p := range ps.ListProjects() {
		fmt.Println(p.Path)
	}
}

func Add(path string) error {
	p := project.Project{}
	var projectPath string
	if err := p.InitProject(path); err != nil {
		return err
	}

	// NOTE: Get the path in argument if there is no project? (this is bad ig)
	// if p.Path == "/" {
	// 	projectPath = cmd.Args().Get(1)
	// } else {
	// 	projectPath = p.Path
	// }
	projectPath = p.Path

	if err := ps.AddProject(projectPath); err != nil {
		return err
	}

	return nil
}

func Remove(path string) error {
	if err := ps.RemoveProject(path); err != nil {
		return err
	}
	return nil
}

func Exec(path string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("command required")
	}
	var projectPath string

	p, err := tryGetProject(path)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	projectPath = p.Path

	// Create and configure command
	runner := exec.Command(args[0], args[1:]...)
	runner.Dir = projectPath
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	runner.Stdin = os.Stdin

	// Execute command
	return runner.Run()
}

func Run(path string) error {
	if err := ps.RunProject(path); err != nil {
		return err
	}
	return nil
}

func Update(path string) error {
	if err := ps.UpdateProject(path, "BuildCommand", "go run ."); err != nil {
		return err
	}
	return nil
}

func ListBlacklist() error {
	blacklist, err := ps.ShowBlacklist()
	if err != nil {
		return err
	}
	for _, path := range blacklist {
		fmt.Println(path)
	}
	return nil
}

func RemoveBlacklist(path string) error {
	if err := ps.ManageBlacklist(path, false); err != nil {
		return err
	}
	return nil
}

func AddBlacklist(path string) error {
	if err := ps.ManageBlacklist(path, true); err != nil {
		return err
	}
	return nil
}

func FindFile(path string) error {
	pf := NewProjectFinder()
	// TODO: Should nested stuff also increment? If so, find project here too
	p, err := tryGetProject(path)
	if err != nil {
		return err
	}
	file, err := pf.FindFile(p.Path)
	if err != nil {
		return err
	}
	fmt.Println(file)
	return nil
}

func FindProject() error {
	pf := NewProjectFinder()
	file, err := pf.FindProject()
	if err != nil {
		return err
	}
	fmt.Println(file)
	return nil
}

func GrepFile(path string) error {
	pf := NewProjectFinder()
	// TODO: Should nested stuff also increment? If so, find project here too
	p, err := tryGetProject(path)
	if err != nil {
		return err
	}
	err = pf.GrepEdit(p.Path)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Fix name
func Main(path string) error {
	p, err := tryGetProject(path)
	if err != nil {
		return err
	}

	err = state.GetConfig()
	if err != nil {
		return err
	}

	fmt.Println(p.Path)
	return nil
}

// func Get() {
// }
// func StateEdit() {
// }
