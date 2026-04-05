package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/diogoftm/protonpe/internal/handlers"
)

// CLI entrypoint
func main() {

	cmd := &cli.Command{
		Name:    "protonpe",
		Usage:   "read secrets from Proton Pass exports",
		Version: "v0.1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "path to Proton Pass export",
			},
		},
		HideHelpCommand: true,
		Commands: []*cli.Command{
			{
				Name:         "info",
				Usage:        "Show general information about the export file and vaults",
				Action:       handlers.Info,
				OnUsageError: bypassDefaultErrorHandling,
			},
			{
				Name:  "aliases",
				Usage: "Retrieve all email aliases",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "vault",
						Aliases: []string{"v"},
						Usage:   "filter by vault name or ID",
					},
				},
				Action:       handlers.Aliases,
				OnUsageError: bypassDefaultErrorHandling,
			},
			{
				Name:  "note",
				Usage: "Retrieve note",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "vault",
						Aliases: []string{"v"},
						Usage:   "filter by vault name or ID",
					},
				},
				OnUsageError: bypassDefaultErrorHandling,
				Commands: []*cli.Command{
					{
						Name:  "clip",
						Usage: "copy note content to clipboard",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "<itemId>",
							},
						},
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "ttl",
								Aliases: []string{"l"},
								Usage:   "clipboard copy time-to-live before deletion",
							},
						},
						Action:       handlers.NoteClip,
						OnUsageError: bypassDefaultErrorHandling,
					},
					{
						Name:  "show",
						Usage: "print note content",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "<itemId>",
							},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "dump",
								Aliases: []string{"d"},
								Usage:   "show all content related to this login",
							},
						},
						Action:       handlers.NoteShow,
						OnUsageError: bypassDefaultErrorHandling,
					},
					{
						Name:         "list",
						Usage:        "list all notes",
						Action:       handlers.NoteList,
						OnUsageError: bypassDefaultErrorHandling,
					},
				},
			},
			{
				Name:  "login",
				Usage: "Retrieve login",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "vault",
						Aliases: []string{"v"},
						Usage:   "filter by vault name or ID",
					},
				},
				OnUsageError: bypassDefaultErrorHandling,
				Commands: []*cli.Command{
					{
						Name:  "clip",
						Usage: "copy login password or TOTP to clipboard",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "<itemId>",
							},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "totp",
								Aliases: []string{"t"},
								Usage:   "copy TOTP instead of password",
							},
							&cli.IntFlag{
								Name:    "ttl",
								Aliases: []string{"l"},
								Usage:   "clipboard copy time-to-live before deletion",
							},
						},
						Action:       handlers.LoginClip,
						OnUsageError: bypassDefaultErrorHandling,
					},
					{
						Name:  "show",
						Usage: "print login details or TOTP",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "<itemId>",
							},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "totp",
								Aliases: []string{"t"},
								Usage:   "show TOTP instead of credentials",
							},
							&cli.BoolFlag{
								Name:    "dump",
								Aliases: []string{"d"},
								Usage:   "show all content related to this login",
							},
						},
						Action:       handlers.LoginShow,
						OnUsageError: bypassDefaultErrorHandling,
					},
					{
						Name:         "list",
						Usage:        "list all logins",
						Action:       handlers.LoginList,
						OnUsageError: bypassDefaultErrorHandling,
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Used to bypass CLI default error output messages
func bypassDefaultErrorHandling(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
	return err
}
