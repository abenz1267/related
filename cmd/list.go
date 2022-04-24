package cmd

import (
	"github.com/abenz1267/related/list"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.AddCommand(listFragmentsCmd)
	listCmd.AddCommand(listGroupsCmd)
	listCmd.AddCommand(listParentsCmd)
}

var listCmd = &cobra.Command{
	Use:   "list [fragments, groups, templates, scripts]",
	Short: "list templates, scripts, types or groups",
}

var listFragmentsCmd = &cobra.Command{
	Use:   "fragments",
	Short: "list all available fragments",
	Run: func(_ *cobra.Command, _ []string) {
		list.Fragments()
	},
}

var listParentsCmd = &cobra.Command{
	Use:   "parents",
	Short: "list all available parents",
	Run: func(_ *cobra.Command, _ []string) {
		list.Parents()
	},
}

var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "list all available groups",
	Run: func(_ *cobra.Command, _ []string) {
		list.Groups()
	},
}
