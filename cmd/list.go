package cmd

import (
	"github.com/abenz1267/related/list"
	"github.com/spf13/cobra"
)

var List = &cobra.Command{
	Use:   "list [fragments, groups, templates, scripts]",
	Short: "list templates, scripts, types or groups",
}

var Fragments = &cobra.Command{
	Use:   "fragments",
	Short: "list all available fragments",
	Run: func(cmd *cobra.Command, args []string) {
		list.Fragments()
	},
}

var Parents = &cobra.Command{
	Use:   "parents",
	Short: "list all available parents",
	Run: func(cmd *cobra.Command, args []string) {
		list.Parents()
	},
}

var Groups = &cobra.Command{
	Use:   "groups",
	Short: "list all available groups",
	Run: func(cmd *cobra.Command, args []string) {
		list.Groups()
	},
}
