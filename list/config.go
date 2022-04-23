package list

import (
	"errors"
	"fmt"
	"log"

	"github.com/abenz1267/related/config"
)

func Fragments() {
	cfg, err := config.Get()
	if err != nil && !errors.Is(err, config.ErrFragmentNotFound) {
		log.Panic(err)
	}

	for _, v := range cfg.Fragments {
		fmt.Printf("%s (%s)\n", v.Name, v.ConfigFile) //nolint
	}
}

func Parents() {
	cfg, err := config.Get()
	if err != nil && !errors.Is(err, config.ErrFragmentNotFound) {
		log.Panic(err)
	}

	for _, v := range cfg.Parents {
		fmt.Printf("%s (%s)\n", v.Name, v.ConfigFile) //nolint
	}
}

func Groups() {
	cfg, err := config.Get()
	if err != nil && !errors.Is(err, config.ErrFragmentNotFound) {
		log.Panic(err)
	}

	for _, v := range cfg.Groups {
		fmt.Printf("%s (%s)\n", v.Name, v.ConfigFile) //nolint
	}
}
