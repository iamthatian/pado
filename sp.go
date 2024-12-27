package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "pd",
		Usage: "project root finder",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list projects",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var ps ProjectState
					if err := ps.LoadState(); err != nil {
						log.Fatal(err)
					}

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
					var ps ProjectState
					if err := ps.LoadState(); err != nil {
						log.Fatal(err)
					}
					if err := ps.AddProject(cmd.Args().Get(0)); err != nil {
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
					var ps ProjectState
					if err := ps.LoadState(); err != nil {
						log.Fatal(err)
					}
					if err := ps.RemoveProject(cmd.Args().Get(0)); err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "update project",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var ps ProjectState
					if err := ps.LoadState(); err != nil {
						log.Fatal(err)
					}
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
							var ps ProjectState
							if err := ps.LoadState(); err != nil {
								log.Fatal(err)
							}
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
							var ps ProjectState
							if err := ps.LoadState(); err != nil {
								log.Fatal(err)
							}
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
							var ps ProjectState
							if err := ps.LoadState(); err != nil {
								log.Fatal(err)
							}
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
			var ps ProjectState
			if err := ps.LoadState(); err != nil {
				log.Fatal(err)
			}

			projectPath := cmd.Args().Get(0)
			project, err := ps.GetProject(projectPath)
			if err != nil {
				log.Fatal(err)
			}

			if project.IsEmpty() {
				if err := project.FindProject(projectPath, 100); err != nil {
					log.Fatal(err)
				}
			}

			fmt.Println(project.Path)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
