package services

import (
	"os"
	"strings"
)

const TAB = "\t"
const NL = "\n"
var any = "*"

var AUTO_INIT bool
var Debug bool

func init() {
	AUTO_INIT = true
	Debug = false
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
