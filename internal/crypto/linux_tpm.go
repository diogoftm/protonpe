//go:build linux
// +build linux

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport"
	"github.com/google/go-tpm/tpm2/transport/linuxtpm"
	"golang.org/x/crypto/hkdf"
)

// Combines the TPM command interface with Close for resource cleanup.
type tpmCloser interface {
	transport.TPM
	io.Closer
}

// TPM NV slot where the ECC key lives after provisioning.
// This might be configurable in the future
const persistentHandle = tpm2.TPMHandle(0x81000100)

// P-256 ECDH child key template
var eccTemplate = tpm2.TPMTPublic{
	Type:    tpm2.TPMAlgECC,
	NameAlg: tpm2.TPMAlgSHA256,
	ObjectAttributes: tpm2.TPMAObject{
		FixedTPM:            true,
		FixedParent:         true,
		SensitiveDataOrigin: true,
		UserWithAuth:        true,
		Decrypt:             true,
	},
	Parameters: tpm2.NewTPMUPublicParms(tpm2.TPMAlgECC, &tpm2.TPMSECCParms{
		Scheme: tpm2.TPMTECCScheme{
			Scheme:  tpm2.TPMAlgECDH,
			Details: tpm2.NewTPMUAsymScheme(tpm2.TPMAlgECDH, &tpm2.TPMSKeySchemeECDH{HashAlg: tpm2.TPMAlgSHA256}),
		},
		CurveID: tpm2.TPMECCNistP256,
	}),
	Unique: tpm2.NewTPMUPublicID(tpm2.TPMAlgECC, &tpm2.TPMSECCPoint{}),
}

// Opens the TPM resource manager device, falling back to the raw device.
func openTPM() (tpmCloser, error) {
	tpm, err := linuxtpm.Open("/dev/tpmrm0")
	if err != nil {
		tpm, err = linuxtpm.Open("/dev/tpm0")
		if err != nil {
			return nil, fmt.Errorf("failed to open TPM (tried /dev/tpmrm0 and /dev/tpm0)")
		}
	}
	return tpm, nil
}

// IsTPMKeyProvisioned returns true if a key already exists at persistentHandle.
func IsTPMKeyProvisioned() (bool, error) {
	tpm, err := openTPM()
	if err != nil {
		return false, err
	}
	defer tpm.Close()

	_, err = tpm2.ReadPublic{ObjectHandle: persistentHandle}.Execute(tpm)
	return err == nil, nil
}

// Creates a P-256 ECDH key and persists it to the TPM's NV storage
func ProvisionTPMKey() error {
	tpm, err := openTPM()
	if err != nil {
		return err
	}
	defer tpm.Close()

	_, err = tpm2.ReadPublic{ObjectHandle: persistentHandle}.Execute(tpm)
	if err == nil {
		return fmt.Errorf("TPM key already provisioned at handle 0x%08X — evict it first", persistentHandle)
	}

	primaryResp, err := tpm2.CreatePrimary{
		PrimaryHandle: tpm2.TPMRHOwner,
		InPublic:      tpm2.New2B(tpm2.ECCSRKTemplate),
	}.Execute(tpm)
	if err != nil {
		return fmt.Errorf("failed to create TPM primary key")
	}
	defer tpm2.FlushContext{FlushHandle: primaryResp.ObjectHandle}.Execute(tpm)

	createResp, err := tpm2.Create{
		ParentHandle: tpm2.AuthHandle{
			Handle: primaryResp.ObjectHandle,
			Name:   primaryResp.Name,
			Auth:   tpm2.PasswordAuth(nil),
		},
		InPublic: tpm2.New2B(eccTemplate),
	}.Execute(tpm)
	if err != nil {
		return fmt.Errorf("failed to create ECC child key")
	}

	loadResp, err := tpm2.Load{
		ParentHandle: tpm2.AuthHandle{
			Handle: primaryResp.ObjectHandle,
			Name:   primaryResp.Name,
			Auth:   tpm2.PasswordAuth(nil),
		},
		InPrivate: createResp.OutPrivate,
		InPublic:  createResp.OutPublic,
	}.Execute(tpm)
	if err != nil {
		return fmt.Errorf("failed to load ECC child key")
	}
	defer tpm2.FlushContext{FlushHandle: loadResp.ObjectHandle}.Execute(tpm)

	_, err = tpm2.EvictControl{
		Auth: tpm2.AuthHandle{
			Handle: tpm2.TPMRHOwner,
			Auth:   tpm2.PasswordAuth(nil),
		},
		ObjectHandle: tpm2.NamedHandle{
			Handle: loadResp.ObjectHandle,
			Name:   loadResp.Name,
		},
		PersistentHandle: persistentHandle,
	}.Execute(tpm)
	if err != nil {
		return fmt.Errorf("failed to persist ECC key to TPM NV: %w", err)
	}

	return nil
}

