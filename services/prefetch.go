package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/megaproaktiv/awclip/cache"
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

func PrefetchEc2DescribeInstancesProxyWrapper(newCacheEntry *cache.CacheEntry, args []string) {
	if Debug {
		fmt.Println("Start prefetch")
	}
	newCacheEntry.Provider = "go-prefetch"
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		// Specify the shared configuration profile to load.
		config.WithSharedConfigProfile(*newCacheEntry.Parameters.Profile),
	)
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	PrefetchEc2DescribeInstancesProxy(newCacheEntry, cfg)
}

func PrefetchLambdaListFunctionsProxyWrapper(metadata *cache.CacheEntry) {
	if Debug {
		fmt.Println("Start prefetch")
	}
	metadata.Provider = "go"
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		// Specify the shared configuration profile to load.
		config.WithSharedConfigProfile(*metadata.Parameters.Profile),
	)
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := lambda.NewFromConfig(cfg)
	if metadata.Parameters.Region == nil {
		metadata.Parameters.Region = &cfg.Region
	}
	if Debug {
		fmt.Println("prefetch line 75 Region:", *metadata.Parameters.Region)
	}
	PrefetchLambdaListFunctionsProxy(metadata, client)
}

func PrefetchEc2DescribeInstancesProxy(metadata *cache.CacheEntry, cfg aws.Config) error {
	var wg sync.WaitGroup

	for _, region := range regions {

		if Debug {
			fmt.Println("Range Region: ", region)
		}

		metadata.Parameters.Region = &region
		//reCalc ID
		id := metadata.Parameters.HashValue()

		regionalEntry := *metadata
		regionalEntry.Parameters.Region = &region

		if cache.CacheMiss(id) {

			wg.Add(1)
			go calcInstances(&wg, id, metadata, cfg)

		}
		if Debug {
			fmt.Println("Prefetch - cache miss: ", region)
		}

	}
	wg.Wait()

	return nil

}
func PrefetchLambdaListFunctionsProxy(metadata *cache.CacheEntry, client LambdaInterface) error {
	var wg sync.WaitGroup

	var MetadataRegionMap = make(map[string]*cache.CacheEntry)
	for _, region := range regions {

		if Debug {
			fmt.Println("Range Region: ", region)
		}
		//reCalc ID
		MetadataRegionMap[region] = metadata.Copy()

		// If region is bound to goroutine, the value will change
		MetadataRegionMap[region].Parameters.Region = aws.String(region)

		MetadataRegionMap[region].Id = MetadataRegionMap[region].Parameters.HashValue()

		if Debug {
			fmt.Println("prefetch:129 Lambda  id:", *MetadataRegionMap[region].Id)
			fmt.Println("prefetch:130 Lambda  region:", *MetadataRegionMap[region].Parameters.Region)
			fmt.Printf("prefetch:131 Meta %v \nLocalmeta %v", MetadataRegionMap[region], metadata)
		}

		if cache.CacheMiss(MetadataRegionMap[region].Id) {
			wg.Add(1)
			go calcFunctions(&wg, MetadataRegionMap[region])

		}
		if Debug {
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

func calcInstances(wg *sync.WaitGroup, id *string, metadata *cache.CacheEntry, cfg aws.Config) {
	defer wg.Done()

	Ec2DescribeInstancesProxy(metadata, cfg)

	start := time.Now()

	metadata.Created = start
	metadata.LastAccessed = start
	metadata.AccessCounter = 0
	metadata.Provider = "go - prefetch"
	cache.WriteMetadata(metadata)
}

func calcFunctions(wg *sync.WaitGroup, meta *cache.CacheEntry) {
	defer wg.Done()

	if Debug {
		fmt.Println("prefetch:184 ", meta.Parameters.Region)
		fmt.Println("prefetch Line 183 Region:", *meta.Parameters.Region)
		fmt.Println("prefetch Line 187 metadata *:", meta)
	}
	LambdaListFunctionsProxy(meta, aws.Config{})

	start := time.Now()
	meta.Created = start
	meta.LastAccessed = start
	meta.AccessCounter = 0
	meta.Provider = "go - prefetch"

	cache.WriteMetadata(meta)
}
