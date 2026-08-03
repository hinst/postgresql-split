package main

import "strings"

// Returns folder and file to store section. Folder: might be empty
func getQualifiedName(sectionName string) (folder string, file string) {
	i := 0
	for i < len(sectionName) && ((sectionName[i] >= 'A' && sectionName[i] <= 'Z') || (sectionName[i] >= 'a' && sectionName[i] <= 'z')) {
		i++
	}
	if i == 0 || i == len(sectionName) {
		file = sectionName
	} else {
		folder = sectionName[:i]
		file = sectionName[i:]
	}
	file = getFileName(file)
	folder = getFileName(folder)
	return
}

func getFileName(name string) string {
	return shrinkFilename(sanitizeFilename(name))
}

func shrinkFilename(name string) string {
	return strings.TrimLeft(name, "_")
}

func sanitizeFilename(name string) string {
	var b strings.Builder
	b.Grow(len(name))
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	return b.String()
}
