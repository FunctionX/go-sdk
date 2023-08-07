package secp256k1

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"

	secp256k1 "github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ripemd160" // nolint: staticcheck // necessary for Bitcoin address format, keep around for backwards compatibility

	cosmoscrypto "github.com/functionx/go-sdk/cosmos/crypto"
)

const KeyType cosmoscrypto.Algorithm = "secp256k1"

func init() {
	cosmoscrypto.RegisterAlgo(KeyType, NewPrivKey)
}

// Bytes returns the byte representation of the Private Key.
func (privKey *PrivKey) Bytes() []byte {
	return privKey.Key
}

// PubKey performs the point-scalar multiplication from the privKey on the
// generator point to get the pubkey.
func (privKey *PrivKey) PubKey() cosmoscrypto.PubKey {
	_, pubkeyObject := secp256k1.PrivKeyFromBytes(privKey.Key)
	pk := pubkeyObject.SerializeCompressed()
	return &PubKey{Key: pk}
}

// Equals - you probably don't need to use this.
// Runs in constant time based on length of the
func (privKey *PrivKey) Equals(other cosmoscrypto.PrivKey) bool {
	return privKey.Type() == other.Type() && subtle.ConstantTimeCompare(privKey.Bytes(), other.Bytes()) == 1
}

func (privKey *PrivKey) Type() cosmoscrypto.Algorithm {
	return KeyType
}

func NewPrivKey(privateKey *ecdsa.PrivateKey) cosmoscrypto.PrivKey {
	return &PrivKey{Key: crypto.FromECDSA(privateKey)}
}

// PubKeySize is comprised of 32 bytes for one field element
// (the x-coordinate), plus one byte for the parity of the y-coordinate.
const PubKeySize = 33

// Address returns a Bitcoin style addresses: RIPEMD160(SHA256(pubkey))
func (pubKey *PubKey) Address() cosmoscrypto.Address {
	if len(pubKey.Key) != PubKeySize {
		panic("length of pubkey is incorrect")
	}

	sha := sha256.Sum256(pubKey.Key)
	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha[:]) // does not error
	return hasherRIPEMD160.Sum(nil)
}

// Bytes returns the pubkey byte format.
func (pubKey *PubKey) Bytes() []byte {
	return pubKey.Key
}

func (pubKey *PubKey) String() string {
	return fmt.Sprintf("PubKeySecp256k1{%X}", pubKey.Key)
}

func (pubKey *PubKey) Type() cosmoscrypto.Algorithm {
	return KeyType
}

func (pubKey *PubKey) Equals(other cosmoscrypto.PubKey) bool {
	return pubKey.Type() == other.Type() && bytes.Equal(pubKey.Bytes(), other.Bytes())
}

func Sha256(bytes []byte) []byte {
	hasher := sha256.New()
	hasher.Write(bytes)
	return hasher.Sum(nil)
}
