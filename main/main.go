package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/aws/smithy-go/middleware"
	"github.com/megaproaktiv/awclip"
	"github.com/megaproaktiv/awclip/cache"
	"github.com/megaproaktiv/awclip/services"
)

const debug = true

func main() {
	services.Debug = debug

	prg := "aws"

	service := os.Args[1]
	action := os.Args[2]

	if debug {
		log.Println("Parameters: ", os.Args)
	}
	args := awclip.ArrangeParameters(os.Args)
	
	parameters := cache.Parameters{}
	parameters.ArgumentsToCachedEntry(args)
	
	id := parameters.HashValue()
	

	metadata := &cache.CacheEntry{
		Id:         id,
		Cmd:        parameters.CommandLine(),
		Parameters: parameters,
		Provider:   "awscli",
	}

	if debug {
		fmt.Println("Main cmd: \n", *metadata.Cmd)
		fmt.Println("Main id: \n", *id)
	}
	cmd := exec.Command(prg, args[1:]...)

	var content *string

	discriminated := awclip.DiscriminatedCommand(&service, &action)
	if cache.CacheHit(id) && !discriminated {
		// Hit
		content, _ = awclip.ReadContentUpdate(id)

	}
	var cfg aws.Config
	var err error 
	if len(*metadata.Parameters.Profile) > 2 {
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			// Specify the shared configuration profile to load.
			config.WithSharedConfigProfile(*metadata.Parameters.Profile),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
		)
	}
	if err != nil {
		panic(err)
	}

	
	// Prefetch check116
	if ( services.RawCacheMiss(metadata) && *metadata.Parameters.Service == "iam" && *metadata.Parameters.Action == "list-attached-user-policies" ){
		client := iam.NewFromConfig(cfg)
		users,err := services.IamListUsers(metadata, client)
		cfg.APIOptions = append(cfg.APIOptions, func(stack *middleware.Stack) error {
			// Attach the custom middleware to the beginning of the Deserialize step
			return stack.Deserialize.Add(services.HandleDeserialize, middleware.After)
		})
		services.CallManager(metadata, cfg)
		if err == nil{
			if debug {
				log.Println("Prefetch")
				metadata.Print()
			}
			err = services.PrefetchIamListAttachedUserPoliciesProxy(metadata, client, users)
			if err != nil{
				log.Fatal(err)
			}

		}
	}

	// if services.RawCacheMiss(metadata) && !discriminated {
	// 	if err != nil {
	// 		panic("configuration error, " + err.Error())
	// 	}
	// 	cfg.APIOptions = append(cfg.APIOptions, func(stack *middleware.Stack) error {
	// 		// Attach the custom middleware to the beginning of the Desrialize step
	// 		return stack.Deserialize.Add(services.HandleDeserialize, middleware.After)
	// 	})
	// 	cfg.APIOptions = append(cfg.APIOptions, func(stack *middleware.Stack) error {
	// 		// Attach the custom middleware to the beginning of the Desrialize step
	// 		return stack.Finalize.Add(services.HandleFinalize,middleware.Before)
	// 	})
	// 	services.CallManager(metadata, cfg)
	// }

	//Miss => create entry
	if cache.CacheMiss(id) && !discriminated {
		// fastproxy available?
		newCacheEntry := metadata
		newParms := metadata.Parameters

		// CheckPrefetch
		if *newCacheEntry.Parameters.Region == services.FirstRegion {
			if newParms.AlmostEqual(services.Ec2DescribeInstancesParameter) {
				services.PrefetchEc2DescribeInstancesProxyWrapper(newCacheEntry, args)
			}
			if newParms.AlmostEqual((services.LambdaListFunctionsParameter)) {
				services.PrefetchLambdaListFunctionsProxyWrapper(newCacheEntry)
			}
		}

		// Actions implemented in go
		if newParms.AlmostEqual(services.Ec2DescribeInstancesParameter) {
			content = awclip.CallQuery(metadata)
		}
		if newParms.AlmostEqual(services.LambdaListFunctionsParameter) {
			//services.LambdaListFunctionsProxy(newCacheEntry,lambda.NewFromConfig(cfg))
			//log.Printf("Query - lambda")
			content = awclip.CallQuery(metadata)
		}

		// if newParms.AlmostEqual(services.StsGetCallerIdentityParameter) {
		// 	rawContent = services.StsGetCallerIdentityProxy(newCacheEntry, sts.NewFromConfig(cfg))
		// }

		// if newParms.AlmostEqual(services.Ec2DescribeRegionsParameter) {
		// 	rawContent = services.Ec2DescribeRegionsProxy(newCacheEntry,ec2.NewFromConfig(cfg))
		// }

		// if newParms.AlmostEqualWithParameters(services.IamListUserPoliciesParameter){
		// 	rawContent = services.IamListUserPoliciesProxy(newCacheEntry, iam.NewFromConfig(cfg))
		// }



		// if newParms.AlmostEqualWithParameters(services.IamListAttachedUserPoliciesParameter){
		// 	rawContent = services.IamListAttachedUserPoliciesProxy(newCacheEntry,iam.NewFromConfig(cfg))
		// }

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
		cache.WriteMetadata(metadata)
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
