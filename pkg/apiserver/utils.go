package apiserver

import "regexp"

func getDmPathFromAbsPath(path string) string {
	m := regexp.MustCompile(`[0-9]+\.`)
	s := m.ReplaceAllString(path, "")
	return s
}
