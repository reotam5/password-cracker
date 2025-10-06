package internal

import (
	"os/exec"
	"strings"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/md5_crypt"
	"github.com/go-crypt/crypt/algorithm/bcrypt"
	"github.com/openwall/yescrypt-go"
)

func MakeHash(plaintext string, algorithm string, salt string, parameters string) (string, error) {
	var err error
	var hash []byte

	// mkpasswd doesn't seem to support yescrypt, md5, or bcrypt
	switch algorithm {
	case "yescrypt":
		saltBytes := []byte("$y$" + parameters + "$" + salt)
		hash, err = yescrypt.Hash([]byte(plaintext), saltBytes)
	case "md5":
		crypt := crypt.MD5.New()
		var hashStr string
		hashStr, err = crypt.Generate([]byte(plaintext), []byte("$1$"+salt))
		hash = []byte(hashStr)
	case "bcrypt":
		hasher, _ := bcrypt.New()
		digest, _ := hasher.Hash(plaintext)
		hash = []byte(digest.Encode())
	default:
		cmd := exec.Command("mkpasswd", "-m", algorithm, "-s", "-S", salt)
		cmd.Stdin = strings.NewReader(plaintext)
		hash, err = cmd.Output()
	}

	if err != nil {
		return "", err
	}

	parts := strings.Split(strings.TrimSpace(string(hash)), "$")
	return parts[len(parts)-1], nil
}
