package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/smithy-go/middleware"
	"github.com/jmespath/go-jmespath"

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

	//(base call)
	if Debug {
		fmt.Println("lambda:35 Region:", *llfpm.Parameters.Region)
		fmt.Println("lambda:36 *:", llfpm)
	}

	llfpm.Provider = "go"
	var err error
	client := ProfiledLambdaClient(llfpm)
	//var response *lambda.ListFunctionsOutput
	parms := &lambda.ListFunctionsInput{}
	if( Debug){
		fmt.Println("lambda line 39 Region:", *llfpm.Parameters.Region)
	}

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
	prefetchName := ApiCallDumpFileNameString(llfpm.Parameters.Service,llfpm.Parameters.Action,llfpm.Parameters.Region)
	jsondata, err := ioutil.ReadFile(*prefetchName)
    check(err)
    var data interface{}
	_ = json.Unmarshal(jsondata, &data)
	result, _ := jmespath.Search("Functions[*].{R:Runtime,F:FunctionName}", data)
	
	// to text
	content := ""
	for _,item:=range result.([]interface{}) {
		thisMap := item.(map[string]interface{})
		content += fmt.Sprintf("%v\t%v\n",thisMap["F"],thisMap["R"])
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
	
	// fmt.Printf("%T\n",response.Body)
	// fmt.Printf("%v\n",response.Body)
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(response.Body, ringBuffer)
	// check errors

	
	prefetchName := ApiCallDumpFileNameCtx(ctx)

	file, err := os.Create(*prefetchName)
    if err != nil {
		log.Fatal(err)
    }
	_, err = io.Copy(file, body)
	
	defer file.Close()
	
	// jsondata, err := ioutil.ReadFile(*prefetchName)
    // check(err)
    // var data interface{}
	// err = json.Unmarshal(jsondata, &data)
	// result, err := jmespath.Search("Functions[*].{R:Runtime,F:FunctionName}", data)
	// // result, err := jmespath.Search("Functions[*].{R:Runtime,N:FunctionName}", data)

	// fmt.Println("---")
	// for _,item:=range result.([]interface{}) {
	// 	thisMap := item.(map[string]interface{})
	// 	fmt.Printf("%v\t%v\n",thisMap["F"],thisMap["R"])
	// }
	// fmt.Println("---")
	
	// Middleware must call the next middleware to be executed in order to continue execution of the stack.
	// If an error occurs, you can return to prevent further execution.
	return next.HandleDeserialize(ctx, in)
})

func check(e error) {
    if e != nil {
        panic(e)
    }
}
	