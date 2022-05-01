package cmd

import (
	"errors"
	"log"

	"github.com/abenz1267/related/config"
	"github.com/spf13/cobra"
)

var ValidateConfig = &cobra.Command{
	Use:   "validate",
	Short: "validate config files",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil && !errors.Is(err, config.ErrFragmentNotFound) {
			log.Panic(err)
		}

		cfg.Validate()
	},
}
