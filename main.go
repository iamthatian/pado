package main

import (
	"context"
	"log"
	"os"

	"github.com/duckonomy/parkour/cmd"
	"github.com/urfave/cli/v3"
)

func main() {
	var pathFlag string

	// TODO: Load Config as well and pass in to LoadState
	if err := cmd.Init(); err != nil {
		log.Fatal(err)
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
					return cmd.FindProject()
				},
			},
			{
				Name:  "find-file",
				Usage: "find file",
				Action: func(ctx context.Context, command *cli.Command) error {
					return cmd.FindFile(command.Args().Get(0))
				},
			},
			{
				Name:  "grep-file",
				Usage: "grep file",
				Action: func(ctx context.Context, command *cli.Command) error {
					return cmd.GrepFile(command.Args().Get(0))
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
			return cmd.Main(command.Args().Get(0))
		},
	}

	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
