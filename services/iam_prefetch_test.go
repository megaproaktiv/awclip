package services_test

import (
	"context"
	"github.com/megaproaktiv/awclip/services"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/megaproaktiv/awclip/cache"

	"github.com/aws/smithy-go/middleware"

)

func TestPrefetchIamListAttachedUserPoliciesProxy(t *testing.T) {
	var metadata *cache.CacheEntry = &cache.CacheEntry{
		Id:            aws.String("c5f4bf7f592e18f8465b1f33a0acf37b"),
		Cmd:           aws.String("aws iam list-attached-user-policies    user-name jenny"),
		AccessCounter: 0,
		Parameters:    cache.Parameters{
			Service:              aws.String("iam"),
			Action:               aws.String("list-attached-user-policies"),
			Output:               new(string),
			Region:               new(string),
			Profile:              new(string),
			AdditionalParameters: map[string]*string{
				 
				"user-name" : aws.String("jenny") ,
				
			},
			
		},
		
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	client := iam.NewFromConfig(cfg)
	users,err := services.IamListUsers(metadata, client)
	cfg.APIOptions = append(cfg.APIOptions, func(stack *middleware.Stack) error {
		// Attach the custom middleware to the beginning of the Deserialize step
		return stack.Deserialize.Add(services.HandleDeserialize, middleware.After)
	})
	services.CallManager(metadata, cfg)
	err = services.PrefetchIamListAttachedUserPoliciesProxy(metadata, client, users)
	
}
