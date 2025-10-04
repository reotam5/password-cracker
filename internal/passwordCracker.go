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
	MaxLength  int
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

func (pc *PasswordCracker) getAndIncrementStartingChars() []rune {
	pc.AssignmentMutex.Lock()
	defer pc.AssignmentMutex.Unlock()

	current := make([]rune, len(pc.NextStartingChars))
	copy(current, pc.NextStartingChars)

	pc.NextStartingChars = utils.RotateString(pc.NextStartingChars, pc.BatchSize, pc.Charset)

	return current
}

func (pc *PasswordCracker) worker(validator func(string) (bool, error), wg *sync.WaitGroup) {
	defer wg.Done()

	// while other threads haven't found the password
	for !pc.Found {
		currentPassword := pc.getAndIncrementStartingChars()

		for i := 0; i < pc.BatchSize; i++ {
			// other thread found the password
			if pc.Found {
				return
			}

			// password not found within max length
			if len(currentPassword) > pc.MaxLength {
				return
			}

			found, err := validator(string(currentPassword))

			if err != nil {
				fmt.Println("Error in validator:", err)
				return
			}

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

			if found {
				pc.FoundChan <- string(currentPassword)
				pc.Found = true
				return
			}

			currentPassword = utils.RotateString(currentPassword, 1, pc.Charset)
		}
	}
}

func (pc *PasswordCracker) Crack(validator func(string) (bool, error)) string {
	pc.ProgressBar = progressbar.NewOptions(-1)

	var wg sync.WaitGroup
	wg.Add(pc.NumWorkers)

	for i := 0; i < pc.NumWorkers; i++ {
		go pc.worker(validator, &wg)
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
