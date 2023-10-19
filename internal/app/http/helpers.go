package http

import (
	"strings"
)

func IsHeaderContainsMIMETypes(headerValues []string, searchValues []string) bool {
	for _, headerVal := range headerValues {
		for _, searchVal := range searchValues {
			if strings.Contains(headerVal, searchVal) {
				return true
			}
		}
	}
	return false
}
