package keygen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

// Algorithm represents an SSH key algorithm.
type Algorithm string

const (
	AlgorithmEd25519 Algorithm = "ed25519"
	AlgorithmRSA     Algorithm = "rsa"
)

// Options configure key generation.
type Options struct {
	Algorithm  Algorithm
	Comment    string
	Passphrase string
	Bits       int
}

// KeyPair holds generated key material in memory.
type KeyPair struct {
	PrivateKeyPEM []byte
	PublicKeySSH  string
	Algorithm     Algorithm
	Comment       string
}

// Generate creates a new SSH key pair according to the provided options.
// Passphrase support is currently reserved for future implementation.
func Generate(opts Options) (*KeyPair, error) {
	var (
		privateKey interface{}
		err        error
	)

	switch opts.Algorithm {
	case AlgorithmEd25519:
		privateKey, err = ed25519Generate()
	case AlgorithmRSA, "":
		bits := opts.Bits
		if bits == 0 {
			bits = 4096
		}
		privateKey, err = generateRSA(bits)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", opts.Algorithm)
	}
	if err != nil {
		return nil, fmt.Errorf("generate %s key: %w", opts.Algorithm, err)
	}

	privDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("marshal private key: %w", err)
	}

	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})

	pubKey, ok := privateKey.(interface{ Public() interface{} })
	if !ok {
		return nil, fmt.Errorf("private key does not implement Public()")
	}
	sshPub, err := ssh.NewPublicKey(pubKey.Public())
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}

	comment := opts.Comment
	if comment == "" {
		comment = "FyneSSH generated key"
	}

	authorized := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(sshPub)))
	if comment != "" {
		authorized = authorized + " " + comment
	}
	authorized = authorized + "\n"

	return &KeyPair{
		PrivateKeyPEM: privPEM,
		PublicKeySSH:  authorized,
		Algorithm:     opts.Algorithm,
		Comment:       comment,
	}, nil
}

// Save writes the private and public key files to disk with restrictive permissions.
func (k *KeyPair) Save(privatePath, publicPath string) error {
	if err := writeFile(privatePath, k.PrivateKeyPEM, 0600); err != nil {
		return fmt.Errorf("write private key: %w", err)
	}
	if err := writeFile(publicPath, []byte(k.PublicKeySSH), 0644); err != nil {
		return fmt.Errorf("write public key: %w", err)
	}
	return nil
}

func generateRSA(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

func writeFile(path string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(sshKeyDir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}

func sshKeyDir(_ string) string {
	// Placeholder: directory creation is handled by MkdirAll.
	// Future: extract and enforce ~/.ssh base dir permissions.
	return ""
}
