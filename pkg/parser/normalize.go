package parser

import "strings"

func NormalizeTeacherName(name string) string {
	name = strings.TrimSpace(name)

	name = strings.ReplaceAll(name, ". ", ".")

	for strings.Contains(name, "  ") {
		name = strings.ReplaceAll(name, "  ", " ")
	}

	return name
}
