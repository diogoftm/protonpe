//go:build windows
// +build windows

package crypto

import (
	"fmt"
	"io"
)

func IsTPMKeyProvisioned() (bool, error) {
	return false, nil
}

func ProvisionTPMKey() error {
	return fmt.Errorf("TPM not supported on Windows")
}

func EvictTPMKey() error {
	return fmt.Errorf("TPM not supported on Windows")
}

func encryptFileTPMECC(inputPath string, outputPath string) error {
	return fmt.Errorf("TPM not supported on Windows")
}

func decryptTPMECC(r io.Reader) ([]byte, error) {
	return nil, fmt.Errorf("TPM not supported on Windows")
}
