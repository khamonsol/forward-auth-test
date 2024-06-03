package util

import "strings"

func Intersection(slice1, slice2 []string) []string {
	elements := make(map[string]bool)
	for _, v := range slice1 {
		elements[strings.ToLower(v)] = true
	}

	intersection := make([]string, 0)
	for _, v := range slice2 {
		if elements[strings.ToLower(v)] {
			intersection = append(intersection, v)
			elements[strings.ToLower(v)] = false // to handle duplicates
		}
	}

	return intersection
}
