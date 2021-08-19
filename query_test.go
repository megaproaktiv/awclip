package awclip_test

import (
	"os"
	"testing"
	"gotest.tools/assert"
	"github.com/megaproaktiv/awclip"

)

func TestQueryText(t *testing.T) {

	jsonLambdaListB, err := os.ReadFile("services/testdata/lambda/lambda_listfunctions_eu-west-1.json")
	if err != nil {
		t.Error("Cant read input testdata")
		t.Error(err)
	}
	jsonLambdaList := string(jsonLambdaListB)

		
	textLambdaListB, err := os.ReadFile("services/testdata/lambda/lambda_listfunctions_eu-west-1_cli.txt")
	if err != nil {
		t.Error("Cant read output testdata")
		t.Error(err)
	}
	textLambdaList := string(textLambdaListB)
	
	queryLambdaList := "Functions[*].{R:Runtime,N:FunctionName}"

	type args struct {
		json  *string
		query *string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "Lambda List text query",
			args: args{
				json:  &jsonLambdaList,
				query: &queryLambdaList,
			},
			want: &textLambdaList,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := awclip.QueryText(tt.args.json, tt.args.query)
			assert.DeepEqual(t,  *tt.want, *got)
			
		})
	}
}

func TestQueryTextSmall(t *testing.T) {

	jsonLambdaListB, err := os.ReadFile("services/testdata/lambda/small.json")
	if err != nil {
		t.Error("Cant read input testdata")
		t.Error(err)
	}
	jsonLambdaList := string(jsonLambdaListB)

	
	textLambdaListB, err := os.ReadFile("services/testdata/lambda/small.txt")
	if err != nil {
		t.Error("Cant read output testdata")
		t.Error(err)
	}
	textLambdaList := string(textLambdaListB)
	
	queryLambdaList := "Functions[*].{R:Runtime,N:FunctionName}"

	type args struct {
		json  *string
		query *string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "Lambda List text query",
			args: args{
				json:  &jsonLambdaList,
				query: &queryLambdaList,
			},
			want: &textLambdaList,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := awclip.QueryText(tt.args.json, tt.args.query)
			assert.DeepEqual(t,  *tt.want, *got)
			
		})
	}
}

func TestQueryEc2DescribeFunc(t *testing.T) {

	jsonDataB, err := os.ReadFile("services/testdata/ec2/ec2_describeinstances_eu-central-1.json")
	if err != nil {
		t.Error("Cant read input testdata")
		t.Error(err)
	}
	jsonData := string(jsonDataB)

	
	textDataB, err := os.ReadFile("services/testdata/ec2/ec2_describeinstances_eu-central-1_cli.txt")
	if err != nil {
		t.Error("Cant read output testdata")
		t.Error(err)
	}
	textData := string(textDataB)
	
	queryString := "Reservations[*].Instances[*].[InstanceId]"

	type args struct {
		json  *string
		query *string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "Lambda List text query",
			args: args{
				json:  &jsonData,
				query: &queryString,
			},
			want: &textData,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := awclip.QueryText(tt.args.json, tt.args.query)
			assert.DeepEqual(t,  *tt.want, *got)
			
		})
	}
}
