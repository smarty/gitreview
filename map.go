package main

import (
	"log"
	"sort"
)

func sortUniqueKeys(maps ...map[string]string) (unique []string) {
	combined := make(map[string]struct{})
	for _, m := range maps {
		for key := range m {
			combined[key] = struct{}{}
		}
	}
	for key := range combined {
		unique = append(unique, key)
	}
	sort.Strings(unique)
	return unique
}

func mapKeys(m map[string]string) (keys []string) {
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
func printMapKeys(m map[string]string, preamble string) {
	printStrings(mapKeys(m), preamble)
}
func printStrings(paths []string, preamble string) {
	log.Printf(preamble, len(paths))
	if len(paths) == 0 {
		return
	}
	for _, path := range paths {
		log.Println("  " + path)
	}
}
