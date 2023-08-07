package hd

import (
	"errors"

	"github.com/tyler-smith/go-bip39"
)

// Derive derives and returns the secp256k1 private key for the given seed and HD path.
func Derive(mnemonic, bip39Passphrase, hdPath string) ([]byte, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalid mnemonic")
	}
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
	if err != nil {
		return nil, err
	}

	masterPriv, ch := ComputeMastersFromSeed(seed)
	if len(hdPath) == 0 {
		return masterPriv[:], nil
	}
	derivedKey, err := DerivePrivateKeyForPath(masterPriv, ch, hdPath)

	return derivedKey, err
}
