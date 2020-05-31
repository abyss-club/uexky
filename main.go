package main

import (
	"math/rand"
	"time"

	"gitlab.com/abyss.club/uexky/cmd"
)

func main() {
	rand.Seed(time.Now().Unix())
	cmd.Execute()
}
