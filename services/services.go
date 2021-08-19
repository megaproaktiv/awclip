package services

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/megaproaktiv/awclip/cache"

	gomiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	// smithy "github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithyio "github.com/aws/smithy-go/io"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const TAB = "\t"
const NL = "\n"
var any = "*"
const DATADIR = ".awclip"

var AUTO_INIT bool
var Debug bool
var ServiceMap map[string]interface{}


func init() {
	AUTO_INIT = true
	Debug = false
	ServiceMap  = make(map[string]interface{})
	ServiceMap["ec2:describe-instances"] = Ec2DescribeInstancesProxy
	// No Calls with Parameters
	// ServiceMap["iam:list-attached-user-policies"] = IamListUserPoliciesProxy
	ServiceMap["iam:list-users"] = IamListUserProxy
	ServiceMap["iam:list-user-policies"] = IamListUserPoliciesProxy
	ServiceMap["lambda:list-functions"] = LambdaListFunctionsProxy
	ServiceMap["sts:get-caller-identity"] = StsGetCallerIdentityProxy


}

func Autoinit() bool {
	key := "AUTO_INIT"
	if value, ok := os.LookupEnv(key); ok {
		if strings.EqualFold(value, "false") {
			return false
		}
	}
	return true
}

func ApiCallDumpFileNameCtx(ctx context.Context) *string{
	serviceId := gomiddleware.GetServiceID(ctx)
	operationName := gomiddleware.GetOperationName(ctx)
	region := gomiddleware.GetRegion(ctx)
	name:= DATADIR+"/"+serviceId+"_"+operationName+"_"+region+".json"
	normalized := strings.ToLower(name)
	return &normalized
}


func ApiCallDumpFileNameString(serviceId *string, operationName *string, region *string) *string{
	ops := strings.Replace(*operationName, "-","",1)
	name:= DATADIR+"/"+*serviceId+"_"+ops+"_"+*region+".json"
	normalized := strings.ToLower(name)
	return &normalized
	
}

func RawCacheHit(entry *cache.CacheEntry) bool{
	prefetchName := ApiCallDumpFileNameString(entry.Parameters.Service,entry.Parameters.Action,entry.Parameters.Region)
	location := DATADIR + string(os.PathSeparator) + *prefetchName
	info, err := os.Stat(location)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}


// handleDeserialize to save the raw api call json output 
var HandleDeserialize = middleware.DeserializeMiddlewareFunc("dumpjson", func(
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
