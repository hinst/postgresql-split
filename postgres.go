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

func getSectionName(headers []string) (sectionName string) {
	for _, header := range headers {
		if strings.HasPrefix(header, "-- Name: ") {
			sectionName = strings.TrimPrefix(header, "-- Name: ")
		} else if strings.HasPrefix(header, "-- Data for Name: ") {
			sectionName = strings.TrimPrefix(header, "-- Data for Name: ")
		}
	}
	return
}
