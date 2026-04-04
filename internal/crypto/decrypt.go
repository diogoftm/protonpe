package crypto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"

	"github.com/diogoftm/protonpe/internal/models"
	"github.com/diogoftm/protonpe/internal/utils"
)

// Decrypt PGP file (if needed) and unmarshal it to internal struct
func Decrypt(filePath string) (models.VaultFile, error) {
	var result models.VaultFile
	f, err := os.Open(filePath)
	if err != nil {
		return result, fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	buf := make([]byte, 32)
	n, _ := f.Read(buf)
	f.Seek(0, 0)

	isPGP := false

	if bytes.HasPrefix(buf[:n], []byte("-----BEGIN PGP")) {
		isPGP = true
	}

	if !isPGP && n > 0 && (buf[0]&0x80) != 0 {
		isPGP = true
	}

	var data []byte

	if isPGP {
		reader, err := utils.DetectInput(f)
		if err != nil {
			return result, fmt.Errorf("error reading input: %w", err)
		}

		fmt.Fprint(os.Stderr, "Password: ")
		pass, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return result, fmt.Errorf("error reading password: %w", err)
		}

		tries := 0
		md, err := openpgp.ReadMessage(
			reader,
			nil,
			func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
				if tries == 1 {
					utils.Zero(pass)
					return nil, fmt.Errorf("wrong password")
				}
				tries = 1
				return pass, nil
			},
			nil,
		)
		utils.Zero(pass)

		if err != nil {
			return result, fmt.Errorf("decryption failed: %w", err)
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, md.UnverifiedBody); err != nil {
			return result, fmt.Errorf("error reading decrypted data: %w", err)
		}

		data = buf.Bytes()

	} else {
		raw, err := io.ReadAll(f)
		if err != nil {
			return result, fmt.Errorf("error reading file: %w", err)
		}
		data = raw
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("error parsing JSON: %w", err)
	}

	utils.Zero(data)

	return result, nil
}

// Get source file from flag or env var, decrypt it (if needed),
// and unmarshal it to internal struct
func GetFileContent(c *cli.Command) (models.VaultFile, error) {
	filePath, err := utils.GetFilePath(c)
	if err != nil {
		return models.VaultFile{}, err
	}

	return Decrypt(filePath)
}
