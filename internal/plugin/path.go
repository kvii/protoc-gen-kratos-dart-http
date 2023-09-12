package plugin

import (
	"path/filepath"
	"strings"
)

func changeExtension(path string, ext string) string {
	return strings.TrimSuffix(path, filepath.Ext(path)) + ext
}

func fileName(path string) string {
	path = filepath.Base(path)
	return strings.TrimSuffix(path, filepath.Ext(path))
}

func firstCharLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}
