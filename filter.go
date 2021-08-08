package awclip

import (
	"github.com/aws/aws-sdk-go-v2/aws"

)

var discriminatedCommandsMap = map[string]*string{
	"iam": aws.String("generate-credential-report"),
}

// DiscriminatedCommand
// check whether a command is cachable
func DiscriminatedCommand(service *string, action *string) bool {
	value, exists := discriminatedCommandsMap[*service]
	if exists {
		if *value == * action {
			return true
		}
	} 
	return false
	
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
