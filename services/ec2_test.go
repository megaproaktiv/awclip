package services_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"gotest.tools/assert"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/megaproaktiv/awclip"
	"github.com/megaproaktiv/awclip/services"
)

func TestEc2DescribeInstancesProxy(t *testing.T) {
	// make and configure a mocked Ec2Interface
	mockedEc2Interface := &services.Ec2InterfaceMock{
		DescribeInstancesFunc: func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
			var output ec2.DescribeInstancesOutput
			data, err := os.ReadFile("testdata/ec2-describe-instances-input.json")
			if err != nil {
				t.Error("Cant read input testdata")
				t.Error(err)
			}
			json.Unmarshal(data, &output)
			return &output, nil
		},
	}

	// use mockedEc2Interface in code that requires Ec2Interface
	// and then make assertions.
	os.Setenv("AUTO_INIT", "false")
	config := &awclip.CacheEntry{
		Parameters: awclip.Parameters{
			Action:  aws.String("describe-instances"),
			Output:  aws.String("text"),
			Region:  aws.String("eu-central-1"),
			Profile: aws.String("dummy"),
			Query:   aws.String("Reservations[*].Instances[*].[InstanceId]"),
		},
	}
	content := services.Ec2DescribeInstancesProxy(config, mockedEc2Interface)
	assert.Equal(t, "i-038834a1e9d61882a", *content)
}
