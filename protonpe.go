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
		Name:            "protonpe",
		Usage:           "read secrets from Proton Pass exports",
		Version:         "v0.2.0",
		HideHelpCommand: true,
		Commands: []*cli.Command{
			{
				Name:  "aliases",
				Usage: "Retrieve all email aliases",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "path to Proton Pass export",
					},
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
				Name:  "card",
				Usage: "Retrieve card",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "path to Proton Pass export",
					},
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
						Usage: "copy card information (card number by default) to clipboard",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "<itemId>",
							},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "csc",
								Aliases: []string{"c"},
								Usage:   "copy card security code",
							},
							&cli.BoolFlag{
								Name:    "pin",
								Aliases: []string{"p"},
								Usage:   "copy PIN code",
							},
							&cli.IntFlag{
								Name:    "ttl",
								Aliases: []string{"l"},
								Usage:   "clipboard copy time-to-live before deletion",
							},
						},
						Action:       handlers.CardClip,
						OnUsageError: bypassDefaultErrorHandling,
					},
					{
						Name:         "list",
						Usage:        "list all notes",
						Action:       handlers.CardList,
						OnUsageError: bypassDefaultErrorHandling,
					},
					{
						Name:  "show",
						Usage: "print card content",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "<itemId>",
							},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "dump",
								Aliases: []string{"d"},
								Usage:   "show all content related to this card",
							},
						},
						Action:       handlers.CardShow,
						OnUsageError: bypassDefaultErrorHandling,
					},
				},
			},
			{
				Name:  "harden",
				Usage: "Enhance the security of a vault file by encrypting it again using a different method",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "encryptionMethod",
						Aliases: []string{"e"},
						Usage: `Encryption method to use. Supported options:
- TPM_ECCP256_AES256GCM_SHA256
	Encrypts the file with AES-256-GCM and protects the encryption keys using the TPM's ECC P-256 key.
	The private key never leaves the TPM, providing hardware-backed security.`,
						Value: "TPM_ECCP256_AES256GCM_SHA256",
					},
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "path to Proton Pass export",
					},
				},
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "<sourceFilePath>",
					},
					&cli.StringArg{
						Name: "<destinationFilePath>",
					},
				},
				Action:       handlers.Harden,
				OnUsageError: bypassDefaultErrorHandling,
			},
			{
				Name:  "id",
				Usage: "Retrieve identity",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "path to Proton Pass export",
					},
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
						Usage: "copy identity information (full name by default) to clipboard",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "<itemId>",
							},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "address",
								Aliases: []string{"a"},
								Usage:   "copy address",
							},
							&cli.BoolFlag{
								Name:    "email",
								Aliases: []string{"e"},
								Usage:   "copy email",
							},
							&cli.BoolFlag{
								Name:    "license",
								Aliases: []string{"li"},
								Usage:   "copy license number",
							},
							&cli.BoolFlag{
								Name:    "passport",
								Aliases: []string{"p"},
								Usage:   "copy passport number",
							},
							&cli.BoolFlag{
								Name:    "phone",
								Aliases: []string{"p"},
								Usage:   "copy main phone number",
							},
							&cli.BoolFlag{
								Name:    "ssn",
								Aliases: []string{"s"},
								Usage:   "copy social security number",
							},
							&cli.IntFlag{
								Name:    "ttl",
								Aliases: []string{"l"},
								Usage:   "clipboard copy time-to-live before deletion",
							},
							&cli.BoolFlag{
								Name:    "workEmail",
								Aliases: []string{"we"},
								Usage:   "copy email",
							},
							&cli.BoolFlag{
								Name:    "workPhone",
								Aliases: []string{"wp"},
								Usage:   "copy work phone number",
							},
						},
						Action:       handlers.IdClip,
						OnUsageError: bypassDefaultErrorHandling,
					},
					{
						Name:         "list",
						Usage:        "list all identities",
						Action:       handlers.IdList,
						OnUsageError: bypassDefaultErrorHandling,
					},
					{
						Name:  "show",
						Usage: "print identity content",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "<itemId>",
							},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "dump",
								Aliases: []string{"d"},
								Usage:   "show all content related to this identity",
							},
						},
						Action:       handlers.IdShow,
						OnUsageError: bypassDefaultErrorHandling,
					},
				},
			},
			{
				Name:  "info",
				Usage: "Show general information about the export file and vaults",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "path to Proton Pass export",
					},
				},
				Action:       handlers.Info,
				OnUsageError: bypassDefaultErrorHandling,
			},
			{
				Name:  "login",
				Usage: "Retrieve login",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "path to Proton Pass export",
					},
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
						Usage: "copy login password (default) or TOTP to clipboard",
						Arguments: []cli.Argument{
							&cli.StringArg{
								Name: "<itemId>",
							},
						},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "totp",
								Aliases: []string{"t"},
								Usage:   "copy TOTP",
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
						Name:         "list",
						Usage:        "list all logins",
						Action:       handlers.LoginList,
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
								Name:    "dump",
								Aliases: []string{"d"},
								Usage:   "show all content related to this login",
							},
							&cli.BoolFlag{
								Name:    "totp",
								Aliases: []string{"t"},
								Usage:   "show TOTP instead of credentials",
							},
						},
						Action:       handlers.LoginShow,
						OnUsageError: bypassDefaultErrorHandling,
					},
				},
			},
			{
				Name:  "note",
				Usage: "Retrieve note",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "path to Proton Pass export",
					},
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
						Name:         "list",
						Usage:        "list all notes",
						Action:       handlers.NoteList,
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
