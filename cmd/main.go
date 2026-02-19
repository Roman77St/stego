package main

import (
	"log"

	"github.com/Roman77St/stego/internal/cli"
)

func main() {
	err := cli.RunCLI()
	if err != nil {
		log.Fatal(err)
	}
}
