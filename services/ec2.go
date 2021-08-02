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
}

// func init() {
// 	autoinit := Autoinit()
// 	if autoinit {
// 		cfg, err := config.LoadDefaultConfig(context.TODO())
// 		if err != nil {
// 			panic("configuration error, " + err.Error())
// 		}
// 		clientEc2 = ec2.NewFromConfig(cfg)
// 	}
// }

var Ec2DescribeInstancesParameter = &awclip.Parameters{
	Action:  aws.String("describe-instances"),
	Output:  aws.String("text"),
	Region:  aws.String("eu-central-1"),
	Profile: aws.String("dummy"),
	Query:   aws.String("Reservations[*].Instances[*].[InstanceId]"),
}

func Ec2DescribeInstancesProxy(config *awclip.CacheEntry, client Ec2Interface) *string {

	resp, err := client.DescribeInstances(context.TODO(), nil, func(o *ec2.Options) {
		o.Region = *config.Parameters.Region
	})

	if err != nil {
		fmt.Println("Cant connect ec2")
		log.Fatal(err)
	}
	content := ""
	if *config.Parameters.Query == "Reservations[*].Instances[*].[InstanceId]" {
		for _, v := range resp.Reservations {
			for _, k := range v.Instances {
				content = content + *k.InstanceId
			}
		}
	}
	return &content
}

