package internal

import (
	"strings"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/md5_crypt"
	_ "github.com/GehirnInc/crypt/sha256_crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
	"github.com/openwall/yescrypt-go"
)

func MakeHash(plaintext string, algorithm string, salt string, parameters string) (string, error) {
	var err error
	var hash []byte

	switch algorithm {
	case "yescrypt":
		saltBytes := []byte("$y$" + parameters + "$" + salt)
		hash, err = yescrypt.Hash([]byte(plaintext), saltBytes)
	case "md5":
		crypt := crypt.MD5.New()
		var hashStr string
		hashStr, err = crypt.Generate([]byte(plaintext), []byte("$1$"+salt))
		hash = []byte(hashStr)
	case "sha-256":
		crypt := crypt.SHA256.New()
		var hashStr string
		hashStr, err = crypt.Generate([]byte(plaintext), []byte("$5$"+salt))
		hash = []byte(hashStr)
	case "sha-512":
		crypt := crypt.SHA512.New()
		var hashStr string
		hashStr, err = crypt.Generate([]byte(plaintext), []byte("$6$"+salt))
		hash = []byte(hashStr)
	}

	if err != nil {
		return "", err
	}

	parts := strings.Split(strings.TrimSpace(string(hash)), "$")
	return parts[len(parts)-1], nil
}
