// NOTE: see deps https://github.com/junegunn/fzf/blob/master/go.mod
// https://github.com/bbatsov/projectile for inspiration (not intend to make features 1-to-1)
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

// TODO-DONE: Add [path] parameter or --path [path] idk thinking on it
// I think going with the former is nicer
// TODO-DEFERRED: now that we have "basic" functionality, the project intended to do add TDD

// Intermission: Start version control to github without SEO information to make in indiscoverable with version scheme and build instructions

// TODO-NEXT-NEXT-NEXT-NEXT: Design structure of a "Project"
// TODO-NEXT-NEXT-NEXT-NEXT-NEXT: Build binary cache and state
// TODO-NEXT-NEXT-NEXT-NEXT-NEXT-NEXT: DB Operations on basic project cache
// TODO-NEXT-NEXT-NEXT-NEXT-NEXT-NEXT-NEXT: Utilize cache in search algorithm
// TODO-NEXT-NEXT-NEXT-NEXT-NEXT-NEXT-NEXT-NEXT:
// MAYBE-LATER: Split match list into multiple matchers and
// MAYBE-LATER: asynchronously handle each match candidate and exit quick once one is found and has least upward move (since these are io-heavy)

func main() {
	// NOTE: should be extended using config or env variables
	cmd := &cli.Command{
		Name:  "pd",
		Usage: "project root finder",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list projects",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					// fmt.Println("listing test: ", cmd.Args().First())
					// listProjects()
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
					fmt.Println("adding test: ", cmd.Args().First())
					// normalizePath(cmd.Args().First())
					err := Add("/Users/ianpark/Development/")
					if err != nil {
						fmt.Println("awesome")
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
					fmt.Println("removing test: ", cmd.Args().First())
					// normalizePath(cmd.Args().First())
					// addProject()
					err := Remove("/Users/ianpark/Development/")
					if err != nil {
						log.Fatal(err)
					}
					return nil
				},
			},
			// {
			// 	Name:    "blacklist",
			// 	Aliases: []string{"a"},
			// 	Usage:   "add project",
			// 	Action: func(ctx context.Context, cmd *cli.Command) error {
			// 		fmt.Println("adding test: ", cmd.Args().First())
			// 		// normalizePath(cmd.Args().First())
			// 		addProject()
			// 		return nil
			// 	},
			// },
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// pathArg := cmd.Args().Get(0)
			// var cwd string
			// var err error
			//
			// if pathArg == "" {
			// 	cwd, err = os.Getwd()
			// 	if err != nil {
			// 		fmt.Printf("Failed %v\n", err)
			// 	}
			// } else {
			// 	cwd = pathArg
			// }
			//
			normPath, err := normalizePath(cmd.Args().Get(0))
			if err != nil {
				log.Fatal(err)
			}

			foundPath := searchAncestors(normPath, getProjectFiles())
			fmt.Println(foundPath)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
