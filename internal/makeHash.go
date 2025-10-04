package internal

import (
	"os/exec"
	"strings"
)

func MakeHash(plaintext, algorithm, salt string) (string, error) {
	hash, err := exec.Command("mkpasswd", "-m", "-s", algorithm, "--", plaintext, salt).Output()

	if err != nil {
		return "", err
	}

	parts := strings.Split(strings.TrimSpace(string(hash)), "$")
	return parts[len(parts)-1], nil
}
