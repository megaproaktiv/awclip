package services

import (
	"github.com/megaproaktiv/awclip"
)

func CallManager( metadata *awclip.CacheEntry){
	serviceCall := *metadata.Parameters.Service + ":" + *metadata.Parameters.Action
	for service := range ServiceMap {
		if ServiceMap[service] == serviceCall {

		}
	}
}