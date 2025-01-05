package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	command := &cli.Command{
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
					cmd.List()
					return nil
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add project",
				Action: func(ctx context.Context, command *cli.Command) error {
					return cmd.Add(command.Args().Get(0))
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "remove project",
				Action: func(ctx context.Context, command *cli.Command) error {
					return cmd.Remove(command.Args().Get(0))
				},
			},
			{
				Name:            "exec",
				Usage:           "Execute a command with arguments",
				SkipFlagParsing: true,
				ArgsUsage:       "[PATH] COMMAND [ARGS...]",
				Action: func(ctx context.Context, command *cli.Command) error {
					return cmd.Exec(pathFlag, command.Args().Slice())
				},
			},
			{
				Name:    "run",
				Aliases: []string{"ru"},
				Usage:   "run project command",
				Action: func(ctx context.Context, command *cli.Command) error {
					return cmd.Run(command.Args().Get(0))
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
					return cmd.Update(command.Args().Get(0))
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
							return cmd.AddBlacklist(command.Args().Get(0))
						},
					},
					{
						Name:  "remove",
						Usage: "remove from blacklist",
						Action: func(ctx context.Context, command *cli.Command) error {
							return cmd.RemoveBlacklist(command.Args().Get(0))
						},
					},
					{
						Name:  "show",
						Usage: "show blacklist",
						Action: func(ctx context.Context, command *cli.Command) error {
							return cmd.ListBlacklist()
						},
					},
				},
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			projectPath := command.Args().Get(0)
			// TOOD Should nested stuff also increment? If so, find project here too
			p, err := try_get_project(projectPath)
			if err != nil {
				return err
			}

			state.GetConfig()

			fmt.Println(p.Path)

			return nil
		},
	}

	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
