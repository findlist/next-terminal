package main

import (
	"github.com/labstack/gommon/log"
	"next-terminal/server/app"
)

func main() {

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
