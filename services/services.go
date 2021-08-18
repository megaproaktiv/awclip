package services

import (
	"context"
	"os"
	"strings"

	gomiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/megaproaktiv/awclip"
)

const TAB = "\t"
const NL = "\n"
var any = "*"

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
	name:= awclip.DATADIR+"/"+serviceId+"_"+operationName+"_"+region+".json"
	normalized := strings.ToLower(name)
	return &normalized
}


func ApiCallDumpFileNameString(serviceId *string, operationName *string, region *string) *string{
	ops := strings.Replace(*operationName, "-","",1)
	name:= awclip.DATADIR+"/"+*serviceId+"_"+ops+"_"+*region+".json"
	normalized := strings.ToLower(name)
	return &normalized
	
}

func RawCacheHit(entry *awclip.CacheEntry) bool{
	prefetchName := ApiCallDumpFileNameString(entry.Parameters.Service,entry.Parameters.Action,entry.Parameters.Region)
	location := awclip.DATADIR + string(os.PathSeparator) + *prefetchName
	info, err := os.Stat(location)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

