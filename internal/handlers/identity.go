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

// Get all items with of the identity type
func getIdentityItems(c *cli.Command) ([]*models.Item, string, *models.VaultFile, error) {
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

	items := utils.GetItems(fileContent, c.String("vault"), itemId, "identity")
	if len(items) == 0 {
		return nil, "", &fileContent, fmt.Errorf("no identity found with name or ID: %s", itemId)
	}

	return items, itemId, &fileContent, nil
}

//
// Commands
//

// `id show` command handler
func IdShow(ctx context.Context, c *cli.Command) error {
	items, itemId, fileContent, err := getIdentityItems(c)
	if err != nil {
		return err
	}

	defer utils.ZeroVaultFile(fileContent)

	return printIdentities(items, itemId, c.Bool("dump"))
}

// `id clip` command handler
func IdClip(ctx context.Context, c *cli.Command) error {
	ttl := c.Int("ttl")
	if ttl < 0 || ttl > 30 {
		return fmt.Errorf("invalid TTL value (0 <= ttl <= 30)")
	}

	items, itemId, fileContent, err := getIdentityItems(c)
	if err != nil {
		return err
	}
	defer utils.ZeroVaultFile(fileContent)

	item := items[0]

	var value string

	if c.Bool("email") {
		value = item.Data.Content.Email
	} else if c.Bool("workEmail") {
		value = item.Data.Content.WorkEmail
	} else if c.Bool("phone") {
		value = item.Data.Content.PhoneNumber
	} else if c.Bool("workPhone") {
		value = item.Data.Content.WorkPhoneNumber
	} else if c.Bool("ssn") {
		value = item.Data.Content.SocialSecurityNumber
	} else if c.Bool("address") {
		value = item.Data.Content.StreetAddress
	} else if c.Bool("passport") {
		value = item.Data.Content.PassportNumber
	} else if c.Bool("license") {
		value = item.Data.Content.LicenseNumber
	} else {
		value = item.Data.Content.FullName
	}

	utils.WriteToClipboard(value, ttl)

	if len(items) > 1 {
		fmt.Printf(
			"Warning: %d matches for '%s'. Using first match:\nID: %s\n",
			len(items), itemId, item.ItemId,
		)
	}

	return nil
}

// `id list` command handler
func IdList(ctx context.Context, c *cli.Command) error {
	filePath, err := utils.GetFilePath(c)
	if err != nil {
		return err
	}

	fileContent, err := crypto.Decrypt(filePath)
	if err != nil {
		return err
	}

	defer utils.ZeroVaultFile(&fileContent)

	items := utils.GetAllItems(fileContent, c.String("vault"), "identity")
	return printIdentityList(items)
}

//
// Output
//

// Print card
func printIdentities(items []*models.Item, itemId string, dump bool) error {
	fmt.Printf("Found %d ID(s) for '%s'\n\n", len(items), itemId)

	for i, item := range items {
		if len(items) > 1 {
			fmt.Printf("=== Card %d (%s) ===\n", i+1, item.ItemId)
		}

		fmt.Println(strings.Repeat("-", 30))
		fmt.Printf("Full Name: %s\n", item.Data.Content.FullName)
		fmt.Printf("Fist Name: %s\n", item.Data.Content.FirstName)
		fmt.Printf("Middle Name: %s\n", item.Data.Content.MiddleName)
		fmt.Printf("Last Name: %s\n", item.Data.Content.LastName)
		fmt.Printf("Birthdate: %s\n", item.Data.Content.Birthdate)
		fmt.Printf("Gender: %s\n", item.Data.Content.Gender)

		fmt.Println("")

		fmt.Printf("Address: %s\n", item.Data.Content.StreetAddress)
		fmt.Printf("Floor: %s\n", item.Data.Content.Floor)
		fmt.Printf("Postal Code: %s\n", item.Data.Content.ZipOrPostalCode)
		fmt.Printf("City: %s\n", item.Data.Content.City)
		fmt.Printf("State or Province: %s\n", item.Data.Content.StateOrProvince)
		fmt.Printf("County: %s\n", item.Data.Content.County)
		fmt.Printf("Country or Region: %s\n", item.Data.Content.CountryOrRegion)

		fmt.Println("")

		fmt.Printf("Phone Number: %s\n", item.Data.Content.PhoneNumber)
		fmt.Printf("Second Phone Number: %s\n", item.Data.Content.SecondPhoneNumber)
		fmt.Printf("Work Phone Number: %s\n", item.Data.Content.WorkPhoneNumber)
		fmt.Printf("Email: %s\n", item.Data.Content.Email)
		fmt.Printf("Work Email: %s\n", item.Data.Content.WorkEmail)
		fmt.Printf("Website: %s\n", item.Data.Content.Website)
		fmt.Printf("Personal Website: %s\n", item.Data.Content.PersonalWebsite)

		fmt.Println("")

		fmt.Printf("Social Security Number: %s\n", item.Data.Content.SocialSecurityNumber)
		fmt.Printf("Passport Number: %s\n", item.Data.Content.PassportNumber)
		fmt.Printf("License Number: %s\n", item.Data.Content.LicenseNumber)

		fmt.Println("")

		fmt.Printf("Company: %s\n", item.Data.Content.Company)
		fmt.Printf("Organization: %s\n", item.Data.Content.Organization)
		fmt.Printf("JobTitle: %s\n", item.Data.Content.JobTitle)

		fmt.Println("")

		fmt.Printf("X: %s\n", item.Data.Content.XHandle)
		fmt.Printf("Linkedin: %s\n", item.Data.Content.Linkedin)
		fmt.Printf("Reddit: %s\n", item.Data.Content.Reddit)
		fmt.Printf("Facebook: %s\n", item.Data.Content.Facebook)
		fmt.Printf("Yahoo: %s\n", item.Data.Content.Yahoo)
		fmt.Printf("Instagram: %s\n", item.Data.Content.Instagram)

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
func printIdentityList(items []*models.Item) error {
	fmt.Printf("Found %d ID(s)\n\n", len(items))

	for _, item := range items {
		fmt.Printf("%s (%s)\n",
			item.Data.Metadata.Name,
			item.ItemId,
		)
	}

	return nil
}
