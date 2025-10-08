package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ShadowResult struct {
	Username   string
	Algorithm  string
	Salt       string
	Hash       string
	Parameters string
	Raw        string
}

var HashAlgorithms = map[string]string{
	"1": "md5",

	"2a": "bcrypt",
	"2b": "bcrypt",
	"2y": "bcrypt",

	"5": "sha-256",

	"6": "sha-512",

	"y": "yescrypt",
}

func ParseShadowForUser(shadowPath string, username string) (*ShadowResult, error) {
	f, err := os.Open(shadowPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open shadow file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	prefix := username + ":"

	// look for the line that starts with the username:
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, prefix) {
			line = l
			break
		}
	}

	if line == "" {
		return nil, fmt.Errorf("user %s not found in shadow file", username)
	}

	// the field that contains the hash is the second field, separated by ':'
	fields := strings.Split(line, ":")
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid shadow file format for user %s", username)
	}

	hashField := fields[1]
	if hashField == "" || hashField == "*" || hashField == "!" {
		return nil, fmt.Errorf("user %s has no password set", username)
	}

	// the hash field contains different parts separated by '$'. (ie, algo, parameters, hash)
	hashParts := strings.Split(hashField, "$")
	shadowResult := &ShadowResult{}
	shadowResult.Raw = hashField

	if len(hashParts) == 4 {
		// $algo$salt$hash
		if _, exists := HashAlgorithms[hashParts[1]]; !exists {
			return nil, fmt.Errorf("unsupported hash algorithm: %s", hashParts[1])
		}

		shadowResult.Username = username
		shadowResult.Algorithm = HashAlgorithms[hashParts[1]]
		shadowResult.Salt = hashParts[2]
		shadowResult.Hash = hashParts[3]
		shadowResult.Parameters = ""
	} else if len(hashParts) == 5 {
		// $algo$params$salt$hash
		if _, exists := HashAlgorithms[hashParts[1]]; !exists {
			return nil, fmt.Errorf("unsupported hash algorithm: %s", hashParts[1])
		}

		shadowResult.Username = username
		shadowResult.Algorithm = HashAlgorithms[hashParts[1]]
		shadowResult.Parameters = hashParts[2]
		shadowResult.Salt = hashParts[3]
		shadowResult.Hash = hashParts[4]
	} else {
		return nil, fmt.Errorf("unsupported hash format for user %s", username)
	}

	return shadowResult, nil
}
