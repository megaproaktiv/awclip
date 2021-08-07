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
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/megaproaktiv/awclip"
	"github.com/megaproaktiv/awclip/services"
)
const debug = false

func main() {

	prg := "aws"

	seperator := "_"
	commandLine := "aws"
	command := os.Args[1]
	
	if debug {
		log.Println("Parameters: ", os.Args)
	}
	args := awclip.ArrangeParameters(os.Args)
	for i := 1; i < len(args)-1; i++ {
		commandLine += seperator + args[1+i]
	}

	if debug {
		log.Println("CommandLine: ", commandLine)
	}
	cmd := exec.Command(prg, args[1:]...)

	hash := md5.Sum([]byte(commandLine))
	hashstring := hex.EncodeToString(hash[:])
	id := &hashstring
	start := time.Now()

	var content *string

	discriminated := awclip.DiscriminatedCommand(&command)
	if awclip.CacheHit(id) && !discriminated {
		// Hit
		content, _ = awclip.ReadContent(id)
		metadata, _ := awclip.ReadMetaData(id)
		err := awclip.UpdateMetaData(metadata)
		if err != nil {
			log.Print(err)
		}
	}

	//Miss => create entry
	if awclip.CacheMiss(id) && !discriminated {
		// fastproxy available?
		newCacheEntry := &awclip.CacheEntry{
			Parameters: awclip.Parameters{
				Service: new(string),
				Action:  new(string),
				Output:  new(string),
				Region:  new(string),
				Profile: new(string),
				Query:   new(string),
			},
			Provider: "tbd",
		}
		newCacheEntry.ArgumentsToCachedEntry(args)
		newParms := newCacheEntry.Parameters
		if newParms.AlmostEqual(services.Ec2DescribeInstancesParameter) {
			newCacheEntry.Provider = "go"
			cfg, err := config.LoadDefaultConfig(
				context.TODO(),
				// Specify the shared configuration profile to load.
				config.WithSharedConfigProfile(*newCacheEntry.Parameters.Profile),
			)
			if err != nil {
				panic("configuration error, " + err.Error())
			}
			client := ec2.NewFromConfig(cfg)
			content = services.Ec2DescribeInstancesProxy(newCacheEntry, client) 
		} else if newParms.AlmostEqual(services.StsGetCallerIdentityParameter) {
			newCacheEntry.Provider = "go"
			cfg, err := config.LoadDefaultConfig(
				context.TODO(),
				config.WithSharedConfigProfile(*newCacheEntry.Parameters.Profile),
			)
			if err != nil {
				panic("configuration error, " + err.Error())
			}
			client := sts.NewFromConfig(cfg)
			content = services.StsGetCallerIdentityProxy(newCacheEntry, client)
		} else if newParms.AlmostEqual(services.Ec2DescribeRegionsParameter) {
			newCacheEntry.Provider = "go"
			cfg, err := config.LoadDefaultConfig(
				context.TODO(),
				// Specify the shared configuration profile to load.
				config.WithSharedConfigProfile(*newCacheEntry.Parameters.Profile),
			)
			if err != nil {
				panic("configuration error, " + err.Error())
			}
			client := ec2.NewFromConfig(cfg)
			content = services.Ec2DescribeRegionsProxy(newCacheEntry, client)
			if debug {
				fmt.Print("Content:",content)
			}
		} else {
			// just python aws cli
			newCacheEntry.Provider = "python"
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
			Parameters:    newParms,
			Provider: newCacheEntry.Provider,
		}
		awclip.WriteMetadata(metadata)
	}
	// Do not cache at all
	if discriminated {
		stdout, err := cmd.Output()
		if err != nil {
			log.Print(err.Error())

		}
		data := string(stdout)
		content = &data
	}

	fmt.Print(*content)

}
