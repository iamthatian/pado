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
					p, err := ListProject()
					if err != nil {
						log.Fatal(err)
					}

					for _, i := range p {
						fmt.Println(i.Path)
					}
					return nil
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add project",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					err := AddProject(cmd.Args().Get(0))
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"a"},
				Usage:   "remove project",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					err := RemoveProject(cmd.Args().Get(0))
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},

			{
				Name:    "update",
				Aliases: []string{"a"},
				Usage:   "update project",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					err := UpdateProject(cmd.Args().Get(0), "name", "Awesome")
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:    "blacklist",
				Aliases: []string{"t"},
				Usage:   "options for task templates",
				Commands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add blacklist",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := AddBlacklist(cmd.Args().Get(0))
							if err != nil {
								log.Fatal(err)
							}
							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "remove blacklist",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							err := RemoveBlacklist(cmd.Args().Get(0))
							if err != nil {
								log.Fatal(err)
							}
							return nil
						},
					},
					{
						Name:  "show",
						Usage: "show blacklist",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							bl, err := ShowBlacklist()
							if err != nil {
								log.Fatal(err)
							}

							for _, i := range bl {
								fmt.Println(i)
							}
							return nil
						},
					},
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var found string
			pf := cmd.Args().Get(0)

			exists, err := ProjectExists(pf)
			if err != nil {
				log.Fatal(err)
			}

			if exists {
				proj, err := GetProject(pf)
				found = proj.Path
				if err != nil {
					log.Fatal(err)
				}
			} else {
				found = searchAncestors(pf, getProjectFiles())
			}

			fmt.Println(found)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
