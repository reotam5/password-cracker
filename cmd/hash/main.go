package main

import (
	"fmt"
	"os"

	"reotamai/assignment3/internal"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: hash <plaintext> <algorithm> <salt>")
		return
	}

	plaintext := os.Args[1]
	algorithm := os.Args[2]
	salt := os.Args[3]

	hash, err := internal.MakeHash(plaintext, algorithm, salt)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Hash: %s\n", hash)
}
