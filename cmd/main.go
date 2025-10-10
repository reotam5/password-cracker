package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"reotamai/assignment3/internal"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// check args
	if len(os.Args) < 4 {
		fmt.Println("Usage: cracker <shadow file> <username> <number of threads>")
		return
	}
	shadowPath := os.Args[1]
	username := os.Args[2]
	numThreads, err := strconv.Atoi(os.Args[3])

	if err != nil || numThreads <= 0 {
		fmt.Println("Invalid number of threads:", os.Args[3])
		return
	}

	// parse shadow file for given user
	shadowResult, err := internal.ParseShadowForUser(shadowPath, username)
	if err != nil {
		fmt.Println("Error parsing shadow file:", err)
		os.Exit(1)
	}

	fmt.Printf("Cracking password for user: %s\n", shadowResult.Username)
	fmt.Printf("Threads: %d\n", numThreads)
	fmt.Printf("Algorithm: %s\n", shadowResult.Algorithm)
	fmt.Printf("Salt: %s\n", shadowResult.Salt)
	fmt.Printf("Hash: %s\n", shadowResult.Hash)

	// initialize a struct to be shared between workers
	pc := &internal.PasswordCracker{
		Charset:    []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#%^&*()_+-=.,:;?"),
		MinLength:  3,
		MaxLength:  256,
		NumWorkers: numThreads,
		BatchSize:  10,

		FoundChan: make(chan string, 1),

		StartTime: time.Now(),
		Attempts:  big.NewInt(0),
	}

	// start with the smallest possible password eg) "aaa"
	pc.NextStartingChars = make([]rune, pc.MinLength)
	for i := range pc.NextStartingChars {
		pc.NextStartingChars[i] = pc.Charset[0]
	}

	password := pc.CreateWorkers(func(attempt string) (bool, error) {
		// bcrypt is handled differently because its hash changes every time
		if shadowResult.Algorithm == "bcrypt" {
			err := bcrypt.CompareHashAndPassword([]byte(shadowResult.Raw), []byte(attempt))
			return err == nil, nil
		}

		// recompute hash and compare with correct hash
		attemptHash, err := internal.MakeHash(attempt, shadowResult.Algorithm, shadowResult.Salt, shadowResult.Parameters)

		if err != nil {
			return false, err
		}

		// if the hashes match, we found the password
		return attemptHash == shadowResult.Hash, nil
	})

	if password == "" {
		fmt.Println("Password not found.")
	} else {
		fmt.Printf("Password found: %s\n", password)
	}

	elapsed := time.Since(pc.StartTime)
	fmt.Printf("Elapsed time: %s\n", elapsed)
	fmt.Printf("Total attempts: %s\n", pc.Attempts.String())
}
