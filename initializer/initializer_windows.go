//go:build windows

package initializer

import "strings"

func cleanFilename(filename string) string {
	result := filename
	if strings.HasPrefix(result, "/") {
		result = result[1:]
	}
	return result
}
