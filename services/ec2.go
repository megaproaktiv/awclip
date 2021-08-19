package services

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/megaproaktiv/awclip/cache"
)

const FirstRegion = "eu-north-1"
const DefaultRegion = "us-west-1"

//go:generate moq -out ec2_moq_test.go . Ec2Interface
type Ec2Interface interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	DescribeRegions(ctx context.Context,
		params *ec2.DescribeRegionsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeRegionsOutput, error)
}

var Ec2DescribeInstancesParameter = &cache.Parameters{
	Service: aws.String("ec2"),
	Action:  aws.String("describe-instances"),
	Output:  aws.String("text"),
	Region:  aws.String("eu-central-1"),
	Profile: aws.String("unknown"),
	Query:   aws.String("Reservations[*].Instances[*].[InstanceId]"),
}

var Ec2DescribeRegionsParameter = &cache.Parameters{
	Service: aws.String("ec2"),
	Action:  aws.String("describe-regions"),
	Output:  aws.String("text"),
	Region:  aws.String("dontcare"),
	Profile: aws.String("unknown"),
	Query:   aws.String("Regions[].RegionName"),
}

func Ec2DescribeInstancesProxy(llfpm *cache.CacheEntry, cfg aws.Config) {

	client := ec2.NewFromConfig(cfg)

	llfpm.Provider = "go"
	var err error

	if len(*llfpm.Parameters.Region) > 4 {
		_, err = client.DescribeInstances(context.TODO(), nil, func(o *ec2.Options) {
			o.Region = *llfpm.Parameters.Region
		})
	} else {
		_, err = client.DescribeInstances(context.TODO(), nil)
	}

	if err != nil {
		log.Println("Cant connect to ec2 service")
		log.Println("Region:", *llfpm.Parameters.Region)
		log.Fatal(err)
	}

}

func Ec2DescribeRegionsProxy(newCacheEntry *cache.CacheEntry, client Ec2Interface) *string {
	if Debug {
		fmt.Println("Start describe regions")
	}

	newCacheEntry.Provider = "go"

	var response *ec2.DescribeRegionsOutput
	var err error
	if len(*newCacheEntry.Parameters.Region) > 4 {
		response, err = client.DescribeRegions(context.TODO(), nil, func(o *ec2.Options) {
			o.Region = *newCacheEntry.Parameters.Region
		})
	} else {
		response, err = client.DescribeRegions(context.TODO(), nil)
	}

	if err != nil {
		log.Println("Cant connect to ec2 service")
		log.Println("Region:", *newCacheEntry.Parameters.Region)
		log.Fatal(err)
	}
	// Content for --query_Reservations[*].Instances[*].[InstanceId]
	content := ""

	length := len(response.Regions) - 1
	for i, v := range response.Regions {
		content = content + *v.RegionName
		if i < length {
			content += "\t"
		}
	}

	content += "\n"
	return &content
}
