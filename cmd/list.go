package cmd

import (
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var List = &cobra.Command{
	Use:   "list [type to list]",
	Short: "list templates, scripts, types or groups",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[41], 100*time.Millisecond) // Build our new spinner
		s.Suffix = "  Processing..."
		s.Start()                   // Start the spinner
		time.Sleep(4 * time.Second) // Run for some time to simulate work
		s.Stop()
		color.Green("Echo: " + strings.Join(args, " "))
	},
}
