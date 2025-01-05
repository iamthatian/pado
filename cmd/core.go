package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/duckonomy/parkour/project"
	"github.com/duckonomy/parkour/state"
)

var ps state.ProjectState

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

func Add(path string) {
	p := project.Project{}
	var projectPath string
	if err := p.InitProject(path); err != nil {
		log.Fatal(err)
	}

	// NOTE: Get the path in argument if there is no project? (this is bad ig)
	// if p.Path == "/" {
	// 	projectPath = cmd.Args().Get(1)
	// } else {
	// 	projectPath = p.Path
	// }
	projectPath = p.Path

	if err := ps.AddProject(projectPath); err != nil {
		log.Fatal(err)
	}
}

func Remove(path string) {
	if err := ps.RemoveProject(path); err != nil {
		log.Fatal(err)
	}
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

func Run(path string) {
	if err := ps.RunProject(path); err != nil {
		log.Fatal(err)
	}
}

func Update(path string) {
	if err := ps.UpdateProject(path, "BuildCommand", "go run ."); err != nil {
		log.Fatal(err)
	}
}

func ListBlacklist() {
	blacklist, err := ps.ShowBlacklist()
	if err != nil {
		log.Fatal(err)
	}
	for _, path := range blacklist {
		fmt.Println(path)
	}
}

func RemoveBlacklist(path string) {
	if err := ps.ManageBlacklist(path, false); err != nil {
		log.Fatal(err)
	}
}

// func Get() {
// }
// func StateEdit() {
// }
