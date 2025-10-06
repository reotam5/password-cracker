package internal

import (
	"os/exec"
	"strings"

	"github.com/openwall/yescrypt-go"
)

func MakeHash(plaintext string, algorithm string, salt string, parameters string) (string, error) {
	var err error
	var hash []byte

	// mkpasswd doesn't seem to support yescrypt
	if algorithm == "yescrypt" {
		saltBytes := []byte("$y$" + parameters + "$" + salt)
		hash, err = yescrypt.Hash([]byte(plaintext), saltBytes)
	} else {
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
