package main

import (
	"math/rand"

	"github.com/segmentio/ksuid"
)

func RandomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

func GenerateID(prefix string) string {
	return prefix + "_" + ksuid.New().String()
}

func RandomFloat(max float32) float32 {
	return rand.Float32() * max
}

func RandomElement(list []string) string {
	return list[RandomInt(0, len(list))]
}
