package services

import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/megaproaktiv/awclip/cache"

)

//go:generate moq -out lambda_moq_test.go . LambdaInterface

type LambdaInterface interface {
	ListFunctions(ctx context.Context,
	params *lambda.ListFunctionsInput,
	optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
}

var LambdaListFunctionsParameter = &cache.Parameters{
	Service: aws.String("lambda"),
	Action: aws.String("list-functions"),
	Output: aws.String("text"),
	Region:     &any,
	Profile:    &any,
	AdditionalParameters: map[string]*string{},
	Query:      aws.String("Functions[*].{R:Runtime,N:FunctionName}"),
}

func LambdaListFunctionsProxy(llfpm *cache.CacheEntry, client LambdaInterface, cfg aws.Config) {

	if client == nil{
		client = lambda.NewFromConfig(cfg)
	}


	llfpm.Provider = "go"
	var err error
	//var response *lambda.ListFunctionsOutput
	parms := &lambda.ListFunctionsInput{}

	if len(*llfpm.Parameters.Region) > 4 {
		_, err = client.ListFunctions(context.TODO(), parms, func(o *lambda.Options) {
			o.Region = *llfpm.Parameters.Region
		})
	} else {
		_, err = client.ListFunctions(context.TODO(), parms)
	}
	if err != nil {
		log.Println("Cant connect to lambda service")
		log.Println("Region:", *llfpm.Parameters.Region)
		log.Fatal(err)
	}

	// Query
	// prefetchName := ApiCallDumpFileNameString(llfpm.Parameters.Service,llfpm.Parameters.Action,llfpm.Parameters.Region)
	// jsondata, err := ioutil.ReadFile(*prefetchName)
	// if err != nil {
	// 	panic("read error, " + err.Error())
	// }
    // data := string(jsondata)
	// // to text
	// content := awclip.QueryText(&data, llfpm.Parameters.Query)
	
	
}


func check(e error) {
    if e != nil {
        panic(e)
    }
}
	