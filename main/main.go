package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/megaproaktiv/awclip"
	// "github.com/aws/aws-sdk-go-v2/aws"
	// "github.com/megaproaktiv/awclip"
)

func main() {

    prg := "aws"
	
	seperator := "_"
	commandLine := "aws"
	command := os.Args[1]
	commandLine += seperator+command
	os.Args = awclip.CleanUp(os.Args)
	for i:= 1; i < len(os.Args)-1 ; i++{
		commandLine += seperator+os.Args[1+i]
	}
	// fmt.Println(prg, ":", args, ":", commandLine)
	
	cmd := exec.Command(prg, os.Args[1:]...)
	
	hash := md5.Sum([]byte(commandLine))
	hashstring := hex.EncodeToString(hash[:])
	id := &hashstring
	start := time.Now()

	var content *string
	
	discriminated := awclip.DiscriminatedCommand(&command)
	if awclip.CacheHit(id) && !discriminated{
		// Hit
		content, _ = awclip.ReadContent(id)
		metadata,_ := awclip.ReadMetaData(id)
		err := awclip.UpdateMetaData(metadata)
		if err != nil {
			log.Print(err)
		}
	}

	if !awclip.CacheHit(id) && !discriminated{
		//Miss
		stdout, err := cmd.Output()
		if err != nil {
		 	log.Print(err.Error())
			
		}
		data := string(stdout)
		content = &data
		awclip.WriteContent(id, content)
		metadata := &awclip.CacheEntry{
			Id:            id,
			Cmd:           &commandLine,
			Created:       start,
			LastAccessed:  start,
			AccessCounter: 0,
		}
		awclip.WriteMetadata(metadata)
	}
	// Do not cache at all
	if discriminated{
		stdout, err := cmd.Output()
		if err != nil {
		 	log.Print(err.Error())
			
		}
		data := string(stdout)
		content = &data
	}

    fmt.Print(*content)
	
}
