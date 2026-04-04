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

// Get items of the "note" type from the encrypted file taking into account
// the provided flags.
func getNoteItems(c *cli.Command) ([]*models.Item, string, *models.VaultFile, error) {
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

	items := utils.GetItems(fileContent, c.String("vault"), itemId, "")
	if len(items) == 0 {
		return nil, "", &fileContent, fmt.Errorf("no note found with name or ID: %s", itemId)
	}

	return items, itemId, &fileContent, nil
}

//
// Commands
//

// `note show` command handler
func NoteShow(ctx context.Context, c *cli.Command) error {
	items, itemId, fileContent, err := getNoteItems(c)
	if err != nil {
		return err
	}

	defer utils.ZeroVaultFile(fileContent)

	return printNotes(items, itemId, c.Bool("dump"))
}

// `note clip` command handler
func NoteClip(ctx context.Context, c *cli.Command) error {
	ttl := c.Int("ttl")
	if ttl < 0 || ttl > 30 {
		return fmt.Errorf("invalid TTL value (0 <= ttl <= 30)")
	}

	items, itemId, fileContent, err := getNoteItems(c)
	if err != nil {
		return err
	}
	defer utils.ZeroVaultFile(fileContent)

	utils.WriteToClipboard(items[0].Data.Metadata.Note, ttl)
	time.Sleep(50 * time.Millisecond)

	if len(items) > 1 {
		fmt.Printf(
			"Warning: %d matches for '%s'. Using first match:\nID: %s\n",
			len(items), itemId, items[0].ItemId,
		)
	}

	return nil
}

// `note list` command handler
func NoteList(ctx context.Context, c *cli.Command) error {
	filePath, err := utils.GetFilePath(c)
	if err != nil {
		return err
	}

	fileContent, err := crypto.Decrypt(filePath)
	if err != nil {
		return err
	}
	defer utils.ZeroVaultFile(&fileContent)

	items := utils.GetAllItems(fileContent, c.String("vault"), "note")
	return printNoteList(items)
}

//
// Output
//

// Print note item content
func printNotes(items []*models.Item, itemId string, dump bool) error {
	fmt.Printf("Found %d note(s) for '%s'\n\n", len(items), itemId)

	for i, item := range items {
		if len(items) > 1 {
			fmt.Printf("=== Note %d (%s) ===\n", i+1, item.ItemId)
		}

		fmt.Println(strings.Repeat("-", 30))
		fmt.Println(item.Data.Metadata.Note)
		fmt.Println(strings.Repeat("-", 30))

		if dump {
			fmt.Printf("\nItemId: %s\n", item.ItemId)
			fmt.Printf("ShareId: %s\n", item.ShareId)
			fmt.Printf("CreateTime: %s\n", time.Unix(int64(item.CreateTime), 0).Format("2006-01-02 15:04:05"))
			fmt.Printf("ModifyTime: %s\n", time.Unix(int64(item.ModifyTime), 0).Format("2006-01-02 15:04:05"))
			fmt.Printf("Pinned: %t\n", item.Pinned)
			fmt.Printf("ShareCount: %d\n", item.ShareCount)
			fmt.Printf("State: %d\n", item.State)
		}
	}

	return nil
}

// Print the list of note items
func printNoteList(items []*models.Item) error {
	fmt.Printf("Found %d note(s)\n\n", len(items))

	for _, item := range items {
		fmt.Printf("%s (%s)\n",
			item.Data.Metadata.Name,
			item.ItemId,
		)
	}

	return nil
}
