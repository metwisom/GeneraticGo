package main

import (
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func Random(max int) int {
	if max < 1 {
		return 0
	}
	return rand.Intn(max)
}
