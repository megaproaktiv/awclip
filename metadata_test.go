package awclip_test

import (
	"testing"
	"gotest.tools/assert"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/megaproaktiv/awclip"
)

func TestArgumentsToCachedEntry(t *testing.T) {
	type args struct {
		args []string
		item *awclip.CacheEntry
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
				item: &awclip.CacheEntry{
					Parameters: &awclip.Parameters{
					Action:        new(string),
					Output:        new(string),
					Region:        new(string),
					Profile:       new(string),
					Query:         new(string),
				},
			},
			},
			want: &awclip.Parameters{
				Action:        aws.String("describe-instances"),
				Output:        aws.String("text"),
				Region:        aws.String("eu-central-1"),
				Profile:       aws.String("myprofile"),
				Query:        aws.String("Reservations[*].Instances[*].[InstanceId]"),
				
			},
		
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.item.ArgumentsToCachedEntry(tt.args.args)
			assert.DeepEqual(t,*tt.args.item.Parameters, *tt.want)
			
		})
	}
}
