package awclip

import (
	"strings"
	"unicode"
)

const DATADIR = ".awclip"
const debug = false

func SpaceStringsBuilder(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}




// ArrangeParameters
// if args contains "--profile" and "profilename", put them at the end
// so optimizer can regognize it
func ArrangeParameters(args []string) []string {
	for i, v := range args {
		if v == "--profile" {
			profileName := args[i+1]
			args = append(args[:i], args[i+2:]...)
			args = append(args, "--profile")
			args = append(args, profileName)
			break
		}
	}
	return args
}
