package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/abenz1267/related/config"
	"github.com/abenz1267/related/creation"
)

func main() {
	wd := flag.String("workingDir", "", "set the working directory")
	flag.Parse()

	if *wd != "" {
		err := os.Chdir(*wd)
		if err != nil {
			log.Fatal(err)
		}
	}

	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatalf("Possible commands: %s\n", strings.Join([]string{config.ListCmd}, ", "))
	}

	switch args[0] {
	case config.ListCmd:
		config.List(args)
	case creation.TypeCmd, creation.GroupCmd:
		cArgs := creation.CmdArgs{
			Kind:      args[0],
			Component: args[1],
			Name:      args[2],
		}

		creation.Create(cArgs)
	default:
		log.Printf("Unknown command '%s'\n", args[0])
	}
}
