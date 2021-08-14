package awclip

import (
	"os"
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

func CacheMiss(id *string) bool {
	return !CacheHit(id)
}

func CacheHit(id *string) bool {
	location := GetLocationMetaData(id)
	info, err := os.Stat(*location)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetLocationData(contentId *string) *string {
	location := DATADIR + string(os.PathSeparator) + *contentId + ".json"
	return &location
}
func GetLocationMetaData(contentId *string) *string {
	location := DATADIR + string(os.PathSeparator) + *contentId + "-db.json"
	return &location
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
