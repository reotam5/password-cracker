package internal

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	utils "reotamai/assignment3/pkg/utils"

	"github.com/schollz/progressbar/v3"
)

type PasswordCracker struct {
	Charset    []rune
	MinLength  int
	NumWorkers int
	BatchSize  int

	AssignmentMutex   sync.Mutex
	NextStartingChars []rune

	Found     bool
	FoundChan chan string

	StartTime     time.Time
	Attempts      *big.Int
	ProgressMutex sync.Mutex
	ProgressBar   *progressbar.ProgressBar
}

func (pc *PasswordCracker) getAssignedCharsAndRotate() []rune {
	// this is sort of like a queue, we assign a batch of passwords to a worker and then rotate the starting chars for the next worker
	// we do this with a mutex
	pc.AssignmentMutex.Lock()
	defer pc.AssignmentMutex.Unlock()

	current := make([]rune, len(pc.NextStartingChars))
	copy(current, pc.NextStartingChars)

	pc.NextStartingChars = utils.RotateString(pc.NextStartingChars, pc.BatchSize, pc.Charset)

	return current
}

func (pc *PasswordCracker) crack(validator func(string) (bool, error), wg *sync.WaitGroup) {
	defer wg.Done()

	// while other threads haven't found the password
	for !pc.Found {
		currentPassword := pc.getAssignedCharsAndRotate()

		for i := 0; i < pc.BatchSize; i++ {
			// other thread found the password
			if pc.Found {
				return
			}

			found, err := validator(string(currentPassword))

			if err != nil {
				fmt.Println("Error in validator:", err)
				return
			}

			// we don't need this mutex if we don't need insight analysis
			pc.ProgressMutex.Lock()
			pc.Attempts = new(big.Int).Add(pc.Attempts, big.NewInt(1))

			// update progress bar every 1000 attempts
			if new(big.Int).Mod(pc.Attempts, big.NewInt(1000)).Cmp(big.NewInt(0)) == 0 {
				elapsed := time.Since(pc.StartTime).Seconds()
				attemptsFloat := new(big.Float).SetInt(pc.Attempts)
				attemptsPerSec := new(big.Float).Quo(attemptsFloat, big.NewFloat(elapsed))
				pc.ProgressBar.Describe(fmt.Sprintf("Attempts: %s, Attempts/sec: %s, Current: %s", pc.Attempts.String(), attemptsPerSec.Text('f', 2), string(currentPassword)))
			}

			pc.ProgressMutex.Unlock()

			// this allows other threads to know we've found the password
			if found {
				pc.FoundChan <- string(currentPassword)
				pc.Found = true
				return
			}

			// get the next password to try (rotate by 1)
			currentPassword = utils.RotateString(currentPassword, 1, pc.Charset)
		}
	}
}

func (pc *PasswordCracker) CreateWorkers(validator func(string) (bool, error)) string {
	pc.ProgressBar = progressbar.NewOptions(-1)

	var wg sync.WaitGroup
	wg.Add(pc.NumWorkers)

	for i := 0; i < pc.NumWorkers; i++ {
		go pc.crack(validator, &wg)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// select will wait for either channel to receive a value
	select {
	case password := <-pc.FoundChan:
		wg.Wait()
		pc.ProgressBar.Finish()

		fmt.Println()
		fmt.Println()
		return password
	case <-done:
		pc.ProgressBar.Finish()

		fmt.Println()
		fmt.Println()
		return ""
	}
}
