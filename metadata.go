package awclip

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type ApiCallProviderName string

const (
	ApiCallProviderNameAws ApiCallProviderName = "aws cli python"
	ApiCallProviderNameGo  ApiCallProviderName = "go sdk v2"
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
	Parameters map[string]*string
	Query      *string
}

// CacheEntry
// Cmd executable command line e.g. aws ec2 describe instances
type CacheEntry struct {
	Id            *string
	Cmd           *string
	Created       time.Time
	LastAccessed  time.Time
	AccessCounter int
	Parameters    Parameters
	Provider      string
}

func WriteMetadata(md *CacheEntry) error {
	location := GetLocationMetaData(md.Id)

	content, err := json.Marshal(md)
	if err != nil {
		log.Fatal(err)
	}
	var fileoptions int
	if CacheHit(md.Id) {
		// Update file
		fileoptions = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	} else {
		fileoptions = os.O_RDWR | os.O_CREATE
	}
	file, err := os.OpenFile(*location, fileoptions, 0755)
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString(string(content))
	defer file.Close()

	if err != nil {
		log.Panicf("failed write file: %s, %s", *location, err)
		return err
	}

	return nil
}

func ReadMetaData(id *string) (*CacheEntry, error) {
	file, err := os.Open(*GetLocationMetaData(id))
	if err != nil {
		log.Panicf("failed open file: %s", err)
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panicf("failed reading file: %s", err)
		return nil, err
	}
	var metadata CacheEntry
	json.Unmarshal(data, &metadata)
	if err != nil {
		log.Printf("failed Unmarshal file: %s", err)
		return nil, err
	}

	return &metadata, nil
}

func UpdateMetaData(md *CacheEntry) error {
	location := GetLocationMetaData(md.Id)
	md.AccessCounter += 1
	content, err := json.Marshal(md)
	if err != nil {
		log.Fatal(err)
	}
	var fileoptions = os.O_RDWR | os.O_CREATE | os.O_TRUNC

	file, err := os.OpenFile(*location, fileoptions, 0755)
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString(string(content))
	defer file.Close()

	if err != nil {
		log.Panicf("failed write file: %s, %s", *location, err)
		return err
	}

	return nil
}

func (item *CacheEntry) ArgumentsToCachedEntry(args []string) {
	service := args[1]
	*item.Parameters.Service = service
	action := args[2]
	*item.Parameters.Action = action
	for i, arg := range args {
		found := false
		if strings.HasPrefix(arg, "--") {
			if arg == "--query" {
				*item.Parameters.Query = args[i+1]
				found = true
			}
			if arg == "--region" {
				*item.Parameters.Region = args[i+1]
				found = true
			}
			if arg == "--profile" {
				*item.Parameters.Profile = args[i+1]
				found = true
			}
			if arg == "--output" {
				*item.Parameters.Output = args[i+1]
				found = true
			}
			if !found {
				// --action
				action := strings.Split(arg, "--")[1] 
				item.Parameters.Parameters[action] = aws.String(args[i+1])
			}
		}
	}
}

func (a *Parameters) AlmostEqual(b *Parameters) bool {
	if *a.Service == *b.Service &&
		*a.Action == *b.Action &&
		*a.Output == *b.Output &&
		*a.Query == *b.Query {
		return true
	}
	return false
}

func (a *Parameters) AlmostEqualWithParameters(b *Parameters) bool {
	if *a.Service == *b.Service &&
		*a.Action == *b.Action &&
		*a.Output == *b.Output &&
		*a.Query == *b.Query {
		if len(a.Parameters) != len(b.Parameters) {
			return false
		}
		// All Parameters exist?
		for key := range a.Parameters {
			value, ok := b.Parameters[key]
			if !ok {
				return false
			}
			if !(*value == "*") {
				if !(*a.Parameters[key] == *value) {
					return false
				}
			}
		}
		return true
	}
	return false
}
