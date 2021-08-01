package services

import (
	"os"
	"strings"
)

var AUTO_INIT bool

func init() {
	AUTO_INIT = true
}

func Autoinit() bool {
	key := "AUTO_INIT"
	if value, ok := os.LookupEnv(key); ok {
		if strings.EqualFold(value, "false") {
			return false
		}
	}
	return true
}
