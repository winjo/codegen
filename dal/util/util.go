package util

import (
	"os"
	"regexp"
	"strings"
)

func AssertNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

func EnsureDir(dir string) {
	if len(dir) == 0 {
		return
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		AssertNotNil(err)
	}
}

func ToCamelCase(str string) string {
	return regexp.MustCompile(`_(\w)`).ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(s[1:])
	})
}

func ToPascalCase(str string) string {
	return UpperFirst(ToCamelCase(str))
}

func UpperFirst(str string) string {
	return strings.ToUpper(str[:1]) + str[1:]
}
