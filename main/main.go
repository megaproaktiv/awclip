package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/megaproaktiv/awclip"
	"github.com/megaproaktiv/awclip/services"
)
const debug = false


func main() {
	services.Debug = false

	prg := "aws"

	service := os.Args[1]
	action := os.Args[2]
	
	if debug {
		log.Println("Parameters: ", os.Args)
	}
	args := awclip.ArrangeParameters(os.Args)
	parameters :=  awclip.Parameters{}
	parameters.ArgumentsToCachedEntry(args)
	id := parameters.HashValue()
	
	metadata := &awclip.CacheEntry{
		Id: id,
		Cmd: parameters.CommandLine(),
		Parameters: parameters,
		Provider: "awscli",
	}
	
	

	if debug {
		fmt.Println("Main cmd: \n",*metadata.Cmd)
		fmt.Println("Main id: \n",*id)
	}
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
		newCacheEntry :=metadata
		newParms := metadata.Parameters

		// CheckPrefetch
		if ( *newCacheEntry.Parameters.Region == services.FirstRegion) {
			if newParms.AlmostEqual(services.Ec2DescribeInstancesParameter ) {
				services.PrefetchEc2DescribeInstancesProxyWrapper(newCacheEntry ,args)
			}
			if newParms.AlmostEqual((services.LambdaListFunctionsParameter)){
				services.PrefetchLambdaListFunctionsProxyWrapper(newCacheEntry)
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

		if newParms.AlmostEqualWithParameters(services.IamListUserPoliciesParameter){
			content = services.IamListUserPoliciesProxy(newCacheEntry)
		}

		if newParms.AlmostEqual((services.IamListUserParameter)){
			content = services.IamListUserProxy(newCacheEntry)
		}
		
		if newParms.AlmostEqualWithParameters(services.IamListAttachedUserPoliciesParameter){
			content = services.IamListAttachedUserPoliciesProxy(newCacheEntry)
		}

		if newParms.AlmostEqual((services.LambdaListFunctionsParameter)){
			content = services.LambdaListFunctionsProxy(newCacheEntry)
		}
		
		// no actions in go => use python cli
		if newCacheEntry.Provider == "awscli" {
			// just python aws cli
			stdout, err := cmd.Output()
			if err != nil {
				log.Print(err.Error())

			}
			data := string(stdout)
			content = &data
		}

		metadata.Provider = newCacheEntry.Provider
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
