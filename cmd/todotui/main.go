package main

import (
	"log"

	"github.com/yuucu/todotui/pkg/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
