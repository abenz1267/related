package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var Create = &cobra.Command{ //nolint
	Use:   "create [template, script, group, fragment]",
	Short: "create template, script, group, fragment",
}

var Template = &cobra.Command{ //nolint
	Use:   "template [source] [destination]",
	Short: "create template for given file",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// bla
		log.Println(args)
		log.Println(cmd.Flags().GetBool("global"))
	},
}

func init() {
	Template.Flags().BoolP("global", "g", false, "if set, generated file will be placed in user configuration folder")
}
