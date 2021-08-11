package services

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/megaproaktiv/awclip"
)

//go:generate moq -out lambda_moq_test.go . LambdaInterface

type LambdaInterface interface {
	ListFunctions(ctx context.Context,
	params *lambda.ListFunctionsInput,
	optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
}

var LambdaListFunctionsParameter = &awclip.Parameters{
	Service: aws.String("lambda"),
	Action: aws.String("list-functions"),
	Output: aws.String("text"),
	Region:     &any,
	Profile:    &any,
	AdditionalParameters: map[string]*string{},
	Query:      aws.String("Functions[*].{R:Runtime,N:FunctionName}"),
}

func LambdaListFunctionsProxy(llfpm *awclip.CacheEntry) *string{

	if Debug {
		fmt.Println("lambda:35 Region:", *llfpm.Parameters.Region)
		fmt.Println("lambda:36 *:", llfpm)
	}

	llfpm.Provider = "go"
	var err error
	client := ProfiledLambdaClient(llfpm)
	var response *lambda.ListFunctionsOutput
	parms := &lambda.ListFunctionsInput{}
	if( Debug){
		fmt.Println("lambda line 39 Region:", *llfpm.Parameters.Region)
	}

	if len(*llfpm.Parameters.Region) > 4 {
		response, err = client.ListFunctions(context.TODO(), parms, func(o *lambda.Options) {
			o.Region = *llfpm.Parameters.Region
		})
	} else {
		response, err = client.ListFunctions(context.TODO(), parms)
	}
	if err != nil {
		log.Println("Cant connect to lambda service")
		log.Println("Region:", *llfpm.Parameters.Region)
		log.Fatal(err)
	}
	content := ""
	
	for _, v := range response.Functions {
		content = content + *v.FunctionName + TAB + string(v.Runtime)+NL
	}

	
	return &content
}

func ProfiledLambdaClient(entry *awclip.CacheEntry) LambdaInterface{
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		// Specify the shared configuration profile to load.
		config.WithSharedConfigProfile(*entry.Parameters.Profile),
	)
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	return  lambda.NewFromConfig(cfg)
}