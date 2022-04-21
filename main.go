package main

import (
	"log"

	"github.com/abenz1267/related/cmd"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "related",
		Short: "Related is a file generator.",
		Long: `Related is a tool that helps you quickly create files based on definitions. For more information
	checkout https://github.com/abenz1267/related`,
	}

	root.AddCommand(cmd.List)

	if err := root.Execute(); err != nil {
		log.Panic(err)
	}
}
