package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "related",
	Short: "Related is a file generator.",
	Long: `Related is a tool that helps you quickly create files based on definitions. For more information
	checkout https://github.com/abenz1267/related`,
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[41], 100*time.Millisecond) // Build our new spinner
		s.Suffix = "  Processing..."
		s.Start()                   // Start the spinner
		time.Sleep(4 * time.Second) // Run for some time to simulate work
		s.Stop()
		color.Green("Successfully Finished\n")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
