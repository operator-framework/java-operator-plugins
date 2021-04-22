package util

import (
	"path/filepath"
	"strings"
)

const (
	filePathSep  = string(filepath.Separator)
	javaPath     = "src" + filePathSep + "main" + filePathSep + "java"
	resourcePath = "src" + filePathSep + "main" + filePathSep + "resources"
)

func PrependJavaPath(filename string, pkg string) string {
	return javaPath + filePathSep + pkg + filePathSep + filename
}

func PrependResourcePath(filename string) string {
	return resourcePath + filePathSep + filename
}

func AsPath(s string) string {
	return strings.ReplaceAll(s, ".", filePathSep)
}
