package services

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/megaproaktiv/awclip/cache"

	gomiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	// smithy "github.com/aws/smithy-go"
	smithyio "github.com/aws/smithy-go/io"
	"github.com/aws/smithy-go/middleware"
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
	ServiceMap = make(map[string]interface{})
	ServiceMap["lambda:list-functions"] = LambdaListFunctionsProxy
	ServiceMap["ec2:describe-instances"] = Ec2DescribeInstancesProxy
	// // No Calls with Parameters
	// // ServiceMap["iam:list-attached-user-policies"] = IamListUserPoliciesProxy
	// ServiceMap["iam:list-users"] = IamListUserProxy
	// ServiceMap["iam:list-user-policies"] = IamListUserPoliciesProxy
	// ServiceMap["sts:get-caller-identity"] = StsGetCallerIdentityProxy

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

func ApiCallDumpFileNameCtx(ctx context.Context) *string {
	serviceId := gomiddleware.GetServiceID(ctx)
	operationName := gomiddleware.GetOperationName(ctx)
	region := gomiddleware.GetRegion(ctx)
	return ApiCallDumpFileNameString(&serviceId, &operationName, &region)
}

func ApiCallDumpFileNameString(serviceId *string, operationName *string, region *string) *string {
	postfix := ".json"
	if strings.ToLower(*serviceId) == "ec2" {
		postfix = ".xml"
	}
	ops := strings.Replace(*operationName, "-", "", 1)
	name := DATADIR + "/" + *serviceId + "_" + ops + "_" + *region + postfix
	normalized := strings.ToLower(name)
	return &normalized

}

func RawCacheMiss(entry *cache.CacheEntry) bool {
	return !RawCacheHit(entry)
}

// Cahce service jsoan answer
func RawCacheHit(entry *cache.CacheEntry) bool {
	prefetchName := ApiCallDumpFileNameString(entry.Parameters.Service, entry.Parameters.Action, entry.Parameters.Region)
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
	response := out.RawResponse.(*smithyhttp.Response)

	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(response.Body, ringBuffer)

	prefetchName := ApiCallDumpFileNameCtx(ctx)
	log.Println("PrefetchName: ",prefetchName)
	file, err := os.Create(*prefetchName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, body)

	defer file.Close()

	return out, metadata, nil
})

var HandleFinalize = middleware.FinalizeMiddlewareFunc("dumpjson", func(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleFinalize(ctx, in)
	if err != nil {
		return out, metadata, err
	}
	//response := out.RawResponse.(*smithyhttp.Response)
	response := out.Result

	ec2_DescribeInstancesOutput := response.(*ec2.DescribeInstancesOutput)

	u, err := json.Marshal(ec2_DescribeInstancesOutput)
	if err != nil {
		panic(err)
	}

	prefetchName := ApiCallDumpFileNameCtx(ctx)

	file, err := os.Create(*prefetchName + ".wew")
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.WriteString(file, string(u))

	defer file.Close()

	return out, metadata, nil
})
