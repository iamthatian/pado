package main

import (
	"context"
	"fmt"
	"log"
	"os"
	// "os/exec"

	"github.com/urfave/cli/v3"
)

var ps ProjectState

func main() {
	// TODO: Load Config as well and pass in to LoadState
	if err := ps.LoadState(); err != nil {
		log.Fatal(err)
	}

	try_get_project := func(path string) (Project, error) {
		project, err := ps.GetProject(path)
		if err != nil {
			return project, fmt.Errorf("%v", err)
		}

		if project.IsEmpty() {
			if err := project.FindProject(path); err != nil {
				return project, fmt.Errorf("%v", err)
			}
		}

		return project, nil
	}

	cmd := &cli.Command{
		Name:  "pd",
		Usage: "project root finder",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list projects",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					for _, project := range ps.ListProjects() {
						fmt.Println(project.Path)
					}
					return nil
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add project",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					project := Project{}
					var projectPath string
					if err := project.FindProject(cmd.Args().Get(0)); err != nil {
						log.Fatal(err)
					}
					if project.Path == "/" {
						projectPath = cmd.Args().Get(1)
					} else {
						projectPath = project.Path
					}

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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if err := ps.RemoveProject(cmd.Args().Get(0)); err != nil {
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					args := cmd.Args().Slice()
					if len(args) < 1 {
						return cli.Exit("Command required", 1)
					}

					// var cmdArgs []string
					// var projectPath string
					//
					// Check if first arg is a path
					// what if path's name is ls?
					// 	possiblePath := args[0]
					// 	if absPath, err := filepath.Abs(expandPath(possiblePath)); err == nil {
					// 		if fi, err := os.Stat(absPath); err == nil && fi.IsDir() {
					// 			projectPath = absPath
					// 			cmdArgs = args[1:] // Skip the path in command args
					// 		} else {
					// 			project, err := try_get_project("")
					// 			if err != nil {
					// 				log.Fatal(err)
					// 			}
					// 			projectPath = project.Path
					// 			cmdArgs = args
					// 		}
					// 	} else {
					// 		project, err := try_get_project("")
					// 		if err != nil {
					// 			log.Fatal(err)
					// 		}
					// 		projectPath = project.Path
					// 		cmdArgs = args
					// 	}
					//
					// 	runner := exec.Command(cmdArgs[0], cmdArgs[1:]...)
					// 	runner.Stdout = os.Stdout
					// 	runner.Stderr = os.Stderr
					// 	runner.Stdin = os.Stdin
					// 	runner.Dir = projectPath
					// 	return runner.Run()
					fmt.Println("awesome")
					return nil
				},
			},
			// {
			// 	Name:            "exec",
			// 	Usage:           "Execute a command with arguments",
			// 	SkipFlagParsing: true, // Only skip for exec command
			// 	ArgsUsage:       "COMMAND [ARGS...]",
			// 	Action: func(ctx context.Context, cmd *cli.Command) error {
			// 		if cmd.Args().Len() < 1 {
			// 			return cli.Exit("Command required", 1)
			// 		}
			//
			// 		project, err := try_get_project("")
			// 		if err != nil {
			// 			log.Fatal(err)
			// 		}
			//
			// 		args := cmd.Args().Slice()
			//
			// 		runner := exec.Command(args[0], args[1:]...)
			// 		runner.Stdout = os.Stdout
			// 		runner.Stderr = os.Stderr
			// 		runner.Stdin = os.Stdin
			// 		runner.Dir = project.Path
			// 		return runner.Run()
			// 	},
			// },
			//
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "update project",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if err := ps.UpdateProject(cmd.Args().Get(0), "name", "Awesome"); err != nil {
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
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if err := ps.ManageBlacklist(cmd.Args().Get(0), true); err != nil {
								log.Fatal(err)
							}
							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "remove from blacklist",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if err := ps.ManageBlacklist(cmd.Args().Get(0), false); err != nil {
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
			project, err := try_get_project(projectPath)
			if err != nil {
				log.Fatal(err)
			}

			// project, err := ps.GetProject(projectPath)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			//
			// if project.IsEmpty() {
			// 	if err := project.FindProject(projectPath); err != nil {
			// 		log.Fatal(err)
			// 	}
			// }

			fmt.Println(project.Path)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
