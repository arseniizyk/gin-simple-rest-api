package main

import (
	"log"

	"github.com/arseniizyk/server"
)

func main() {
	app := server.NewApp()
	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
