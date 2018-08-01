package main

import (
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	var (
		err error
		gib *os.File
	)
	if gib, err = os.Create("gibberish.txt"); err != nil {
		log.Println(err.Error())
		return
	}
	defer gib.Close()

	for i := 0; i < 2e6; i++ {
		gib.WriteString(randStringRunes(20) + "\n")
	}
}

func randStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
