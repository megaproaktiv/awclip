package awclip_test

import (
	"github.com/megaproaktiv/awclip"
	"testing"
	"gotest.tools/assert"
)

func TestDiscriminatedCommand(t *testing.T) {
	args := []string{
		"dist/awclip",
		"iam",
		"generate-credential-report",
		"--output",
		"text",
		"--profile",
		"helmut",
		"--query",
		"state",
		"--region",
		"eu-central-1",
	}
	assert.Equal(t, true, awclip.DiscriminatedCommand(&args[1], &args[2]))
	assert.Equal(t, false, awclip.DiscriminatedCommand(&args[1], &args[3]))
}
