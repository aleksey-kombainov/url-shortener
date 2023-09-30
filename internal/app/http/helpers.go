package http

import "strings"

func ExtractMIMETypeFromStr(str string) string {
	mtypeSlice := strings.Split(str, ";")
	return strings.TrimSpace(mtypeSlice[0])
}
