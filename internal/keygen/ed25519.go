package keygen

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
)

// keypair wraps an ed25519 key pair and implements the interfaces needed
// by internal key generation logic.
type keypair struct {
	private ed25519.PrivateKey
	public  ed25519.PublicKey
}

func ed25519Generate() (*keypair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate ed25519 key: %w", err)
	}
	return &keypair{private: priv, public: pub}, nil
}

func (k *keypair) Public() interface{} {
	return k.public
}
