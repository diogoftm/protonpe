package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

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

	items := utils.GetItems(fileContent, c.String("vault"), itemId, "creditCard")
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
	items, itemId, fileContent, err := getCardItems(c)
	if err != nil {
		return err
	}

	defer utils.ZeroVaultFile(fileContent)

	return printCards(items, itemId, c.Bool("dump"))
}

// `card clip` command handler
func CardClip(ctx context.Context, c *cli.Command) error {
	ttl := c.Int("ttl")
	if ttl < 0 || ttl > 30 {
		return fmt.Errorf("invalid TTL value (0 <= ttl <= 30)")
	}

	items, itemId, fileContent, err := getCardItems(c)
	if err != nil {
		return err
	}
	defer utils.ZeroVaultFile(fileContent)

	item := items[0]

	var value string

	if c.Bool("csc") {
		value = item.Data.Content.VerificationNumber
	} else if c.Bool("pin") {
		value = item.Data.Content.Pin
	} else {
		value = item.Data.Content.Number
	}

	utils.WriteToClipboard(value, ttl)

	if len(items) > 1 {
		fmt.Printf(
			"Warning: %d matches for '%s'. Using first match:\nID: %s | Holder Name: %s | Exp. Date: %s\n",
			len(items), itemId, item.ItemId,
			item.Data.Content.CardholderName,
			item.Data.Content.ExpirationDate,
		)
	}

	return nil
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
	fmt.Printf("Found %d note(s) for '%s'\n\n", len(items), itemId)

	for i, item := range items {
		if len(items) > 1 {
			fmt.Printf("=== Card %d (%s) ===\n", i+1, item.ItemId)
		}

		fmt.Println(strings.Repeat("-", 30))
		fmt.Printf("Card Holder Name: %s\n", item.Data.Content.CardholderName)
		fmt.Printf("Number: %s\n", item.Data.Content.Number)
		fmt.Printf("CSC: %s\n", item.Data.Content.VerificationNumber)
		fmt.Printf("Expiration Date: %s\n", item.Data.Content.ExpirationDate)
		fmt.Printf("PIN: %s\n", item.Data.Content.Pin)

		if dump {
			fmt.Printf("\nItemId: %s\n", item.ItemId)
			fmt.Printf("ShareId: %s\n", item.ShareId)
			fmt.Printf("CreateTime: %s\n", time.Unix(int64(item.CreateTime), 0).Format("2006-01-02 15:04:05"))
			fmt.Printf("ModifyTime: %s\n", time.Unix(int64(item.ModifyTime), 0).Format("2006-01-02 15:04:05"))
			fmt.Printf("Pinned: %t\n", item.Pinned)
			fmt.Printf("ShareCount: %d\n", item.ShareCount)
			fmt.Printf("State: %d\n", item.State)
		}
		fmt.Println(strings.Repeat("-", 30))
	}

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
