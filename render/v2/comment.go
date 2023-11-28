package v2

import "strings"

func parseComment(comment string) string {
	comment = strings.ReplaceAll(comment, "//", "") // 去除前面的 '//'
	return strings.TrimSpace(comment)               // 去除前后空格
}
