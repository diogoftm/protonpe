package handlers

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/diogoftm/protonpe/internal/crypto"
)

//
// Commands
//

// `aliases` command handler
func Harden(ctx context.Context, c *cli.Command) error {
	inFilePath := c.StringArg("<sourceFilePath>")
	outFilePath := c.StringArg("<destinationFilePath>")

	if inFilePath == "" || outFilePath == "" {
		return fmt.Errorf("Mandatory arguments not provided: <sourceFilePath>, <destinationFilePath>")
	}

	err := crypto.EncryptFileContent(inFilePath, outFilePath, c.String("encryptionMethod"))

	if err != nil {
		return err
	} else {
		return printHarden(inFilePath, outFilePath)
	}
}

//
// Output
//

// Print encryption success info
func printHarden(inFilePath string, outFilePath string) error {
	_, err := fmt.Printf("✓ Encrypted '%s' → '%s'\n", inFilePath, outFilePath)

	return err
}
