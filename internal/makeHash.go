package internal

import (
	"os/exec"
	"strings"
)

func MakeHash(plaintext, algorithm, salt string) (string, error) {
	cmd := exec.Command("mkpasswd", "-m", algorithm, "-s", "-S", salt)
	cmd.Stdin = strings.NewReader(plaintext)
	hash, err := cmd.Output()

	if err != nil {
		return "", err
	}

	parts := strings.Split(strings.TrimSpace(string(hash)), "$")
	return parts[len(parts)-1], nil
}
