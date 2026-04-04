package handlers

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/urfave/cli/v3"

	"github.com/diogoftm/protonpe/internal/crypto"
	"github.com/diogoftm/protonpe/internal/models"
	"github.com/diogoftm/protonpe/internal/utils"
)

//
// Helpers
//

// Get all items with of the login type
func getLoginItems(c *cli.Command) ([]*models.Item, string, *models.VaultFile, error) {
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

	items := utils.GetItems(fileContent, c.String("vault"), itemId, "login")
	if len(items) == 0 {
		return nil, "", &fileContent, fmt.Errorf("no login found with name or ID: %s", itemId)
	}

	return items, itemId, &fileContent, nil
}

//
// Commands
//

// `login show` command handler
func LoginShow(ctx context.Context, c *cli.Command) error {
	items, itemId, fileContent, err := getLoginItems(c)
	if err != nil {
		return err
	}
	defer utils.ZeroVaultFile(fileContent)

	if c.Bool("totp") {
		return printTOTP(items, itemId)
	}

	return printLogins(items, itemId, c.Bool("dump"))
}

// `login clip` command handler
func LoginClip(ctx context.Context, c *cli.Command) error {
	ttl := c.Int("ttl")
	if ttl < 0 || ttl > 30 {
		return fmt.Errorf("invalid TTL value (0 <= ttl <= 30)")
	}

	items, itemId, fileContent, err := getLoginItems(c)
	if err != nil {
		return err
	}
	defer utils.ZeroVaultFile(fileContent)

	item := items[0]

	var value string

	if c.Bool("totp") {
		if item.Data.Content.TotpUri == "" {
			return fmt.Errorf("no TOTP configured for '%s'", itemId)
		}

		u, err := url.Parse(item.Data.Content.TotpUri)
		if err != nil {
			return fmt.Errorf("invalid TOTP URI")
		}

		secret := u.Query().Get("secret")
		if secret == "" {
			return fmt.Errorf("invalid TOTP secret")
		}

		code, err := totp.GenerateCode(secret, time.Now())
		if err != nil {
			return fmt.Errorf("failed to generate TOTP")
		}

		value = code
	} else {
		value = item.Data.Content.Password
	}

	utils.WriteToClipboard(value, ttl)

	if len(items) > 1 {
		fmt.Printf(
			"Warning: %d matches for '%s'. Using first match:\nID: %s | Email: %s | Username: %s\n",
			len(items), itemId, item.ItemId,
			item.Data.Content.ItemEmail,
			item.Data.Content.ItemUsername,
		)
	}

	return nil
}

// `note list` command handler
func LoginList(ctx context.Context, c *cli.Command) error {
	filePath, err := utils.GetFilePath(c)
	if err != nil {
		return err
	}

	fileContent, err := crypto.Decrypt(filePath)
	if err != nil {
		return err
	}

	defer utils.ZeroVaultFile(&fileContent)

	items := utils.GetAllItems(fileContent, c.String("vault"), "login")
	return printLoginList(items)
}

//
// Output
//

// Print login
func printLogins(items []*models.Item, itemId string, dump bool) error {
	fmt.Printf("Found %d login(s) for '%s'\n\n", len(items), itemId)

	for i, item := range items {
		if len(items) > 1 {
			fmt.Printf("=== Login %d (%s) ===\n", i+1, item.ItemId)
		}

		fmt.Println(strings.Repeat("-", 30))
		fmt.Printf("Email: %s\n", item.Data.Content.ItemEmail)
		fmt.Printf("Username: %s\n", item.Data.Content.ItemUsername)
		fmt.Printf("Password: %s\n", item.Data.Content.Password)
		if dump {
			fmt.Printf("Urls: %s\n", item.Data.Content.Urls)
			fmt.Printf("TotpUri: %s\n", item.Data.Content.TotpUri)
			fmt.Printf("ItemId: %s\n", item.ItemId)
			fmt.Printf("ShareId: %s\n", item.ShareId)
			fmt.Printf("AliasEmail: %s\n", item.AliasEmail)
			fmt.Printf("CreateTime: %s\n", time.Unix(int64(item.CreateTime), 0).Format("2006-01-02 15:04:05"))
			fmt.Printf("ModifyTime: %s\n", time.Unix(int64(item.ModifyTime), 0).Format("2006-01-02 15:04:05"))
			fmt.Printf("Pinned: %t\n", item.Pinned)
			fmt.Printf("ShareCount: %d\n", item.ShareCount)
			fmt.Printf("State: %d\n", item.State)
			fmt.Printf("Note: %s\n", item.Data.Metadata.Note)
		}
		fmt.Println(strings.Repeat("-", 30))
	}

	return nil
}

// Print the list of login items
func printLoginList(items []*models.Item) error {
	fmt.Printf("Found %d login(s)\n\n", len(items))

	for _, item := range items {
		fmt.Printf("%s (%s)\n",
			item.Data.Metadata.Name,
			item.ItemId,
		)
	}

	return nil
}

// Print item's totp
func printTOTP(items []*models.Item, itemId string) error {
	showId := false
	if len(items) > 1 {
		fmt.Printf("Warning: %d matches for %s\n", len(items), itemId)
		showId = true
	}
	for _, item := range items {
		if item.Data.Content.TotpUri == "" {
			fmt.Printf("%s: no TOTP configured\n", item.ItemId)
			continue
		}

		u, err := url.Parse(item.Data.Content.TotpUri)
		if err != nil {
			continue
		}

		secret := u.Query().Get("secret")

		code, err := totp.GenerateCode(secret, time.Now())
		if err != nil {
			continue
		}
		if showId {
			fmt.Printf("%s: %s\n", item.ItemId, code)
		} else {
			fmt.Printf("%s\n", code)
		}
	}

	return nil
}
