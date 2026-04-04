package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/urfave/cli/v3"
	"golang.design/x/clipboard"

	"github.com/diogoftm/protonpe/internal/models"
)

// Detects ASCII-armored vs binary PGP
func DetectInput(r io.Reader) (io.Reader, error) {
	buf := bufio.NewReader(r)

	peek, err := buf.Peek(30)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if bytes.HasPrefix(peek, []byte("-----BEGIN PGP")) {
		block, err := armor.Decode(buf)
		if err != nil {
			return nil, fmt.Errorf("armor decode failed: %w", err)
		}
		return block.Body, nil
	}

	return buf, nil
}

// Overwrites a byte slice
func Zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// Detect if stdout is redirected
func IsRedirected(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) == 0
}

// Get location of the Proton Pass exported file
func GetFilePath(c *cli.Command) (string, error) {
	var filePath string
	if c.String("file") == "" {
		var found bool
		filePath, found = os.LookupEnv("PROTONPE_FILE")
		if !found {
			return filePath, fmt.Errorf("no file selected: neither environment variable PROTONPE_FILE or --file flag were set")
		}
	} else {
		filePath = c.String("file")
	}
	return filePath, nil
}

// Get a single item from the parsed file matching provided criteria
func GetItem(fileContent models.VaultFile, vaultId string, itemId string, itemType string) *models.Item {
	for id, vault := range fileContent.Vaults {
		if vaultId == "" || id == vaultId || vault.Name == vaultId {
			for i := range vault.Items {
				item := &vault.Items[i]
				if item.Data.Type == itemType && (item.Data.Metadata.Name == itemId || item.ItemId == itemId) {
					return item
				}
			}
		}
	}
	return nil
}

// Get all items from the parsed file matching provided criteria
func GetItems(fileContent models.VaultFile, vaultId string, itemId string, itemType string) []*models.Item {
	var items []*models.Item

	for id, vault := range fileContent.Vaults {
		// If no vault specified, or matches vault ID/Name
		if vaultId == "" || id == vaultId || vault.Name == vaultId {
			for i := range vault.Items {
				item := &vault.Items[i]
				if (itemType == "" || item.Data.Type == itemType) && (item.Data.Metadata.Name == itemId || item.ItemId == itemId) {
					items = append(items, item)
				}
			}
		}
	}

	return items
}

// Get all items from the parsed file matching provided criteria (with no item type restriction)
func GetAllItems(fileContent models.VaultFile, vaultId string, itemType string) []*models.Item {
	var items []*models.Item

	for id, vault := range fileContent.Vaults {
		if vaultId == "" || id == vaultId || vault.Name == vaultId {
			for i := range vault.Items {
				item := &vault.Items[i]
				if item.Data.Type == itemType {
					items = append(items, item)
				}
			}
		}
	}

	return items
}

// Best effort to clear some of the most sensitive data inside the VaultFile struct
func ZeroVaultFile(v *models.VaultFile) {
	for vid := range v.Vaults {
		for i := range v.Vaults[vid].Items {
			item := &v.Vaults[vid].Items[i]

			item.Data.Content.Password = ""
			item.Data.Metadata.Note = ""
			item.Data.Content.TotpUri = ""
			item.Data.Content.ItemEmail = ""
			item.Data.Content.ItemUsername = ""
		}
	}

	v.Vaults = nil
	v.Version = ""
	v.UserId = ""
}

// Write to clipboard and clean it after `ttl`
func WriteToClipboard(text string, ttl int) error {
	if ttl < 0 || ttl > 30 {
		return fmt.Errorf("invalid TTL value (0 <= ttl <= 30)")
	}

	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("failed to init clipboard: %w", err)
	}

	clipboard.Write(clipboard.FmtText, []byte(text))
	time.Sleep(10 * time.Millisecond)

	if ttl != 0 {
		for i := ttl; i > 0; i-- {
			fmt.Printf("\rClearing in %ds...", i)
			time.Sleep(1 * time.Second)
		}
		clipboard.Write(clipboard.FmtText, []byte{})
		time.Sleep((10 * time.Millisecond))
		fmt.Println()
	}

	return nil
}
