package main

import (
	"log"
	"os"
	"s9cmd/cmd"
)

func main() {
	app := cmd.Cmd()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
