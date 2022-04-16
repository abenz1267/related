package files_test

import (
	"testing"

	"github.com/abenz1267/related/files"
	"github.com/abenz1267/related/testingcommons"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestFiles(t *testing.T) {
	t.Cleanup(testingcommons.Cleanup)
	testingcommons.CreateTmpData()

	type test struct {
		dir      files.SystemDir
		expected []string
	}

	tests := []test{
		{
			dir:      files.SystemDir(files.TemplateDir),
			expected: []string{testingcommons.ProjectTemplate, testingcommons.ConfigTemplate},
		},
		{
			dir:      files.SystemDir(files.ScriptDir),
			expected: []string{testingcommons.ProjectScript, testingcommons.ConfigScript},
		},
	}

	for _, v := range tests {
		res := files.Availables(string(v.dir), "")

		for _, n := range v.expected {
			assert.True(t, slices.Contains(res[testingcommons.Parent], n))
		}
	}
}
