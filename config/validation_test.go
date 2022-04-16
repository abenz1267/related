package config_test

import (
	"testing"

	"github.com/abenz1267/related/config"
	"github.com/abenz1267/related/testingcommons"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestValidation(t *testing.T) {
	testingcommons.CreateTmpData()
	t.Parallel()
	t.Cleanup(testingcommons.Cleanup)

	cfg := config.Config{
		Types: []config.Type{
			{
				Name:     "component",
				Template: testingcommons.GetName(testingcommons.ProjectTemplate),
				Pre:      testingcommons.GetName(testingcommons.ProjectScript),
			},
			{
				Name:     "other",
				Template: testingcommons.GetName(testingcommons.ConfigTemplate),
				Pre:      testingcommons.GetName(testingcommons.ConfigScript),
			},
		},
		Groups: []config.Group{
			{
				Name:  "component",
				Types: []string{"component", "other"},
			},
		},
	}

	t1, c1 := config.Validate(cfg)

	assert.Empty(t, t1)
	assert.Empty(t, c1)

	cfg.Types = append(cfg.Types, config.Type{Name: "another", Template: "parent/missing", Pre: "parent/missing.lua"}) //nolint
	cfg.Groups = append(cfg.Groups, config.Group{Name: "another", Types: []string{"missing"}})                         //nolint

	t2, c2 := config.Validate(cfg)

	assert.True(t, slices.Contains(t2, "parent/missing"))
	assert.True(t, slices.Contains(t2, "parent/missing.lua"))
	assert.True(t, slices.Contains(c2, "missing"))
}
