package awclip_test

import (
	"testing"

	"gotest.tools/assert"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/megaproaktiv/awclip"
	"github.com/megaproaktiv/awclip/services"
)

func TestArgumentsToCachedEntry(t *testing.T) {
	type args struct {
		args []string
		item *awclip.CacheEntry
	}
	newCacheEntry1 := &awclip.CacheEntry{
		Parameters: awclip.Parameters{
			Service: new(string),
			Action:  new(string),
			Output:  new(string),
			Region:  new(string),
			Profile: new(string),
			Parameters: map[string]*string{},
			Query:   new(string),
		},
	}
	newCacheEntry2 := &awclip.CacheEntry{
		Parameters: awclip.Parameters{
			Service: new(string),
			Action:  new(string),
			Output:  new(string),
			Region:  new(string),
			Profile: new(string),
			Parameters: map[string]*string{},
			Query:   new(string),
		},
	}
	newCacheEntry3 := &awclip.CacheEntry{
		Parameters: awclip.Parameters{
			Service: new(string),
			Action:  new(string),
			Output:  new(string),
			Region:  new(string),
			Profile: new(string),
			Parameters: map[string]*string{},
			Query:   new(string),
		},
	}
	newCacheEntry4 := &awclip.CacheEntry{
		Parameters: awclip.Parameters{
			Service: new(string),
			Action:  new(string),
			Output:  new(string),
			Region:  new(string),
			Profile: new(string),
			Parameters: map[string]*string{},
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
				item: newCacheEntry1,
			},
			want: &awclip.Parameters{
				Service: aws.String("ec2"),
				Action:  aws.String("describe-instances"),
				Output:  aws.String("text"),
				Region:  aws.String("eu-central-1"),
				Profile: aws.String("myprofile"),
				Parameters: map[string]*string{},
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
				item: newCacheEntry2,
			},
			want: &awclip.Parameters{
				Service: aws.String("ec2"),
				Action:  aws.String("describe-regions"),
				Output:  aws.String("text"),
				Region:  aws.String("eu-central-1"),
				Profile: aws.String("helmut"),
				Parameters: map[string]*string{},
				Query:   aws.String("Regions[].RegionName"),
			},
		},
		{
			name: "recognise iam User Parameter",
			args: args{
				args: []string{
					"dist/awclip",
					"iam",
					"list-user-policies",
					"--profile",
					"helmut",
					"--user-name",
					"johndonkey",
					"--region",
					"eu-central-1",
					"--output",
					"text",
				},
				item: newCacheEntry3,
			},
			want: &awclip.Parameters{
				Service:    aws.String("iam"),
				Action:     aws.String("list-user-policies"),
				Output:     aws.String("text"),
				Region:     aws.String("eu-central-1"),
				Profile:    aws.String("helmut"),
				Query:      aws.String(""),
				Parameters: map[string]*string{"user-name": aws.String("johndonkey")},
			},
		},
		{
			name: "recognise list-attached-user-policies",
			args: args{
				args: []string{
					"dist/awclip",
					"iam",
					"list-attached-user-policies",
					"--profile",
					"helmut",
					"--user-name",
					"johndonkey",
					"--region",
					"eu-central-1",
					"--output",
					"text",
				},
				item: newCacheEntry4,
			},
			want: &awclip.Parameters{
				Service:    aws.String("iam"),
				Action:     aws.String("list-attached-user-policies"),
				Output:     aws.String("text"),
				Region:     aws.String("eu-central-1"),
				Profile:    aws.String("helmut"),
				Query:      aws.String(""),
				Parameters: map[string]*string{"user-name": aws.String("johndonkey")},
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


func TestAlmostEqualWithParameters(t *testing.T) {

	newParms := &awclip.Parameters{
		Service:    aws.String("iam"),
		Action:     aws.String("list-attached-user-policies"),
		Output:     aws.String("text"),
		Region:     aws.String("eu-central-1"),
		Profile:    aws.String("helmut"),
		Query:      aws.String(""),
		Parameters: map[string]*string{"user-name": aws.String("johndonkey")},
	}


	ok := newParms.AlmostEqualWithParameters(services.IamListAttachedUserPoliciesParameter)
	assert.Equal(t,true, ok ,"IamListAttachedUserPoliciesParameter should be matched")
}
