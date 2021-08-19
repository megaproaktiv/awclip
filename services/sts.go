package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/megaproaktiv/awclip/cache"

)

//go:generate moq -out sts_moq_test.go . StSInterface
type StSInterface interface {
	GetCallerIdentity(ctx context.Context,
		params *sts.GetCallerIdentityInput,
		optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

var StsGetCallerIdentityParameter = &cache.Parameters{
	Service: aws.String("sts"),
	Action:  aws.String("get-caller-identity"),
	Output:  aws.String(""),
	Region:  aws.String("dontcare"),
	Profile: aws.String("unknown"),
	Query:   aws.String(""),
}

type GetCallerIdentity struct {
	Account *string
	Arn     *string
	UserId  *string
}

func StsGetCallerIdentityProxy(newCacheEntry *cache.CacheEntry, client StSInterface) *string {

	newCacheEntry.Provider = "go"
	
	var response *sts.GetCallerIdentityOutput
	var err error

	if len(*newCacheEntry.Parameters.Region) > 4 {
		response, err = client.GetCallerIdentity(context.TODO(), nil, func(o *sts.Options) {
			o.Region = *newCacheEntry.Parameters.Region
		})
	} else {
		response, err = client.GetCallerIdentity(context.TODO(), nil)
	}

	if err != nil {
		log.Println("Cant connect to sts service")
		log.Println("Region:", *newCacheEntry.Parameters.Region)
		log.Fatal(err)
	}

	caller := &GetCallerIdentity{
		Account: response.Account,
		Arn:     response.Arn,
		UserId:  response.UserId,
	}

	contentB, err := json.Marshal(caller)
	content := string(contentB)
	return &content
}
