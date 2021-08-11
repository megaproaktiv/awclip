package awclip

import (
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
	AdditionalParameters map[string]*string
	Query      *string
}

// CacheEntry
// Cmd executable command line e.g. aws ec2 describe instances
//go:generate ffjson metadata.go

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

	content, err := md.MarshalJSON()
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
	// json.Unmarshal(data, &metadata)
	metadata.UnmarshalJSON(data)
	if err != nil {
		log.Printf("failed Unmarshal file: %s", err)
		return nil, err
	}

	return &metadata, nil
}

func UpdateMetaData(md *CacheEntry) error {
	location := GetLocationMetaData(md.Id)
	md.AccessCounter += 1
	// content, err := json.Marshal(md)
	content, err := md.MarshalJSON()
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

func (Parameters *Parameters) ArgumentsToCachedEntry(args []string) {
	service := args[1]
	Parameters.Service = &service
	action := args[2]
	Parameters.Action = &action
	if Parameters.AdditionalParameters == nil {
		Parameters.AdditionalParameters = make(map[string]*string)
	}
	for i, arg := range args {
		found := false
		if strings.HasPrefix(arg, "--") {
			if arg == "--query" {
				Parameters.Query = &args[i+1]
				found = true
			}
			if arg == "--region" {
				Parameters.Region = &args[i+1]
				found = true
			}
			if arg == "--profile" {
				Parameters.Profile = &args[i+1]
				found = true
			}
			if arg == "--output" {
				Parameters.Output = &args[i+1]
				found = true
			}
			if !found {
				// --action
				action := strings.Split(arg, "--")[1] 
				if len(args) > (i+1)  {
					Parameters.AdditionalParameters[action] = aws.String(args[i+1])
				}
			}
		}
	}
	Parameters.Query = emptyWhenNil(Parameters.Query)
	Parameters.Region = emptyWhenNil(Parameters.Region)
	Parameters.Profile = emptyWhenNil(Parameters.Profile)
	Parameters.Output = emptyWhenNil(Parameters.Output)
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
		if len(a.AdditionalParameters) != len(b.AdditionalParameters) {
			return false
		}
		// All Parameters exist?
		for key := range a.AdditionalParameters {
			value, ok := b.AdditionalParameters[key]
			if !ok {
				return false
			}
			if !(*value == "*") {
				if !(*a.AdditionalParameters[key] == *value) {
					return false
				}
			}
		}
		return true
	}
	return false
}

func (d *CacheEntry) Copy() *CacheEntry{
	if d.Cmd == nil {
		panic("cmd of Cache Entry must be nonzero")
	}
	newD := &CacheEntry{
		Id:            aws.String(*d.Id),
		Cmd:           aws.String(*d.Cmd),
		Created:       time.Now(),
		LastAccessed:  time.Now(),
		AccessCounter: 0,
		Parameters:    *d.Parameters.Copy(),
		Provider:      "",
	}
	return newD
}

func (p *Parameters) Copy() *Parameters{
	addOns := make(map[string]*string)
	for index, element  := range p.AdditionalParameters{        
		addOns[index] = element
   }
	newP := &Parameters{
		Service:    emptyWhenNil(p.Service),
		Action:     emptyWhenNil(p.Action),
		Output:     emptyWhenNil(p.Output),
		Region:     emptyWhenNil(p.Region),
		Profile:    emptyWhenNil(p.Profile),
		AdditionalParameters: map[string]*string{},
		Query:      emptyWhenNil(p.Query),
	}
	return newP
}