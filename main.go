package main

import (
	"math/rand"
	"time"

	"github.com/cnych/sinaimgmover/cmd"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	cmd.Execute()
}
