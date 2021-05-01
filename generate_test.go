// +build dev test

package csrf

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var TestItems = []struct {
	comment     string
	token       string
	salt        string
	envValue    string
	expectHash  string
	expectError error
}{
	{
		comment:     "no env salts set",
		envValue:    ``,
		token:       "thetoken",
		salt:        "",
		expectHash:  "",
		expectError: errors.New("csrf.Generate(): No salts found in ENV"),
	},
	{
		comment:     "one salt set",
		envValue:    `one`,
		token:       "thetoken",
		salt:        "thesalt",
		expectHash:  "US3YTYWkOlwMW8VwYKNROYWdN4esmYGvsRhd/YGigjM=",
		expectError: nil,
	},
	{
		comment:     "two salts set",
		envValue:    `one,two`,
		token:       "thetoken",
		salt:        "thesalt",
		expectHash:  "US3YTYWkOlwMW8VwYKNROYWdN4esmYGvsRhd/YGigjM=",
		expectError: nil,
	},
	{
		comment:     "two salts set, reversed",
		envValue:    `two,one`,
		token:       "thetoken",
		salt:        "thesalt",
		expectHash:  "US3YTYWkOlwMW8VwYKNROYWdN4esmYGvsRhd/YGigjM=",
		expectError: nil,
	},
}

func TestValidate(t *testing.T) {

	os.Setenv(SaltsEnvKey, "thesalt")

	token, hash, err := Generate()
	assert.Equal(t, nil, err, "test validate setup")

	validity, saltIndex := Validate(token, hash)
	assert.Equal(t, true, validity, "test validate true")
	assert.Equal(t, 0, saltIndex, "test validate saltIndex 0")

	os.Setenv(SaltsEnvKey, "thesalt")
	validity, saltIndex = Validate(token[1:], hash)
	assert.Equal(t, false, validity, "test validate wrong token")

	os.Setenv(SaltsEnvKey, "thesalt")
	validity, saltIndex = Validate(token, hash[1:])
	assert.Equal(t, false, validity, "test validate wrong hash")

	os.Setenv(SaltsEnvKey, "thesalt2")
	validity, saltIndex = Validate(token, hash)
	assert.Equal(t, false, validity, "test validate wrong salt")

	os.Unsetenv(SaltsEnvKey)
	validity, saltIndex = Validate(token, hash)
	assert.Equal(t, false, validity, "test validate not salt env var")

	// test that we can validate with an older salt in the list
	os.Setenv(SaltsEnvKey, "thesalt")
	token, hash, err = Generate()
	os.Setenv(SaltsEnvKey, "thesalt2,thesalt")
	validity, saltIndex = Validate(token, hash)
	assert.Equal(t, true, validity, "test validate true multiple salts")
	assert.Equal(t, 1, saltIndex, "test validate saltIndex multiple salts")

}

func TestPublicGenerate(t *testing.T) {

	for _, item := range TestItems {
		os.Setenv(SaltsEnvKey, item.envValue)

		token, _, err := Generate()

		assert.Equal(t, item.expectError, err, item.comment)
		assert.Equal(t, TokenLen, len(token), item.comment)
	}
}

func TestPrivateGenerate(t *testing.T) {

	for _, item := range TestItems {
		os.Setenv(SaltsEnvKey, item.envValue)
		_, hash, err := generate(item.token, item.salt)
		assert.Equal(t, item.expectHash, hash, item.comment)
		assert.Equal(t, item.expectError, err, item.comment)

	}
}

func TestGetSaltList(t *testing.T) {

	os.Unsetenv(SaltsEnvKey)
	salts := getSaltList()
	assert.Equal(t, []string{}, salts, "unset salts env var")

	os.Setenv(SaltsEnvKey, ``)
	salts = getSaltList()
	assert.Equal(t, []string{}, salts, "empty salts env var")

	os.Setenv(SaltsEnvKey, `one`)
	salts = getSaltList()
	assert.Equal(t, []string{"one"}, salts, "one salt")

	os.Setenv(SaltsEnvKey, `one,two`)
	salts = getSaltList()
	assert.Equal(t, []string{"one", "two"}, salts, "two salts")

	os.Setenv(SaltsEnvKey, `two`)
	salts = getSaltList()
	assert.Equal(t, []string{"two"}, salts, "remove first salt")

}
