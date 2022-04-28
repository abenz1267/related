package list

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abenz1267/related/config"
	"golang.org/x/exp/slices"
)

// unreliable for anything that doesn't end with <name>.<ext>.
func Files(file string) error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("Coulnd't get config: %w", err)
	}

	dir := filepath.Dir(file)

	fragments := []string{}

	for _, v := range cfg.Fragments {
		if v.Path == dir {
			fragments = append(fragments, v.Name)
		}
	}

	relatedFragments := []string{}

	for _, f := range fragments {
		for _, v := range cfg.Groups {
			if slices.Contains(v.Fragments, f) {
				relatedFragments = append(relatedFragments, v.Fragments...)
			}
		}
	}

	fragmentMap := map[string]struct{}{}

	for _, v := range relatedFragments {
		if !slices.Contains(fragments, v) {
			fragmentMap[v] = struct{}{}
		}
	}

	name := strings.Split(filepath.Base(file), ".")[0]

	res := []string{}

	for k := range fragmentMap {
		fragment, err := cfg.GetFragment(k, cfg.Fragments)
		if err != nil {
			return fmt.Errorf("couldn't get fragment: %w", err)
		}

		path := strings.ReplaceAll(fragment.Path, "<name>", strings.ToLower(name))
		filename := strings.ReplaceAll(fragment.Filename, "<name>", name)

		fullPath := filepath.Join(path, filename)

		if _, err := os.Stat(fullPath); !errors.Is(err, os.ErrNotExist) {
			res = append(res, fullPath)
		}
	}

	for _, v := range res {
		fmt.Println(v) //nolint
	}

	return nil
}
