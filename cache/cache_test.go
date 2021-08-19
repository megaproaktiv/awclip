package cache_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/megaproaktiv/awclip/services"
	"github.com/megaproaktiv/awclip/cache"
	"gotest.tools/assert"
)

func TestArgumentsToCachedEntry(t *testing.T) {
	type args struct {
		args []string
		item *cache.CacheEntry
	}
	newCacheEntry1 := &cache.CacheEntry{
		Parameters: cache.Parameters{
			Service:    new(string),
			Action:     new(string),
			Output:     new(string),
			Region:     new(string),
			Profile:    new(string),
			AdditionalParameters: map[string]*string{},
			Query:      new(string),
		},
	}
	newCacheEntry2 := &cache.CacheEntry{
		Parameters: cache.Parameters{
			Service:    new(string),
			Action:     new(string),
			Output:     new(string),
			Region:     new(string),
			Profile:    new(string),
			AdditionalParameters: map[string]*string{},
			Query:      new(string),
		},
	}
	newCacheEntry3 := &cache.CacheEntry{
		Parameters: cache.Parameters{
			Service:    new(string),
			Action:     new(string),
			Output:     new(string),
			Region:     new(string),
			Profile:    new(string),
			AdditionalParameters: map[string]*string{},
			Query:      new(string),
		},
	}
	newCacheEntry4 := &cache.CacheEntry{
		Parameters: cache.Parameters{
			Service:    new(string),
			Action:     new(string),
			Output:     new(string),
			Region:     new(string),
			Profile:    new(string),
			AdditionalParameters: map[string]*string{},
			Query:      new(string),
		},
	}

	tests := []struct {
		name string
		args args
		want *cache.Parameters
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
			want: &cache.Parameters{
				Service:    aws.String("ec2"),
				Action:     aws.String("describe-instances"),
				Output:     aws.String("text"),
				Region:     aws.String("eu-central-1"),
				Profile:    aws.String("myprofile"),
				AdditionalParameters: map[string]*string{},
				Query:      aws.String("Reservations[*].Instances[*].[InstanceId]"),
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
			want: &cache.Parameters{
				Service:    aws.String("ec2"),
				Action:     aws.String("describe-regions"),
				Output:     aws.String("text"),
				Region:     aws.String("eu-central-1"),
				Profile:    aws.String("helmut"),
				AdditionalParameters: map[string]*string{},
				Query:      aws.String("Regions[].RegionName"),
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
			want: &cache.Parameters{
				Service:    aws.String("iam"),
				Action:     aws.String("list-user-policies"),
				Output:     aws.String("text"),
				Region:     aws.String("eu-central-1"),
				Profile:    aws.String("helmut"),
				Query:      aws.String(""),
				AdditionalParameters: map[string]*string{"user-name": aws.String("johndonkey")},
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
			want: &cache.Parameters{
				Service:    aws.String("iam"),
				Action:     aws.String("list-attached-user-policies"),
				Output:     aws.String("text"),
				Region:     aws.String("eu-central-1"),
				Profile:    aws.String("helmut"),
				Query:      aws.String(""),
				AdditionalParameters: map[string]*string{"user-name": aws.String("johndonkey")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			tt.args.item.Parameters.ArgumentsToCachedEntry(tt.args.args)
			assert.DeepEqual(t, tt.args.item.Parameters, *tt.want)

		})
	}
}

func TestAlmostEqualWithParameters(t *testing.T) {

	newParms := &cache.Parameters{
		Service:    aws.String("iam"),
		Action:     aws.String("list-attached-user-policies"),
		Output:     aws.String("text"),
		Region:     aws.String("eu-central-1"),
		Profile:    aws.String("helmut"),
		Query:      aws.String(""),
		AdditionalParameters: map[string]*string{"user-name": aws.String("johndonkey")},
	}

	ok := newParms.AlmostEqualWithParameters(services.IamListAttachedUserPoliciesParameter)
	assert.Equal(t, true, ok, "IamListAttachedUserPoliciesParameter should be matched")
}

func TestCacheEntry_ArgumentsToCachedEntry(t *testing.T) {
	type fields struct {
		Id            *string
		Cmd           *string
		Created       time.Time
		LastAccessed  time.Time
		AccessCounter int
		Parameters    cache.Parameters
		Provider      string
	}
	type args struct {
		args []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := &cache.CacheEntry{
				Id:            tt.fields.Id,
				Cmd:           tt.fields.Cmd,
				Created:       tt.fields.Created,
				LastAccessed:  tt.fields.LastAccessed,
				AccessCounter: tt.fields.AccessCounter,
				Parameters:    tt.fields.Parameters,
				Provider:      tt.fields.Provider,
			}
			metadata.Parameters.ArgumentsToCachedEntry(tt.args.args)
		})
	}
}


// Test for error:
// dist/awclip ec2 describe-regions --query "Regions[].RegionName" --output text --profile ggtrcadmin --region eu-central-1 --region-names
// panic: assignment to entry in nil map

// goroutine 1 [running]:
// github.com/megaproaktiv/awclip.(*CacheEntry).ArgumentsToCachedEntry(0xc00012be58, 0xc0000200c0, 0xc, 0xc)
// 	/Users/silberkopf/letsbuild/awclip/metadata.go:154 +0x3ee
// main.main()
// 	/Users/silberkopf/letsbuild/awclip/main/main.go:31 +0x11d
// Prowler adds "--region-names" to aws cli ?
func TestArgumentsToCachedEntryRegionsnames(t *testing.T){
	metadata := &cache.CacheEntry{}
	var args =  []string{
		"dist/awclip",
		"ec2",
		"list-attached-user-policies",
		"--profile",
		"helmut",
		"--query",
		"Regions[].RegionName",
		"--region",
		"eu-central-1",
		"--output",
		"text",
		"--region-names",
	}

	metadata.Parameters.ArgumentsToCachedEntry(args)
	assert.Equal(t , "eu-central-1", *metadata.Parameters.Region)
}