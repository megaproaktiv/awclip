package services

import (
	// "log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/megaproaktiv/awclip/cache"
)

func CallManager(metadata *cache.CacheEntry, cfg aws.Config) {
	serviceCall := *metadata.Parameters.Service + ":" + *metadata.Parameters.Action
	for service := range ServiceMap {
		// log.Printf("range %v . %v ",service,serviceCall)
		if service == serviceCall {
			//log.Print("service:", service)
			service := ServiceMap[service]
			service.(func(*cache.CacheEntry, aws.Config))(metadata, cfg)
			//log.Print("After")
			break
		}
	}
}
