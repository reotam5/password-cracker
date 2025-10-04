package main

import (
	"fmt"
	"os"

	ParseShadowForUser "reotamai/assignment3/internal"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: readshadow <shadow_file_path> <username>")
		return
	}

	shadowPath := os.Args[1]
	username := os.Args[2]

	result, err := ParseShadowForUser.ParseShadowForUser(shadowPath, username)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Username: %s\n", result.Username)
	fmt.Printf("Algorithm: %s\n", result.Algorithm)
	fmt.Printf("Salt: %s\n", result.Salt)
	fmt.Printf("Hash: %s\n", result.Hash)
}