// Removes the persistent key from TPM NV storage.
// Use this to rotate or clean up the persistent key.
func EvictTPMKey() error {
	tpm, err := openTPM()
	if err != nil {
		return err
	}
	defer tpm.Close()

	_, err = tpm2.EvictControl{
		Auth:             tpm2.AuthHandle{Handle: tpm2.TPMRHOwner},
		ObjectHandle:     persistentHandle,
		PersistentHandle: persistentHandle,
	}.Execute(tpm)
	if err != nil {
		return fmt.Errorf("failed to evict TPM key: %w", err)
	}
	return nil
}

// Encrypts inputPath to outputPath using the persistent P-256 TPM key.
//
// Output file layout — every field is uint32 big-endian length-prefixed:
//
//	[4]    magic tpmMagicECC
//	[4+32] ephemeral pub X  (P-256 coordinate, 32 bytes)
//	[4+32] ephemeral pub Y
//	[4+12] AES-GCM nonce
//	[rest] AES-GCM ciphertext + 16-byte authentication tag
func encryptFileTPMECC(inputPath string, outputPath string) error {
	tpm, err := openTPM()
	if err != nil {
		return err
	}
	defer tpm.Close()

	pubResp, err := tpm2.ReadPublic{ObjectHandle: persistentHandle}.Execute(tpm)
	if err != nil {
		fmt.Printf("TPM key not provisioned at handle 0x%08X. Trying to create it...\n", persistentHandle)
		err = ProvisionTPMKey()
		pubResp, err = tpm2.ReadPublic{ObjectHandle: persistentHandle}.Execute(tpm)
		if err != nil {
			return fmt.Errorf("Error creating TPM parent key: %s", err)
		}
	}

	tpmPub, err := pubResp.OutPublic.Contents()
	if err != nil {
		return fmt.Errorf("failed to parse TPM public area")
	}
	tpmECCPub, err := tpmPub.Unique.ECC()
	if err != nil {
		return fmt.Errorf("failed to extract ECC point")
	}

	tpmPubBytes := make([]byte, 65)
	tpmPubBytes[0] = 0x04
	copy(tpmPubBytes[1:33], pad32(tpmECCPub.X.Buffer))
	copy(tpmPubBytes[33:65], pad32(tpmECCPub.Y.Buffer))
	tpmECDHPub, err := ecdh.P256().NewPublicKey(tpmPubBytes)
	if err != nil {
		return fmt.Errorf("failed to import TPM public key")
	}

	ephemeralKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate ephemeral ECC key")
	}
	ephemeralPub := ephemeralKey.PublicKey()

	sharedSecret, err := ephemeralKey.ECDH(tpmECDHPub)
	if err != nil {
		return fmt.Errorf("ECDH has failed")
	}

	aesKey := make([]byte, 32)
	hkdfReader := hkdf.New(sha256.New, sharedSecret, nil, []byte("tpm-ecc-file-encryption"))
	if _, err := io.ReadFull(hkdfReader, aesKey); err != nil {
		return fmt.Errorf("HKDF has failed")
	}

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read file")
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return fmt.Errorf("failed to create AES cipher")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to generate nonce")
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	rawPub := ephemeralPub.Bytes()
	ephX := rawPub[1:33]
	ephY := rawPub[33:65]

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create encrypted output file")
	}
	defer f.Close()

	if _, err = f.Write(tpmMagicECC); err != nil {
		return fmt.Errorf("I/O error")
	}
	for _, field := range [][]byte{ephX, ephY, nonce} {
		if err = writeField(f, field); err != nil {
			return fmt.Errorf("I/O error")
		}
	}
	if _, err = f.Write(ciphertext); err != nil {
		return fmt.Errorf("I/O error")
	}

	return nil
}

