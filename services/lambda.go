package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	gomiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/smithy-go/middleware"

	// smithy "github.com/aws/smithy-go"
	smithyio "github.com/aws/smithy-go/io"
	smithyhttp "github.com/aws/smithy-go/transport/http"

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

	cfg.APIOptions = append(cfg.APIOptions, func(stack *middleware.Stack) error {
		// Attach the custom middleware to the beginning of the Desrialize step
		return stack.Deserialize.Add(handleDeserialize, middleware.After)
	})
	return  lambda.NewFromConfig(cfg)
}

// handleDeserialize to save the raw api call json output 
var handleDeserialize = middleware.DeserializeMiddlewareFunc("dumpjson", func(
	ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler,
) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err != nil {
		return out, metadata, err
	}
	response :=out.RawResponse.(*smithyhttp.Response)
	// if !ok {
	// 	return out, metadata, &smithy.DeserializationError{Err: fmt.Errorf("unknown transport type %T", out.RawResponse)}
	// }
	
	fmt.Printf("%T\n",response.Body)
	fmt.Printf("%v\n",response.Body)
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(response.Body, ringBuffer)
	// check errors
	serviceId := gomiddleware.GetServiceID(ctx)
	operationName := gomiddleware.GetOperationName(ctx)
	region := gomiddleware.GetRegion(ctx)
	
	file, err := os.Create(awclip.DATADIR+"/"+serviceId+"-"+operationName+"-"+region+".json")
    if err != nil {
		log.Fatal(err)
    }
	defer file.Close()
	_, err = io.Copy(file, body)

	// Middleware must call the next middleware to be executed in order to continue execution of the stack.
	// If an error occurs, you can return to prevent further execution.
	return next.HandleDeserialize(ctx, in)
})

	