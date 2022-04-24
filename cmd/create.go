package cmd

import (
	"log"

	"github.com/abenz1267/related/creation"
	"github.com/spf13/cobra"
)

func init() {
	createCmd.AddCommand(createTemplate)
}

var createCmd = &cobra.Command{ //nolint
	Use:   "create [template, script, group, fragment]",
	Short: "create template, script, group, fragment",
}

var createTemplate = &cobra.Command{ //nolint
	Use:   "template [source] [destination]",
	Short: "create template for given file",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := cmd.Flags().GetBool("result")
		if err != nil {
			log.Panic(err)
		}

		global, err := cmd.Flags().GetBool("global")
		if err != nil {
			log.Panic(err)
		}

		creation.Template(result, global, args[0], args[1])
	},
}

func init() {
	createTemplate.Flags().BoolP("global", "g", false, "if set, generated file will be placed in user configuration folder")
	createTemplate.Flags().BoolP("result", "r", false, "if set, the path to the generated template will be printed back")
}
