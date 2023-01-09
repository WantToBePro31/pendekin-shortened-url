package helper

import (
	"math/rand"
	"time"
)

var charSet = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateRandomUrl() string {
	rand.Seed(time.Now().UnixNano())
	rand_url := ""

	for i := 0; i < 7; i++ {
		rand_url += string(charSet[rand.Intn(len(charSet))])
	}

	runeUrl := []rune(rand_url)
	rand.Shuffle(len(runeUrl), func(i, j int) {
		runeUrl[i], runeUrl[j] = runeUrl[j], runeUrl[i]
	})

	return string(runeUrl)
}
