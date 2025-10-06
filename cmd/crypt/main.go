package main

import (
	"fmt"

	"github.com/openwall/yescrypt-go"
)

func main() {
	salt := []byte("$y$j9T$BdxgWlFiefA2i2DNIvAoa1")
	key, _ := yescrypt.Hash([]byte("a2*"), salt)

	fmt.Printf("hash: %s\n", key)
}
