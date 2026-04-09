package crypto

import (
	"encoding/binary"
	"fmt"
	"io"
	"runtime"
)

// Supported encryption methods
var ENCRYPTION_METHODS = []string{"TPM_ECCP256_AES256GCM_SHA256"}

// Gets the source file from a flag or env var, decrypts it,
// and unmarshals it to the internal struct.
func EncryptFileContent(inFilePath string, outFilePath string, method string) error {
	switch method {
	case ENCRYPTION_METHODS[0]:
		if runtime.GOOS == "linux" {
			return encryptFileTPMECC(inFilePath, outFilePath)
		} else {
			return fmt.Errorf("%s currently only supported on Linux", method)
		}
	default:
		return fmt.Errorf("Invalid encryption method")
	}
}

// pad32 zero-pads b on the left to exactly 32 bytes (P-256 coordinate size).
func pad32(b []byte) []byte {
	out := make([]byte, 32)
	copy(out[32-len(b):], b)
	return out
}

// writeField writes data preceded by its uint32 big-endian length.
func writeField(w io.Writer, data []byte) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(len(data)))
	if _, err := w.Write(buf); err != nil {
		return err
	}
	_, err := w.Write(data)
	return err
}
