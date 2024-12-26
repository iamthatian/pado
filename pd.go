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
					p, err := List()
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println(p)
					return nil
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add project",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("Adding...", cmd.Args().First())
					err := Add(cmd.Args().Get(1))
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"a"},
				Usage:   "add project",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("Removing...", cmd.Args().First())
					err := Remove(cmd.Args().Get(1))
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},

			{
				Name:    "update",
				Aliases: []string{"a"},
				Usage:   "add project",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("removing test: ", cmd.Args().First())
					// normalizePath(cmd.Args().First())
					// addProject()
					err := Update(cmd.Args().Get(1), "name", "Awesome")
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var found string
			pf := cmd.Args().Get(0)

			exists, err := Exists(pf)
			if err != nil {
				log.Fatal(err)
			}

			if exists {
				proj, err := Get(pf)
				found = proj.Path
				if err != nil {
					log.Fatal(err)
				}
			} else {
				found = searchAncestors(pf, getProjectFiles())
				err := Add(found)
				if err != nil {
					log.Fatal(err)
				}
			}

			fmt.Println(found)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
