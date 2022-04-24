package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	Root.AddCommand(listCmd)
	Root.AddCommand(ValidateConfig)
	Root.AddCommand(createCmd)
}

var Root = &cobra.Command{
	Use:   "related",
	Short: "Related is a file generator.",
	Long: `Related is a tool that helps you quickly create files based on definitions. For more information
	checkout https://github.com/abenz1267/related`,
}
