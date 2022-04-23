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

	cmd.List.AddCommand(cmd.Fragments)
	cmd.List.AddCommand(cmd.Groups)
	cmd.List.AddCommand(cmd.Parents)

	root.AddCommand(cmd.List)
	root.AddCommand(cmd.ValidateConfig)

	if err := root.Execute(); err != nil {
		log.Panic(err)
	}
}
