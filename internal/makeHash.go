package internal

import (
	"os/exec"
	"strings"
)

func MakeHash(plaintext, algorithm, salt string) (string, error) {
	hash, err := exec.Command("mkpasswd", "-m", algorithm, plaintext, salt, "-s").Output()

	if err != nil {
		return "", err
	}

	parts := strings.Split(strings.TrimSpace(string(hash)), "$")
	return parts[len(parts)-1], nil
}
