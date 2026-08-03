package main

import "strings"

func checkContainsToc(headers []string) bool {
	for _, h := range headers {
		if strings.HasPrefix(h, "-- TOC entry") {
			return true
		}
	}
	return false
}

var name_markers = []string{
	"-- Name: ",
	"-- Data for Name: ",
}

func getSectionName(headers []string) (sectionName string) {
	for _, header := range headers {
		for _, marker := range name_markers {
			if after, ok := strings.CutPrefix(header, marker); ok {
				return after
			}
		}
	}
	return
}
