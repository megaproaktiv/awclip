package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/megaproaktiv/awclip/tools"
)


func WriteMetadata(md *CacheEntry) error {
	location := GetLocationMetaData(md.Id)
	var err error
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
	var err error
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
	Parameters.Query = tools.EmptyWhenNil(Parameters.Query)
	Parameters.Region = tools.EmptyWhenNil(Parameters.Region)
	Parameters.Profile = tools.EmptyWhenNil(Parameters.Profile)
	Parameters.Output = tools.EmptyWhenNil(Parameters.Output)
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

func (d *CacheEntry) Print(){
	fmt.Println("Id    : ",*d.Id)
	fmt.Println("Cmd   : ",*d.Cmd)
	d.Parameters.Print()
}


func (p *Parameters) Copy() *Parameters{
	addOns := make(map[string]*string)
	for index, element  := range p.AdditionalParameters{        
		addOns[index] = element
   }
	newP := &Parameters{
		Service:    tools.EmptyWhenNil(p.Service),
		Action:     tools.EmptyWhenNil(p.Action),
		Output:     tools.EmptyWhenNil(p.Output),
		Region:     tools.EmptyWhenNil(p.Region),
		Profile:    tools.EmptyWhenNil(p.Profile),
		AdditionalParameters: map[string]*string{},
		Query:      tools.EmptyWhenNil(p.Query),
	}
	return newP
}

func (p *Parameters) Print(){
	fmt.Println("Service: ",*p.Service)
	fmt.Println("Action : ",*p.Action)
	for key, element := range p.AdditionalParameters {
		fmt.Println(key, " : ", *element)
	}
}


func GetLocationData(contentId *string) *string {
	location := DATADIR + string(os.PathSeparator) + *contentId + ".json"
	return &location
}
func GetLocationMetaData(contentId *string) *string {
	location := DATADIR + string(os.PathSeparator) + *contentId + "-db.json"
	return &location
}

const SEP = " "
var empty = ""

func (parms *Parameters) CommandLine() *string {
	commandLine := "aws "
	
	commandLine += *tools.EmptyWhenNil(parms.Service) + SEP + 
	*tools.EmptyWhenNil(parms.Action) + SEP + 
	*tools.EmptyWhenNil(parms.Region) + SEP +
	*tools.EmptyWhenNil(parms.Output) + SEP +
	*tools.EmptyWhenNil(parms.Query)
	for key, element := range parms.AdditionalParameters {
		commandLine += SEP + key + SEP + *element
	}
	return &commandLine
}

// HashValue calculates the id of the chache entries
func (parms *Parameters) HashValue() *string {
	commandLine := parms.CommandLine()
	hash := md5.Sum([]byte(*commandLine))
	hashstring := hex.EncodeToString(hash[:])
	return &hashstring
	
}


func CacheMiss(id *string) bool {
	return !CacheHit(id)
}

// CacheHit - true if a file with contant is already there
func CacheHit(id *string) bool {
	location := GetLocationMetaData(id)
	info, err := os.Stat(*location)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}




