package main

import (
	"math/rand"

	"github.com/segmentio/ksuid"
)

func randomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

func generateID(prefix string) string {
	return prefix + "_" + ksuid.New().String()
}

func randomFloat(max float32) float32 {
	return rand.Float32() * max
}

func randomElement(list []string) string {
	return list[randomInt(0, len(list))]
}
