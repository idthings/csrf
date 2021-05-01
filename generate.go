package csrf

import (
	"encoding/base64"
	"errors"
	"github.com/idthings/alphanum"
	"golang.org/x/crypto/scrypt"
	"log"
	"os"
	"strings"
)

func Validate(inputToken string, inputHash string) (bool, int) {

	salts := getSaltList()
	if len(salts) < 1 {
		return false, -1
	}

	for idx, salt := range salts {
		_, computedHash, err := generate(inputToken, salt)
		if err != nil {
			return false, -1
		}
		if computedHash == inputHash {
			return true, idx
		}
	}
	return false, -1
}

func Generate() (string, string, error) {

	token := alphanum.New(TokenLen)

	salts := getSaltList()
	if len(salts) < 1 || salts[0] == "" {
		return token, "", errors.New("csrf.Generate(): No salts found in ENV")
	}

	return generate(token, salts[0])
}

func generate(token string, salt string) (string, string, error) {

	if salt == "" {
		return token, "", errors.New("csrf.Generate(): No salts found in ENV")
	}

	hash := ""

	dk, err := scrypt.Key([]byte(token), []byte(salt), N, R, P, KeyLen)
	if err != nil {
		log.Fatal(err)
	}
	hash = base64.StdEncoding.EncodeToString(dk)

	return token, hash, nil
}

func getSaltList() []string {
	value := os.Getenv(SaltsEnvKey)
	salts := strings.Split(value, ",")
	if salts[0] == "" {
		return []string{} // explicitly eturn a empty list, getenv always returns at least ""
	}
	return salts
}
