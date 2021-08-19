package cache

import (
	"time"
)

type ApiCallProviderName string

const (
	ApiCallProviderNameAws ApiCallProviderName = "aws cli python"
	ApiCallProviderNameGo  ApiCallProviderName = "go sdk v2"
	DATADIR = ".awclip"

)

type ApiCallProvider struct {
	// Open for extensions like provided calls
	Name ApiCallProviderName
}

type Parameters struct {
	Service    *string
	Action     *string
	Output     *string
	Region     *string
	Profile    *string
	AdditionalParameters map[string]*string
	Query      *string
}

// CacheEntry
// Cmd executable command line e.g. aws ec2 describe instances
//go:generate ffjson cache_struct.go

type CacheEntry struct {
	Id            *string
	Cmd           *string
	Created       time.Time
	LastAccessed  time.Time
	AccessCounter int
	Parameters    Parameters
	Provider      string
}
