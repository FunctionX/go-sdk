package crypto

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/pkg/errors"
)

type newPrivKeyFunc func(*ecdsa.PrivateKey) PrivKey

var algoList map[Algorithm]newPrivKeyFunc

func RegisterAlgo(name Algorithm, f newPrivKeyFunc) {
	if algoList == nil {
		algoList = make(map[Algorithm]newPrivKeyFunc)
	}
	algoList[name] = f
}

func NewPrivKeyFromKeyStore(keystoreFile, passwordFile string, algo Algorithm, needPass bool) (PrivKey, error) {
	keyJson, err := os.ReadFile(keystoreFile)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read the keyfile at %s", keystoreFile))
	}
	passphrase := ""
	if needPass {
		passphrase, err = getPassphrase(passwordFile)
		if err != nil {
			return nil, err
		}
	}
	key, err := keystore.DecryptKey(keyJson, passphrase)
	if err != nil {
		return nil, errors.Wrap(err, "error decrypting key")
	}
	newPrivKey, ok := algoList[algo]
	if !ok {
		return nil, errors.New("invalid algorithm")
	}
	return newPrivKey(key.PrivateKey), nil
}

func getPassphrase(passwordFile string) (string, error) {
	if passwordFile == "" {
		return utils.GetPassPhrase("", false), nil
	}
	password, err := os.ReadFile(passwordFile)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to read the keyfile at %s", passwordFile))
	}
	return strings.TrimRight(string(password), "\r\n"), nil
}
