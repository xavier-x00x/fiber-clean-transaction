package utils

import "sort"

func Contains(s []string, str string) bool {
	sort.Strings(s)
	return sort.SearchStrings(s, str) != len(s)
}