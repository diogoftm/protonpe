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
// Helpers
//

// Get all items with of the card type
func getCardItems(c *cli.Command) ([]*models.Item, string, *models.VaultFile, error) {
	itemId := c.StringArg("<itemId>")
	if itemId == "" {
		return nil, "", &models.VaultFile{}, fmt.Errorf("itemId argument is mandatory")
	}

	filePath, err := utils.GetFilePath(c)
	if err != nil {
		return nil, "", &models.VaultFile{}, err
	}

	fileContent, err := crypto.Decrypt(filePath)
	if err != nil {
		return nil, "", &fileContent, err
	}

	items := utils.GetItems(fileContent, c.String("vault"), itemId, "card")
	if len(items) == 0 {
		return nil, "", &fileContent, fmt.Errorf("no card found with name or ID: %s", itemId)
	}

	return items, itemId, &fileContent, nil
}

//
// Commands
//

// `card show` command handler
func CardShow(ctx context.Context, c *cli.Command) error {
	return fmt.Errorf("not implemented")
}

// `login clip` command handler
func CardClip(ctx context.Context, c *cli.Command) error {
	return fmt.Errorf("not implemented")
}

// `card list` command handler
func CardList(ctx context.Context, c *cli.Command) error {
	filePath, err := utils.GetFilePath(c)
	if err != nil {
		return err
	}

	fileContent, err := crypto.Decrypt(filePath)
	if err != nil {
		return err
	}

	defer utils.ZeroVaultFile(&fileContent)

	items := utils.GetAllItems(fileContent, c.String("vault"), "creditCard")
	return printCardList(items)
}

//
// Output
//

// Print card
func printCards(items []*models.Item, itemId string, dump bool) error {

	return nil
}

// Print the list of card items
func printCardList(items []*models.Item) error {
	fmt.Printf("Found %d card(s)\n\n", len(items))

	for _, item := range items {
		fmt.Printf("%s (%s)\n",
			item.Data.Metadata.Name,
			item.ItemId,
		)
	}

	return nil
}
