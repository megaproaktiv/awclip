package awclip

import (
	"reflect"
	"testing"
)

func TestCleanUp(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "GetCallerIdentity",
			args: args{
				args: []string{
					"get-caller-identity",
					"--profile",
					"myprofile",
					"--region",
					"eu-central-1",
				},
			},
			want: []string{
				"get-caller-identity",
				"--region",
				"eu-central-1",
				"--profile",
				"myprofile",
			},
		},
		
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanUp(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CleanUp() = %v, want %v", got, tt.want)
			}
		})
	}
}