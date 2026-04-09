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

var tpmMagicECC = []byte("TPM\x02")

// Decrypt opens filePath, auto-detects the format, and returns the vault contents.
func Decrypt(filePath string) (models.VaultFile, error) {
	var result models.VaultFile

	f, err := os.Open(filePath)
	if err != nil {
		return result, fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	// Peek at the first 32 bytes to identify the format, then rewind.
	buf := make([]byte, 32)
	n, _ := f.Read(buf)
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return result, fmt.Errorf("error seeking file: %w", err)
	}
	header := buf[:n]

	var data []byte

	var tpmEnabled = false

	if bytes.HasPrefix(header, tpmMagicECC) {
		data, err = decryptTPMECC(f)
		if err != nil {
			return result, fmt.Errorf("TPM decryption failed: %w", err)
		}
		tpmEnabled = true
	}

	if isPGP(header) || isPGP(data) {
		if tpmEnabled {
			data, err = decryptPGPBytes(data)
		} else {
			data, err = decryptPGP(f)
		}

		if err != nil {
			return result, fmt.Errorf("PGP decryption failed: %w", err)
		}
	} else {
		data, err = io.ReadAll(f)
		if err != nil {
			return result, fmt.Errorf("error reading file: %w", err)
		}
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("error parsing JSON: %w", err)
	}
	utils.Zero(data)
	return result, nil
}

// isPGP returns true for both ASCII-armoured and binary OpenPGP data.
func isPGP(header []byte) bool {
	if bytes.HasPrefix(header, []byte("-----BEGIN PGP")) {
		return true
	}
	return len(header) > 0 && (header[0]&0x80) != 0
}

// Handles the decryption of password-protected PGP files
func decryptPGP(r io.ReadSeeker) ([]byte, error) {
	reader, err := utils.DetectInput(r)
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	fmt.Fprint(os.Stderr, "Password: ")
	pass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("error reading password: %w", err)
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
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, md.UnverifiedBody); err != nil {
		return nil, fmt.Errorf("error reading decrypted data: %w", err)
	}
	return buf.Bytes(), nil
}

// Handles the decryption of password-protected PGP file
// received as bytes.
func decryptPGPBytes(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)

	reader, err := utils.DetectInput(r)
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	fmt.Fprint(os.Stderr, "Password: ")
	pass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("error reading password: %w", err)
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
			tries++
			return pass, nil
		},
		nil,
	)
	utils.Zero(pass)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, md.UnverifiedBody); err != nil {
		return nil, fmt.Errorf("error reading decrypted data: %w", err)
	}

	return buf.Bytes(), nil
}

// Gets the source file from a flag or env var, decrypts it,
// and unmarshals it to the internal struct.
func GetFileContent(c *cli.Command) (models.VaultFile, error) {
	filePath, err := utils.GetFilePath(c)
	if err != nil {
		return models.VaultFile{}, err
	}
	return Decrypt(filePath)
}
