package main

import (
	"time"

	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/megaproaktiv/awclip"
	"github.com/megaproaktiv/awclip/services"
)
const debug = false

func main() {

	prg := "aws"

	service := os.Args[1]
	action := os.Args[2]
	
	if debug {
		log.Println("Parameters: ", os.Args)
	}
	args := awclip.ArrangeParameters(os.Args)
	commandLine := awclip.CommandLine(args)
	id := awclip.HashValue(commandLine)

	cmd := exec.Command(prg, args[1:]...)

	var content *string

	discriminated := awclip.DiscriminatedCommand(&service,&action)
	if awclip.CacheHit(id) && !discriminated {
		// Hit
		content, _ = awclip.ReadContentUpdate(id)
		
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
				Parameters: map[string]*string{},
				Query:   new(string),
			},
			Provider: "tbd",
		}
		newCacheEntry.ArgumentsToCachedEntry(args)
		newParms := newCacheEntry.Parameters

		// CheckPrefetch
		if ( *newCacheEntry.Parameters.Region == services.FirstRegion) {
			if newParms.AlmostEqual(services.Ec2DescribeInstancesParameter) {
				services.PrefetchEc2DescribeInstancesProxyWrapper(newCacheEntry ,args)
			}
		}

		// Actions implemented in go
		if newParms.AlmostEqual(services.Ec2DescribeInstancesParameter) {
			content = services.Ec2DescribeInstancesProxy(newCacheEntry) 
		} 
		
		if newParms.AlmostEqual(services.StsGetCallerIdentityParameter) {
			content = services.StsGetCallerIdentityProxy(newCacheEntry)
		} 
		
		if newParms.AlmostEqual(services.Ec2DescribeRegionsParameter) {
			content = services.Ec2DescribeRegionsProxy(newCacheEntry)
		} 

		if newParms.AlmostEqual(services.IamListUserPoliciesParamater){
			content = services.IamListUserPoliciesProxy(newCacheEntry)
		}
		
		// no actions in go => use python cli
		if newCacheEntry.Provider == "tbd" {
			// just python aws cli
			newCacheEntry.Provider = "python"
			stdout, err := cmd.Output()
			if err != nil {
				log.Print(err.Error())

			}
			data := string(stdout)
			content = &data
		}

		start := time.Now()
		metadata := &awclip.CacheEntry{
			Id:            id,
			Cmd:           commandLine,
			Created:       start,
			LastAccessed:  start,
			AccessCounter: 0,
			Parameters:    newParms,
			Provider: newCacheEntry.Provider,
		}
		awclip.WriteMetadata(metadata)
		awclip.WriteContent(id, content)
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
