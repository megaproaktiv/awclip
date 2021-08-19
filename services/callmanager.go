package services

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/megaproaktiv/awclip/cache"

)

func CallManager( metadata *cache.CacheEntry, cfg aws.Config){
	serviceCall := *metadata.Parameters.Service + ":" + *metadata.Parameters.Action
	for service := range ServiceMap {
		if ServiceMap[service] == serviceCall {
			service := ServiceMap[service]
			service.(func(*cache.CacheEntry, interface{}, aws.Config))(metadata,nil,cfg)
		}
	}
}