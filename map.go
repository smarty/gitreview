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

func printMap(m map[string]string, preamble string) {
	if len(m) == 0 {
		return
	}
	log.Printf(preamble, len(m))
	log.Println()
	for path := range m {
		log.Println(path)
	}
	log.Println()
}
