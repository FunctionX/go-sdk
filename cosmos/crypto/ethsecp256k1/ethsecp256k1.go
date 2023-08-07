package ethsecp256k1

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/subtle"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"

	cosmoscrypto "github.com/functionx/go-sdk/cosmos/crypto"
)

const (
	// KeyType is the string constant for the Secp256k1 algorithm
	KeyType cosmoscrypto.Algorithm = "eth_secp256k1"
)

func init() {
	cosmoscrypto.RegisterAlgo(KeyType, NewPrivKey)
}

// ----------------------------------------------------------------------------
// secp256k1 Private Key

var _ cosmoscrypto.PrivKey = &PrivKey{}

func NewPrivKey(privateKey *ecdsa.PrivateKey) cosmoscrypto.PrivKey {
	return &PrivKey{Key: crypto.FromECDSA(privateKey)}
}

// Bytes returns the byte representation of the ECDSA Private Key.
func (privKey *PrivKey) Bytes() []byte {
	bz := make([]byte, len(privKey.Key))
	copy(bz, privKey.Key)

	return bz
}

// PubKey returns the ECDSA private key's public key. If the privkey is not valid
// it returns a nil value.
func (privKey *PrivKey) PubKey() cosmoscrypto.PubKey {
	ecdsaPrivKey, err := privKey.ToECDSA()
	if err != nil {
		return nil
	}

	return &PubKey{
		Key: crypto.CompressPubkey(&ecdsaPrivKey.PublicKey),
	}
}

// Equals returns true if two ECDSA private keys are equal and false otherwise.
func (privKey *PrivKey) Equals(other cosmoscrypto.PrivKey) bool {
	return privKey.Type() == other.Type() && subtle.ConstantTimeCompare(privKey.Bytes(), other.Bytes()) == 1
}

// Type returns eth_secp256k1
func (privKey *PrivKey) Type() cosmoscrypto.Algorithm {
	return KeyType
}

// Sign creates a recoverable ECDSA signature on the secp256k1 curve over the
// provided hash of the message. The produced signature is 65 bytes
// where the last byte contains the recovery ID.
func (privKey *PrivKey) Sign(digestBz []byte) ([]byte, error) {
	if len(digestBz) != crypto.DigestLength {
		digestBz = crypto.Keccak256Hash(digestBz).Bytes()
	}

	key, err := privKey.ToECDSA()
	if err != nil {
		return nil, err
	}

	return crypto.Sign(digestBz, key)
}

// ToECDSA returns the ECDSA private key as a reference to ecdsa.PrivateKey type.
func (privKey *PrivKey) ToECDSA() (*ecdsa.PrivateKey, error) {
	return crypto.ToECDSA(privKey.Bytes())
}

// ----------------------------------------------------------------------------
// secp256k1 Public Key

var _ cosmoscrypto.PubKey = &PubKey{}

// Address returns the address of the ECDSA public key.
// The function will return an empty address if the public key is invalid.
func (pubKey *PubKey) Address() cosmoscrypto.Address {
	pubk, err := crypto.DecompressPubkey(pubKey.Key)
	if err != nil {
		return nil
	}

	return crypto.PubkeyToAddress(*pubk).Bytes()
}

// Bytes returns the raw bytes of the ECDSA public key.
func (pubKey *PubKey) Bytes() []byte {
	bz := make([]byte, len(pubKey.Key))
	copy(bz, pubKey.Key)

	return bz
}

// String implements the fmt.Stringer interface.
func (pubKey *PubKey) String() string {
	return fmt.Sprintf("EthPubKeySecp256k1{%X}", pubKey.Key)
}

// Type returns eth_secp256k1
func (pubKey *PubKey) Type() cosmoscrypto.Algorithm {
	return KeyType
}

// Equals returns true if the pubkey type is the same and their bytes are deeply equal.
func (pubKey *PubKey) Equals(other cosmoscrypto.PubKey) bool {
	return pubKey.Type() == other.Type() && bytes.Equal(pubKey.Bytes(), other.Bytes())
}

// VerifySignature verifies that the ECDSA public key created a given signature over
// the provided message. It will calculate the Keccak256 hash of the message
// prior to verification and approve verification if the signature can be verified
// from either the original message or its EIP-712 representation.
//
// CONTRACT: The signature should be in [R || S] format.
func (pubKey *PubKey) VerifySignature(msg, sig []byte) bool {
	return pubKey.verifySignatureECDSA(msg, sig) /*|| pubKey.verifySignatureAsEIP712(msg, sig)*/
}

// Verifies the signature as an EIP-712 signature by first converting the message payload
// to EIP-712 object bytes, then performing ECDSA verification on the hash. This is to support
// signing a Cosmos payload using EIP-712.
// func (pubKey *PubKey) verifySignatureAsEIP712(msg, sig []byte) bool {
// 	eip712Bytes, err := eip712.GetEIP712BytesForMsg(msg)
// 	if err != nil {
// 		return false
// 	}
//
// 	if pubKey.verifySignatureECDSA(eip712Bytes, sig) {
// 		return true
// 	}
//
// 	// Try verifying the signature using the legacy EIP-712 encoding
// 	legacyEIP712Bytes, err := eip712.LegacyGetEIP712BytesForMsg(msg)
// 	if err != nil {
// 		return false
// 	}
//
// 	return pubKey.verifySignatureECDSA(legacyEIP712Bytes, sig)
// }

// Perform standard ECDSA signature verification for the given raw bytes and signature.
func (pubKey *PubKey) verifySignatureECDSA(msg, sig []byte) bool {
	if len(sig) == crypto.SignatureLength {
		// remove recovery ID (V) if contained in the signature
		sig = sig[:len(sig)-1]
	}

	// the signature needs to be in [R || S] format when provided to VerifySignature
	return crypto.VerifySignature(pubKey.Key, crypto.Keccak256Hash(msg).Bytes(), sig)
}
