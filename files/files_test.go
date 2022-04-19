package files_test

import (
	"path/filepath"
	"testing"

	"github.com/abenz1267/related/files"
	"github.com/abenz1267/related/testingcommons"
	"github.com/stretchr/testify/assert"
)

func TestFindFile(t *testing.T) {
	t.Cleanup(testingcommons.Cleanup)
	testingcommons.CreateTmpData()

	type test struct {
		file     string
		dir      files.TypeDir
		expected string
	}

	tests := []test{
		{
			file:     testingcommons.GetName(testingcommons.ProjectTemplate),
			dir:      files.TemplateDir,
			expected: files.NormalizeForFS(filepath.Join(testingcommons.TemplateDir, testingcommons.Parent, testingcommons.ProjectTemplate)),
		},
		{
			file:     testingcommons.GetName(testingcommons.ConfigTemplate),
			dir:      files.TemplateDir,
			expected: files.NormalizeForFS(filepath.Join(testingcommons.TemplateDir, testingcommons.Parent, testingcommons.ConfigTemplate)),
		},
		{
			file:     "parent/missing",
			dir:      files.TemplateDir,
			expected: "",
		},
		{
			file:     testingcommons.GetName(testingcommons.ProjectScript),
			dir:      files.ScriptDir,
			expected: files.NormalizeForFS(filepath.Join(testingcommons.ScriptDir, testingcommons.Parent, testingcommons.ProjectScript)),
		},
		{
			file:     testingcommons.GetName(testingcommons.ConfigScript),
			dir:      files.ScriptDir,
			expected: files.NormalizeForFS(filepath.Join(testingcommons.ScriptDir, testingcommons.Parent, testingcommons.ConfigScript)),
		},
		{
			file:     "parent/missing.lua",
			dir:      files.ScriptDir,
			expected: "",
		},
	}

	for _, v := range tests {
		path, system := files.FindFile(v.file, v.dir)
		assert.Equal(t, v.expected, path, v)

		if path == "" {
			assert.Nil(t, system)
		} else {
			assert.NotNil(t, system)
		}
	}
}
