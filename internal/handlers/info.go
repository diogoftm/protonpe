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

// `info` command handler
func Info(ctx context.Context, c *cli.Command) error {
	fileContent, err := crypto.GetFileContent(c)
	if err != nil {
		return err
	}
	defer utils.ZeroVaultFile(&fileContent)

	return printInfo(fileContent)
}

//
// Outputs
//

// Print general information about the export and available vaults
func printInfo(fileContent models.VaultFile) error {
	// Version
	fmt.Printf("Version: %s\n", fileContent.Version)

	// Vault count
	fmt.Printf("Vaults: %d\n\n", len(fileContent.Vaults))

	// Vault details
	for id, vault := range fileContent.Vaults {
		fmt.Printf("• %s\n", vault.Name)
		fmt.Printf("  ID: %s\n", id)
		fmt.Printf("  Description: %s\n", vault.Description)
		fmt.Printf("  Items: %d\n\n", len(vault.Items))
	}

	return nil
}
