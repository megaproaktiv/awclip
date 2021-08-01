package awclip

import (
	"os"
	"strings"
	"unicode"
)

const tmpdir = ".awclip"



func SpaceStringsBuilder(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}else{
            b.WriteRune('_')
        }
	}
	return b.String()
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
	location := tmpdir + string(os.PathSeparator) + *contentId+ ".json"
	return &location
}
func GetLocationMetaData(contentId *string) *string {
	location := tmpdir + string(os.PathSeparator) + *contentId + "-db.json"
	return &location
}
