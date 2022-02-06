package services

import (
	"context"
	"sync"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/megaproaktiv/awclip/cache"
)


func PrefetchIamListAttachedUserPoliciesProxy(entry *cache.CacheEntry, client *iam.Client, users *[]*string) error {
	entry.Provider = "go"
	var err error
	var wg sync.WaitGroup
	for _, user := range *users{
		wg.Add(1)
		go func(user *string, client *iam.Client) {
			defer wg.Done()
			// do something
			iamParms := &iam.ListAttachedUserPoliciesInput{
				UserName: user,
			}
			if len(*entry.Parameters.Region) > 4 {
				_, err = client.ListAttachedUserPolicies(context.TODO(), iamParms, func(o *iam.Options) {
					o.Region = *entry.Parameters.Region
				})
				} else {
					_, err = client.ListAttachedUserPolicies(context.TODO(), iamParms)
				}
			if err != nil {
				log.Println("Cant connect to iam service")
				log.Println("Region:", *entry.Parameters.Region)
				log.Println("Parms:", *iamParms.UserName)
				log.Fatal(err)
			}
		}(user,client)
	}
	wg.Wait()
	
	return nil
}


func IamListUsers(entry *cache.CacheEntry, client *iam.Client) (*[]*string, error) {
	var err error
	var response *iam.ListUsersOutput
	iamParams := &iam.ListUsersInput{}
	if len(*entry.Parameters.Region) > 4 {
		response, err = client.ListUsers(context.TODO(), iamParams, func(o *iam.Options) {
			o.Region = *entry.Parameters.Region
		})
		} else {
			response, err = client.ListUsers(context.TODO(), iamParams)
		}
		if err != nil {
			log.Println("Cant connect to iam service - listusers")
			log.Println("Region:", *entry.Parameters.Region)
			return nil, err
		}
		
		users := make([]*string, len(response.Users))
		for i, k := range response.Users {
		users[i] = k.UserName
	}

	return &users, nil
}
