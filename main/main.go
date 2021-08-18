package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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

	var rawContent *string
	var content *string

	discriminated := awclip.DiscriminatedCommand(&service,&action)
	if awclip.CacheHit(id) && !discriminated {
		// Hit
		content, _ = awclip.ReadContentUpdate(id)
		
	}

	if services.RawCacheHit(metadata) && !discriminated{
		services.CallManager(metadata)
	}


	//Miss => create entry
	if awclip.CacheMiss(id) && !discriminated {
		// fastproxy available?
		newCacheEntry :=metadata
		newParms := metadata.Parameters
		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			// Specify the shared configuration profile to load.
			config.WithSharedConfigProfile(*newParms.Profile),
		)
		if err != nil {
			panic("configuration error, " + err.Error())
		}
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
			rawContent = services.Ec2DescribeInstancesProxy(newCacheEntry,ec2.NewFromConfig(cfg)) 
		} 
		
		if newParms.AlmostEqual(services.StsGetCallerIdentityParameter) {
			rawContent = services.StsGetCallerIdentityProxy(newCacheEntry, sts.NewFromConfig(cfg))
		} 
		
		if newParms.AlmostEqual(services.Ec2DescribeRegionsParameter) {
			rawContent = services.Ec2DescribeRegionsProxy(newCacheEntry,ec2.NewFromConfig(cfg))
		} 

		if newParms.AlmostEqualWithParameters(services.IamListUserPoliciesParameter){
			rawContent = services.IamListUserPoliciesProxy(newCacheEntry, iam.NewFromConfig(cfg))
		}

		if newParms.AlmostEqual((services.IamListUserParameter)){
			rawContent = services.IamListUserProxy(newCacheEntry,iam.NewFromConfig(cfg))
		}
		
		if newParms.AlmostEqualWithParameters(services.IamListAttachedUserPoliciesParameter){
			rawContent = services.IamListAttachedUserPoliciesProxy(newCacheEntry,iam.NewFromConfig(cfg))
		}

		if newParms.AlmostEqual((services.LambdaListFunctionsParameter)){
			rawContent = services.LambdaListFunctionsProxy(newCacheEntry,lambda.NewFromConfig(cfg))
		}
		
		// no actions in go => use python cli
		if newCacheEntry.Provider == "awscli" {
			// just python aws cli
			stdout, err := cmd.Output()
			if err != nil {
				log.Print(err.Error())

			}
			data := string(stdout)
			rawContent = &data
		}

		metadata.Provider = newCacheEntry.Provider
		awclip.WriteMetadata(metadata)
		awclip.WriteContent(id, rawContent)
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
