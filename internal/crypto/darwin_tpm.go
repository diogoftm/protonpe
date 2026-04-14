//go:build darwin
// +build darwin

package crypto

import (
	"fmt"
	"io"
)

func IsTPMKeyProvisioned() (bool, error) {
	return false, nil
}

func ProvisionTPMKey() error {
	return fmt.Errorf("TPM not supported on macOS")
}

func EvictTPMKey() error {
	return fmt.Errorf("TPM not supported on macOS")
}

func encryptFileTPMECC(inputPath string, outputPath string) error {
	return fmt.Errorf("TPM not supported on macOS")
}

func decryptTPMECC(r io.Reader) ([]byte, error) {
	return nil, fmt.Errorf("TPM not supported on macOS")
}
