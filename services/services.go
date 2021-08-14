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

func init() {
	AUTO_INIT = true
	Debug = false
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
