package cmd

import (
	"log"

	"github.com/abenz1267/related/config"
	"github.com/spf13/cobra"
)

var List = &cobra.Command{
	Use:   "list [type to list]",
	Short: "list templates, scripts, types or groups",
	// Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// s := spinner.New(spinner.CharSets[41], 100*time.Millisecond) // Build our new spinner
		// s.Suffix = "  Processing..."
		// s.Start()                   // Start the spinner
		// time.Sleep(4 * time.Second) // Run for some time to simulate work
		// s.Stop()
		// color.Green("Echo: " + strings.Join(args, " "))
		config, err := config.Get()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%+v", config)
	},
}
