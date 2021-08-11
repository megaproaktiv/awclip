package services

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/megaproaktiv/awclip"
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

var Ec2DescribeInstancesParameter = &awclip.Parameters{
	Service: aws.String("ec2"),
	Action:  aws.String("describe-instances"),
	Output:  aws.String("text"),
	Region:  aws.String("eu-central-1"),
	Profile: aws.String("unknown"),
	Query:   aws.String("Reservations[*].Instances[*].[InstanceId]"),
}

var Ec2DescribeRegionsParameter = &awclip.Parameters{
	Service: aws.String("ec2"),
	Action:  aws.String("describe-regions"),
	Output:  aws.String("text"),
	Region:  aws.String("dontcare"),
	Profile: aws.String("unknown"),
	Query:   aws.String("Regions[].RegionName"),
}

func Ec2DescribeInstancesProxy(entry *awclip.CacheEntry) *string {

	if Debug {
		fmt.Println("Ec2DescribeInstancesProxy - Start : ", *entry.Parameters.Region)
	}

	entry.Provider = "go"
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		// Specify the shared configuration profile to load.
		config.WithSharedConfigProfile(*entry.Parameters.Profile),
	)
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := ec2.NewFromConfig(cfg)

	var response *ec2.DescribeInstancesOutput

	if len(*entry.Parameters.Region) > 4 {
		response, err = client.DescribeInstances(context.TODO(), nil, func(o *ec2.Options) {
			o.Region = *entry.Parameters.Region
		})
	} else {
		response, err = client.DescribeInstances(context.TODO(), nil)
	}

	if err != nil {
		log.Println("Cant connect to ec2 service")
		log.Println("Region:", *entry.Parameters.Region)
		log.Fatal(err)
	}
	// Content for --query_Reservations[*].Instances[*].[InstanceId]
	content := ""
	if *entry.Parameters.Query == "Reservations[*].Instances[*].[InstanceId]" {
		for _, v := range response.Reservations {
			for _, k := range v.Instances {
				content = content + *k.InstanceId
			}
		}
	}
	content += "\n"
	if Debug {
		fmt.Println("Ec2DescribeInstancesProxy - End : ", *entry.Parameters.Region)
	}
	return &content
}

func Ec2DescribeRegionsProxy(newCacheEntry *awclip.CacheEntry) *string {
	if Debug {
		fmt.Println("Start describe regions")
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

	var response *ec2.DescribeRegionsOutput

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

	length := len(response.Regions)-1
	for i, v := range response.Regions {
		content = content + *v.RegionName
		if i < length {
			content += "\t"
		}
	}

	content += "\n"
	return &content
}
