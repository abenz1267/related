package main

import (
	"log"

	"github.com/abenz1267/related/cmd"
)

func main() {
	if err := cmd.Root.Execute(); err != nil {
		log.Panic(err)
	}
}
