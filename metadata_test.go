package awclip

import (
	"testing"
	"gotest.tools/assert"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestArgumentsToCachedEntry(t *testing.T) {
	type args struct {
		args []string
		item *CacheEntry
	}
	tests := []struct {
		name string
		args args
		want *CacheEntry
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
				item: &CacheEntry{
					Id:            aws.String("abc"),
					Action:        new(string),
					Output:        new(string),
					Region:        new(string),
					Profile:       new(string),
					Query:         new(string),
				},
			},
			want: &CacheEntry{
				Id:            aws.String("abc"),
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
			assert.Equal(t,*tt.args.item.Action, *tt.want.Action)
			assert.Equal(t,*tt.args.item.Output, *tt.want.Output)
			assert.Equal(t,*tt.args.item.Region, *tt.want.Region)
			assert.Equal(t,*tt.args.item.Profile, *tt.want.Profile)
			assert.Equal(t,*tt.args.item.Query, *tt.want.Query)
		})
	}
}
