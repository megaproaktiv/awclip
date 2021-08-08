package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/megaproaktiv/awclip"
)

var (
	// dynamically create regions
	regions = []string{
		"eu-north-1",
		"ap-south-1",
		"eu-west-3",
		"eu-west-2",
		"eu-west-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"ap-northeast-1",
		"sa-east-1",
		"ca-central-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"eu-central-1",
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
	}
)

func PrefetchEc2DescribeInstancesProxyWrapper(newCacheEntry *awclip.CacheEntry, args []string) {
	if debug {
		fmt.Println("Start prefetch")
	}
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
	PrefetchEc2DescribeInstancesProxy(newCacheEntry, client, args)
}

func PrefetchEc2DescribeInstancesProxy(config *awclip.CacheEntry, client Ec2Interface, args []string) error {
	var wg sync.WaitGroup

	seperator := "_"

	for _, region := range regions {
		region := region

		if debug {
			fmt.Println("Range Region: ", region)
		}
		//reCalc ID
		commandLine := "aws"
		args := replaceRegion(args, region)
		for i := 1; i < len(args)-1; i++ {
			commandLine += seperator + args[1+i]
		}
		hash := md5.Sum([]byte(commandLine))
		hashstring := hex.EncodeToString(hash[:])
		id := &hashstring

		regionalEntry := *config
		regionalEntry.Parameters.Region = &region

		if awclip.CacheMiss(id) {

			wg.Add(1)
			go calcInstances(&wg, id, args, region, client)

		}
		if debug {
			fmt.Println("Prefetch - cache miss: ", region)
		}

	}
	wg.Wait()

	return nil

}

func replaceRegion(args []string, region string) []string {
	for i, v := range args {
		if v == "--region" {
			args[i+1] = region

			break
		}
	}
	return args
}

func calcInstances(wg *sync.WaitGroup, id *string, args []string, region string, client Ec2Interface) {
	defer wg.Done()
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
	newCacheEntry.Parameters.Region = &region
	newParms := newCacheEntry.Parameters
	newCacheEntry.Provider = "go"
	if debug {
		fmt.Println("Prefetch - Call proxy: ", region)
	}

	localContent := Ec2DescribeInstancesProxy(newCacheEntry)
	if debug {
		fmt.Println("Prefetch - localContent: ", len(*localContent))
	}

	awclip.WriteContent(id, localContent)
	start := time.Now()
	metadata := &awclip.CacheEntry{
		Id:            id,
		Cmd:           aws.String("prefetch"),
		Created:       start,
		LastAccessed:  start,
		AccessCounter: 0,
		Parameters:    newParms,
		Provider:      newCacheEntry.Provider,
	}
	awclip.WriteMetadata(metadata)
}
