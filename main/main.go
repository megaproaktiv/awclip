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
    args := os.Args[1:]
	
	commandLine := "aws"
	command := os.Args[2]
	for i:= 1; i < len(os.Args)-1 ; i++{
		commandLine += os.Args[1+i]
	}
	// fmt.Println(prg, ":", args, ":", commandLine)
	var cmd *exec.Cmd
	switch countArgs := len(args) ; countArgs {
	case 1:
		cmd = exec.Command(prg, os.Args[1])
	case 2:
		cmd = exec.Command(prg, os.Args[1], os.Args[2])
	case 3:
		cmd = exec.Command(prg, os.Args[1],os.Args[2], os.Args[3])
	case 4:
		cmd = exec.Command(prg, os.Args[1],os.Args[2], os.Args[3], os.Args[4])
	case 5:
		cmd = exec.Command(prg,  os.Args[1],os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	case 6:
		cmd = exec.Command(prg, os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6])
	case 7:
		cmd = exec.Command(prg, os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6],os.Args[7])
	case 8:
		cmd = exec.Command(prg, os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6],os.Args[7],os.Args[8])
	case 9:
		cmd = exec.Command(prg, os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6],os.Args[7],os.Args[8],os.Args[9])
	case 10:
		cmd = exec.Command(prg, os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6],os.Args[7],os.Args[8],os.Args[9],os.Args[10])
	case 11:
		cmd = exec.Command(prg, os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6],os.Args[7],os.Args[8],os.Args[9],os.Args[10],os.Args[11])
	case 12:
		cmd = exec.Command(prg, os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6],os.Args[7],os.Args[8],os.Args[9],os.Args[10],os.Args[11],os.Args[12])
	case 13:
		cmd = exec.Command(prg, os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6],os.Args[7],os.Args[8],os.Args[9],os.Args[10],os.Args[11],os.Args[12],os.Args[13])
	case 14:
		cmd = exec.Command(prg, os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6],os.Args[7],os.Args[8],os.Args[9],os.Args[10],os.Args[11],os.Args[12],os.Args[13],os.Args[14])
	default:
		panic("Too much arguments, more than 14")
	}

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
