package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/duckonomy/parkour/cmd"
	"github.com/duckonomy/parkour/project"
	"github.com/duckonomy/parkour/state"
	"github.com/urfave/cli/v3"
)

func main() {
	var ps state.ProjectState
	var pathFlag string

	// TODO: Load Config as well and pass in to LoadState
	if err := ps.LoadState(); err != nil {
		log.Fatal(err)
	}

	try_get_project := func(path string) (project.Project, error) {
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

	cmd := &cli.Command{
		Name:  "pd",
		Usage: "project root finder",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Usage:       "Project path",
				Destination: &pathFlag,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list projects",
				Action: func(ctx context.Context, command *cli.Command) error {
					for _, p := range ps.ListProjects() {
						fmt.Println(p.Path)
					}
					return nil
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add project",
				Action: func(ctx context.Context, command *cli.Command) error {
					p := project.Project{}
					var projectPath string
					if err := p.InitProject(command.Args().Get(0)); err != nil {
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
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "remove project",
				Action: func(ctx context.Context, command *cli.Command) error {
					if err := ps.RemoveProject(command.Args().Get(0)); err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:            "exec",
				Usage:           "Execute a command with arguments",
				SkipFlagParsing: true,
				ArgsUsage:       "[PATH] COMMAND [ARGS...]",
				Action: func(ctx context.Context, command *cli.Command) error {
					args := command.Args().Slice()
					if len(args) < 1 {
						return cli.Exit("Command required", 1)
					}
					var projectPath string

					p, err := try_get_project(pathFlag)
					if err != nil {
						return cli.Exit(err.Error(), 1)
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
				},
			},
			{
				Name:    "run",
				Aliases: []string{"ru"},
				Usage:   "run project command",
				Action: func(ctx context.Context, command *cli.Command) error {
					if err := ps.RunProject(command.Args().Get(0)); err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},

			{
				Name:  "find",
				Usage: "find project",
				Action: func(ctx context.Context, command *cli.Command) error {
					pf := cmd.NewProjectFinder()
					file, err := pf.FindProject()
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println(file)
					return nil
				},
			},
			{
				Name:  "find-file",
				Usage: "find file",
				Action: func(ctx context.Context, command *cli.Command) error {
					pf := cmd.NewProjectFinder()
					projectPath := command.Args().Get(0)
					// TOOD Should nested stuff also increment? If so, find project here too
					p, err := try_get_project(projectPath)
					if err != nil {
						log.Fatal(err)
					}
					file, err := pf.FindFile(p.Path)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println(file)
					return nil
				},
			},
			{
				Name:  "grep-file",
				Usage: "grep file",
				Action: func(ctx context.Context, command *cli.Command) error {
					pf := cmd.NewProjectFinder()
					projectPath := command.Args().Get(0)
					// TOOD Should nested stuff also increment? If so, find project here too
					p, err := try_get_project(projectPath)
					if err != nil {
						log.Fatal(err)
					}
					err = pf.GrepEdit(p.Path)
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "update project",
				Action: func(ctx context.Context, command *cli.Command) error {
					if err := ps.UpdateProject(command.Args().Get(0), "BuildCommand", "go run ."); err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:    "blacklist",
				Aliases: []string{"b"},
				Usage:   "manage blacklist",
				Commands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add to blacklist",
						Action: func(ctx context.Context, command *cli.Command) error {
							if err := ps.ManageBlacklist(command.Args().Get(0), true); err != nil {
								log.Fatal(err)
							}
							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "remove from blacklist",
						Action: func(ctx context.Context, command *cli.Command) error {
							if err := ps.ManageBlacklist(command.Args().Get(0), false); err != nil {
								log.Fatal(err)
							}
							return nil
						},
					},
					{
						Name:  "show",
						Usage: "show blacklist",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							blacklist, err := ps.ShowBlacklist()
							if err != nil {
								log.Fatal(err)
							}
							for _, path := range blacklist {
								fmt.Println(path)
							}
							return nil
						},
					},
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			projectPath := cmd.Args().Get(0)
			// TOOD Should nested stuff also increment? If so, find project here too
			p, err := try_get_project(projectPath)
			if err != nil {
				log.Fatal(err)
			}

			state.GetConfig()

			fmt.Println(p.Path)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
