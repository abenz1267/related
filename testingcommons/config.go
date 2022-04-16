package testingcommons

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/abenz1267/related/config"
)

func CreateValidConfig() {
	cfg := config.Config{
		Types: []config.Type{
			{
				Name:     "component",
				Path:     "testfiles",
				Template: GetName(ProjectTemplate),
				Pre:      "",
				Post:     "",
				Suffix:   ".tsx",
			},
			{
				Name:     "module",
				Path:     "testfiles/styles",
				Template: "",
				Pre:      "",
				Post:     "",
				Suffix:   "",
			},
		},
		Groups: []config.Group{{
			Name:  "component",
			Pre:   "",
			Post:  "",
			Types: []string{"component", "module"},
		}},
	}

	b, err := json.Marshal(cfg)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(config.ConfigFile, b, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
}
