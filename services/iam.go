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
	ListAttachedUserPolicies(ctx context.Context,
		params *iam.ListAttachedUserPoliciesInput,
		optFns ...func(*iam.Options)) (*iam.ListAttachedUserPoliciesOutput, error)
	ListUsers(ctx context.Context,
		params *iam.ListUsersInput,
		optFns ...func(*iam.Options)) (*iam.ListUsersOutput, error)	
}

var IamListUserParameter = &awclip.Parameters{
	Service:    aws.String("iam"),
	Action:     aws.String("list-users"),
	Output:     aws.String("text"),
	Region:     aws.String("*"),
	Profile:    aws.String("*"),
	AdditionalParameters: map[string]*string{"user-name": &any},
	Query:      aws.String("Users[*].UserName"),
}

var IamListUserPoliciesParameter = &awclip.Parameters{
	Service:    aws.String("iam"),
	Action:     aws.String("list-user-policies"),
	Output:     aws.String("text"),
	Region:     aws.String("*"),
	Profile:    aws.String("*"),
	AdditionalParameters: map[string]*string{"user-name": &any},
	Query:      aws.String(""),
}

var IamListAttachedUserPoliciesParameter = &awclip.Parameters{
	Service:    aws.String("iam"),
	Action:     aws.String("list-attached-user-policies"),
	Output:     aws.String("text"),
	Region:     aws.String("*"),
	Profile:    aws.String("*"),
	AdditionalParameters: map[string]*string{"user-name": &any},
	Query:      aws.String(""),
}

func IamListUserPoliciesProxy(entry *awclip.CacheEntry) *string {
	if Debug {
		fmt.Println("IamListUserPoliciesProxy - Start : ", *entry.Parameters.Region)
	}

	entry.Provider = "go"
	var err error
	client := ProfiledIamClient(entry)

	var response *iam.ListUserPoliciesOutput
	iamParms := &iam.ListUserPoliciesInput{
		UserName: entry.Parameters.AdditionalParameters["user-name"],
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

func IamListUserProxy(entry *awclip.CacheEntry) *string {
	entry.Provider = "go"
	var err error
	client := ProfiledIamClient(entry)
	var response  *iam.ListUsersOutput 
	iamParams := &iam.ListUsersInput{}
	if len(*entry.Parameters.Region) > 4 {
		response, err = client.ListUsers(context.TODO(), iamParams, func(o *iam.Options) {
			o.Region = *entry.Parameters.Region
		})
	} else {
		response, err = client.ListUsers(context.TODO(), iamParams)
	}
	if err != nil {
		log.Println("Cant connect to iam service - listusers")
		log.Println("Region:", *entry.Parameters.Region)
		log.Fatal(err)
	}
	content := ""
	
	first := true
	for _, k := range response.Users {	
		if first {
			first = false
		}else{
			content = content + TAB
		}
		content = content + *k.UserName
	}
	
	content = content + NL

	return &content

}

func IamListAttachedUserPoliciesProxy(entry *awclip.CacheEntry) *string {
	entry.Provider = "go"
	var err error
	client := ProfiledIamClient(entry)

	var response *iam.ListAttachedUserPoliciesOutput
	iamParms := &iam.ListAttachedUserPoliciesInput{
		UserName: entry.Parameters.AdditionalParameters["user-name"],
	}
	if len(*entry.Parameters.Region) > 4 {
		response, err = client.ListAttachedUserPolicies(context.TODO(), iamParms, func(o *iam.Options) {
			o.Region = *entry.Parameters.Region
		})
	} else {
		response, err = client.ListAttachedUserPolicies(context.TODO(), iamParms)
	}
	if err != nil {
		log.Println("Cant connect to iam service")
		log.Println("Region:", *entry.Parameters.Region)
		log.Println("Parms:", *iamParms.UserName)
		log.Fatal(err)
	}
	content := ""
	header := "ATTACHEDPOLICIES"

	for _, k := range response.AttachedPolicies {	
		content = content + header+TAB+*k.PolicyArn + TAB + *k.PolicyName + NL
	}
	

	return &content
}


func ProfiledIamClient(entry *awclip.CacheEntry) *iam.Client{
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		// Specify the shared configuration profile to load.
		config.WithSharedConfigProfile(*entry.Parameters.Profile),
	)
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	return  iam.NewFromConfig(cfg)
}