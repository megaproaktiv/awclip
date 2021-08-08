package services

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/megaproaktiv/awclip"
)

//go:generate moq -out iam_moq_test.go . IamInterface
type IamInterface interface {
	ListUserPolicies(ctx context.Context,
		params *iam.ListUserPoliciesInput,
		optFns ...func(*iam.Options)) (*iam.ListUserPoliciesOutput, error)
}

var IamListUserPoliciesParamater = &awclip.Parameters{
	Service: aws.String("iam"),
	Action:  aws.String("list-user-policies"),
	Output:  aws.String("text"),
	Region:  aws.String("*"),
	Profile: aws.String("*"),
	Parameters: map[string]*string{ "user-name" : &any, },
	Query:   aws.String(""),
}

func IamListUserPoliciesProxy(entry *awclip.CacheEntry) *string {
	if debug {
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
	client := iam.NewFromConfig(cfg)
	var response *iam.ListUserPoliciesOutput
	iamParms := &iam.ListUserPoliciesInput{
		UserName: entry.Parameters.Parameters["user-name"],
	}
	if len(*entry.Parameters.Region) > 4 {
		response, err = client.ListUserPolicies(context.TODO(), iamParms, func(o *iam.Options) {
			o.Region = *entry.Parameters.Region
		})
	} else {
		response, err = client.ListUserPolicies(context.TODO(), iamParms)
	}
	if err != nil {
		log.Println("Cant connect to iam service")
		log.Println("Region:", *entry.Parameters.Region)
		log.Println("Parms:", *iamParms.UserName)
		log.Fatal(err)
	}
	content := ""
	for _, k := range response.PolicyNames {
		content = content + k
	}
	content += "\n"
	
	return &content
}