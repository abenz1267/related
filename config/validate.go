package config

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/exp/slices"
)

const listEntry = "%s (%s)"

func (config Config) Validate() {
	errors := 0

	missing := []string{}

	for _, v := range config.Groups {
		missing = append(missing, findMissingFragments(v, config)...)
	}

	errors += printList(missing, "Missing fragments")
	errors += printList(findMissingParents(config), "Missing parents")
	errors += printList(findAmbigiousFragments(config.Fragments), "Ambiguous fragment names")
	errors += printList(findAmbigiousFragments(config.Parents), "Ambiguous parent names")
	errors += printList(findMalformedParentNames(config), "Malformed parent names")
	errors += printList(findMalformedFragmentNames(config), "Malformed fragment names")

	if errors == 0 {
		color.Green("Everything is ok")
	}
}

func findMalformedFragmentNames(config Config) []string {
	res := []string{}

	for _, v := range config.Fragments {
		if len(strings.Split(v.Name, "/")) > 2 {
			res = addToRes(res, v.Name, v.ConfigFile)
		}
	}

	return res
}

func findMalformedParentNames(config Config) []string {
	res := []string{}

	for _, v := range config.Parents {
		if strings.Contains(v.Name, "/") {
			res = addToRes(res, v.Name, v.ConfigFile)
		}
	}

	return res
}

func findMissingParents(config Config) []string {
	res := []string{}

	parents := config.getParentNames()

	for _, v := range config.Fragments {
		if v.getParentName() != "" {
			if !slices.Contains(parents, v.getParentName()) {
				res = addToRes(res, v.getParentName(), v.ConfigFile)
			}
		}
	}

	return res
}

func findMissingFragments(group Group, config Config) []string {
	res := []string{}

	fragments := config.getFragmentNames()
	parents := config.getParentNames()

	for _, v := range group.Fragments {
		if !slices.Contains(fragments, v) && !slices.Contains(parents, v) {
			res = addToRes(res, v, group.ConfigFile)
		}
	}

	return res
}

func findAmbigiousFragments(fragments []Fragment) []string {
	found := []string{}
	res := []string{}

	for _, v := range fragments {
		if slices.Contains(found, v.Name) {
			res = addToRes(res, v.Name, v.ConfigFile)
		} else {
			found = append(found, v.Name)
		}
	}

	return res
}

func addToRes(res []string, name, file string) []string {
	val := fmt.Sprintf(listEntry, name, file)

	if !slices.Contains(res, val) {
		res = append(res, val)
	}

	return res
}

func printList(list []string, title string) int {
	if len(list) == 0 {
		return 0
	}

	fmt.Printf("%s:\n", title) //nolint

	for _, v := range list {
		color.Red("- %s", v)
	}

	fmt.Println() //nolint

	return len(list)
}
