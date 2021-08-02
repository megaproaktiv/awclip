package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/megaproaktiv/awclip"
	"github.com/megaproaktiv/awclip/services"
)

func main() {

    prg := "aws"
	
	seperator := "_"
	commandLine := "aws"
	command := os.Args[1]
	commandLine += seperator+command
	args := awclip.CleanUp(os.Args)
	for i:= 1; i < len(args)-1 ; i++{
		commandLine += seperator+args[1+i]
	}
	
	fmt.Println(prg, ":", args, ":", commandLine)
	
	cmd := exec.Command(prg, args[1:]...)
	
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

	//Miss
	if !awclip.CacheHit(id) && !discriminated{
		// fastproxy available
		testEntry := &awclip.CacheEntry{
			Parameters:    awclip.Parameters{},
		}
		testEntry.ArgumentsToCachedEntry(args)
		testParms := testEntry.Parameters
		if testParms.Equal( services.Ec2DescribeInstancesParameter) {

			cfg,err := config.LoadDefaultConfig(
				context.TODO(),
				// Specify the shared configuration profile to load.
				config.WithSharedConfigProfile(*testEntry.Parameters.Profile),
			)
			if err != nil {
				panic("configuration error, " + err.Error())
			}
			client := ec2.NewFromConfig(cfg)
			content = services.Ec2DescribeInstancesProxy(testEntry, client)

		}else {
			// just python aws cli
			stdout, err := cmd.Output()
			if err != nil {
				log.Print(err.Error())
				
			}
			data := string(stdout)
			content = &data
		}
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
