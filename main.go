package main

import (
	"math/rand"
	"time"

	"github.com/icdb37/kypd/cmd"
)

func main() {
	rand.NewSource(time.Now().UnixNano())
	cmd.Execute()
}
