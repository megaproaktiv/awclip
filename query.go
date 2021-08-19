package awclip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/megaproaktiv/awclip/services"
	"github.com/megaproaktiv/awclip/cache"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jmespath/go-jmespath"
)

const TAB = "\t"
const NL = "\n"
// Query transforms aws json result with the query string (jemspath) as text output
func QueryText(input *string, query *string) *string{
	keys  := Orderkeys(query)
	m := map[string]interface{}{}
	//Parsing/Unmarshalling JSON encoding/json
	err := json.Unmarshal([]byte(*input), &m)
	if err != nil {
		panic(err)
	}
	result, _ := jmespath.Search(*query, m)
	if debug {
		fmt.Println("Result: \n",result)
	}
	
	var buffer bytes.Buffer
	anArray := result.([]interface{})
	parseArray(anArray,keys,&buffer)

	text := buffer.String()
	//fmt.Println("Text: \n",text)
	return &text
}

// Order Problem
// https://golang.org/ref/spec#For_statements
// The iteration order over maps is not specified and is not guaranteed to be the same from one iteration to the next.

func parseMap(aMap map[string]interface{}, keys *[]*string,buffy *bytes.Buffer) {
	first := true
	for _,key := range *keys {
		val := aMap[*key]
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			buffy.WriteString(val.(string))
			parseMap(val.(map[string]interface{}), keys, buffy)
		case []interface{}:
			buffy.WriteString(val.(string))
			parseArray(val.([]interface{}),keys, buffy)
		default:
			if first {
				first = false
			}else{
				buffy.WriteString(TAB)
			}
			value := concreteVal.(string)
			all := buffy.String()
			_ = all
			buffy.WriteString(value)
		}
	}
	buffy.WriteString(NL)
}

func parseArray(anArray []interface{} ,keys *[]*string,buffy *bytes.Buffer) {
	for _, val := range anArray {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			parseMap(val.(map[string]interface{}),keys , buffy)
		case []interface{}:
			parseArray(val.([]interface{}),keys ,buffy)
		default:
			buffy.WriteString(concreteVal.(string))
		}
	}
}


// the AWS CLI got the order wrong
// Example:
// aws lambda list-functions --output text --query "Functions[*].{R:Runtime,N:FunctionName}"
// ask-custom-messeguide-default	nodejs8.10
// amplify-login-custom-message-80c9d7a7	nodejs12.x
// Demo2	nodejs14.x
//
// 
func Orderkeys(query *string) *[]*string{
	lambdaFunctionsKeys := []*string{
		aws.String("N"),
		aws.String("R"),
		
	}
	if *query == "Functions[*].{R:Runtime,N:FunctionName}"{
		return  &lambdaFunctionsKeys
	}
	return nil
}

// Call Query
// read prefetched json data and apply query
func CallQuery(metadata *cache.CacheEntry) *string{
	//	Query
	prefetchName := services.ApiCallDumpFileNameString(metadata.Parameters.Service,metadata.Parameters.Action,metadata.Parameters.Region)
	jsondata, err := ioutil.ReadFile(*prefetchName)
	if err != nil {
		panic("read error, " + err.Error())
	}
    data := string(jsondata)
	// to text
	return QueryText(&data, metadata.Parameters.Query)
}