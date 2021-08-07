package awclip_test

import (
	"gotest.tools/assert"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/megaproaktiv/awclip"
)

func TestArgumentsToCachedEntry(t *testing.T) {
	type args struct {
		args []string
		item *awclip.CacheEntry
	}
	newCacheEntry := &awclip.CacheEntry{
		Parameters: awclip.Parameters{
			Service:  new(string),
			Action:  new(string),
			Output:  new(string),
			Region:  new(string),
			Profile: new(string),
			Query:   new(string),
		},
	}

	tests := []struct {
		name string
		args args
		want *awclip.Parameters
	}{
		{
			name: "recognise ec2 DescribeInstances",
			args: args{
				args: []string{
					"dist/awclip",
					"ec2",
					"describe-instances",
					"--query",
					"Reservations[*].Instances[*].[InstanceId]",
					"--region",
					"eu-central-1",
					"--output",
					"text",
					"--profile",
					"myprofile",
				},
				item: newCacheEntry,
				
			},
			want: &awclip.Parameters{
				Service: aws.String("ec2"),
				Action:  aws.String("describe-instances"),
				Output:  aws.String("text"),
				Region:  aws.String("eu-central-1"),
				Profile: aws.String("myprofile"),
				Query:   aws.String("Reservations[*].Instances[*].[InstanceId]"),
			},
		},		
		{
			name: "recognise ec2 Describe Regions",
			args: args{
				args: []string{
					"dist/awclip",
					"ec2",
					"describe-regions",
					"--profile",
					"helmut",
					"--query",
					"Regions[].RegionName",
					"--region",
					"eu-central-1",
					"--output",
					"text",
				},
				item: newCacheEntry,
			},
			want: &awclip.Parameters{
				Service: aws.String("ec2"),
				Action:  aws.String("describe-regions"),
				Output:  aws.String("text"),
				Region:  aws.String("eu-central-1"),
				Profile: aws.String("helmut"),
				Query:   aws.String("Regions[].RegionName"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			tt.args.item.ArgumentsToCachedEntry(tt.args.args)
			assert.DeepEqual(t, tt.args.item.Parameters, *tt.want)

		})
	}
}
