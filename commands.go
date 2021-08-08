package awclip

import (
	"crypto/md5"
	"encoding/hex"
	"log"
)

func CommandLine(args []string) *string {
	seperator := " "
	commandLine := "aws"
	for i := 0; i < len(args)-1; i++ {
		commandLine += seperator + args[1+i]
	}
	if debug {
		log.Println("CommandLine: ", &commandLine)
	}
	return &commandLine
}

// HashVaule calculates the id of the cahce entries
func HashValue(commandLine *string) *string {
	hash := md5.Sum([]byte(*commandLine))
	hashstring := hex.EncodeToString(hash[:])
	return &hashstring

}
