package services

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/megaproaktiv/awclip"
)

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

func Ec2DescribeInstancesProxy(config *awclip.CacheEntry, client Ec2Interface) *string {

	var response *ec2.DescribeInstancesOutput
	var err error
	
	if len(*config.Parameters.Region) > 4 {
		response, err = client.DescribeInstances(context.TODO(), nil, func(o *ec2.Options) {
			o.Region = *config.Parameters.Region
		})
	}else{
		response, err = client.DescribeInstances(context.TODO(), nil)
	}

	if err != nil {
		log.Println("Cant connect to ec2 service")
		log.Println("Region:",*config.Parameters.Region)
		log.Fatal(err)
	}
	// Content for --query_Reservations[*].Instances[*].[InstanceId]
	content := ""
	if *config.Parameters.Query == "Reservations[*].Instances[*].[InstanceId]" {
		for _, v := range response.Reservations {
			for _, k := range v.Instances {
				content = content + *k.InstanceId
			}
		}
	}
	content += "\n"
	return &content
}


func Ec2DescribeRegionsProxy(config *awclip.CacheEntry, client Ec2Interface) *string {
	if debug {
		fmt.Println("Start describe regions")
	}
	var response *ec2.DescribeRegionsOutput
	var err error
	
	if len(*config.Parameters.Region) > 4 {
		response, err = client.DescribeRegions(context.TODO(), nil, func(o *ec2.Options) {
			o.Region = *config.Parameters.Region
		})
	}else{
		response, err = client.DescribeRegions(context.TODO(), nil)
	}

	if err != nil {
		log.Println("Cant connect to ec2 service")
		log.Println("Region:",*config.Parameters.Region)
		log.Fatal(err)
	}
	// Content for --query_Reservations[*].Instances[*].[InstanceId]
	content := ""
	
	length := len(response.Regions)
	for i, v := range response.Regions {
		content = content + *v.RegionName
		if i < length {
			content += "\t"
		}
	}


	content += "\n"
	return &content
}
