package handlers

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/diogoftm/protonpe/internal/crypto"
	"github.com/diogoftm/protonpe/internal/models"
	"github.com/diogoftm/protonpe/internal/utils"
)

//
// Commands
//

// `aliases` command handler
func Aliases(ctx context.Context, c *cli.Command) error {
	fileContent, err := crypto.GetFileContent(c)
	if err != nil {
		return err
	}
	defer utils.ZeroVaultFile(&fileContent)

	return printAliases(utils.GetAllItems(fileContent, c.String("vault"), "alias"))
}

//
// Output
//

// Print all available aliases
func printAliases(items []*models.Item) error {
	fmt.Printf("Found %d alias(es)\n\n", len(items))

	for _, item := range items {
		fmt.Printf("%s → %s\n",
			item.Data.Metadata.Name,
			item.AliasEmail,
		)
	}

	return nil
}
