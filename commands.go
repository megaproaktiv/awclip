package awclip

import (
	"crypto/md5"
	"encoding/hex"
)

const SEP = " "
var empty = ""

func (parms *Parameters) CommandLine() *string {
	commandLine := "aws "
	
	commandLine += *emptyWhenNil(parms.Service) + SEP + 
	*emptyWhenNil(parms.Action) + SEP + 
	*emptyWhenNil(parms.Region) + SEP +
	*emptyWhenNil(parms.Output) + SEP +
	*emptyWhenNil(parms.Query)
	return &commandLine
}

// HashValue calculates the id of the chache entries
func (parms *Parameters) HashValue() *string {
	commandLine := parms.CommandLine()
	hash := md5.Sum([]byte(*commandLine))
	hashstring := hex.EncodeToString(hash[:])
	return &hashstring
	
}