// Decrypts files encrypted by encryptFileTPMECC.
//
// File layout after the 4-byte magic:
//
//	[4+32] ephemeral pub X (P-256 coordinate)
//	[4+32] ephemeral pub Y
//	[4+12] AES-GCM nonce
//	[rest] AES-GCM ciphertext + 16-byte tag
func decryptTPMECC(r io.Reader) ([]byte, error) {
	magic := make([]byte, len(tpmMagicECC))
	if _, err := io.ReadFull(r, magic); err != nil {
		return nil, fmt.Errorf("reading magic: %w", err)
	}

	ephX, err := readField(r)
	if err != nil {
		return nil, fmt.Errorf("reading ephemeral X: %w", err)
	}
	ephY, err := readField(r)
	if err != nil {
		return nil, fmt.Errorf("reading ephemeral Y: %w", err)
	}
	nonce, err := readField(r)
	if err != nil {
		return nil, fmt.Errorf("reading nonce: %w", err)
	}
	ciphertext, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading ciphertext: %w", err)
	}

	tpm, err := openTPM()
	if err != nil {
		return nil, err
	}
	defer tpm.Close()

	pubResp, err := tpm2.ReadPublic{ObjectHandle: persistentHandle}.Execute(tpm)
	if err != nil {
		return nil, fmt.Errorf("ReadPublic (key not provisioned at 0x%08X?): %w", persistentHandle, err)
	}
	keyPub, err := pubResp.OutPublic.Contents()
	keyName, err := tpm2.ObjectName(keyPub)

	zGenResp, err := tpm2.ECDHZGen{
		KeyHandle: tpm2.AuthHandle{
			Handle: persistentHandle,
			Name:   *keyName,
			Auth:   tpm2.PasswordAuth(nil),
		},
		InPoint: tpm2.New2B(tpm2.TPMSECCPoint{
			X: tpm2.TPM2BECCParameter{Buffer: ephX},
			Y: tpm2.TPM2BECCParameter{Buffer: ephY},
		}),
	}.Execute(tpm)
	if err != nil {
		return nil, fmt.Errorf("ECDHZGen: %w", err)
	}

	sharedSecretTmp, _ := zGenResp.OutPoint.Contents()
	sharedSecret := sharedSecretTmp.X.Buffer

	aesKey := make([]byte, 32)
	hkdfReader := hkdf.New(sha256.New, sharedSecret, nil, []byte("tpm-ecc-file-encryption"))
	if _, err := io.ReadFull(hkdfReader, aesKey); err != nil {
		return nil, fmt.Errorf("HKDF: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}
	if len(nonce) != gcm.NonceSize() {
		return nil, fmt.Errorf("nonce length mismatch: got %d, want %d", len(nonce), gcm.NonceSize())
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("AES-GCM open (wrong TPM, corrupted, or tampered): %w", err)
	}
	return plaintext, nil
}

// Reads a uint32 big-endian length-prefixed byte slice.
func readField(r io.Reader) ([]byte, error) {
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lenBuf); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint32(lenBuf)
	if n > 64*1024*1024 {
		return nil, fmt.Errorf("field length %d exceeds sanity limit (64 MiB)", n)
	}
	data := make([]byte, n)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	return data, nil
}
