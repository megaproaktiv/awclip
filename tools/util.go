package tools

import "bytes"

var empty = ""
const TAB = "\t"
const NL = "\n"

func EmptyWhenNil(s *string) *string {
	if s == nil {
		return &empty
	}
	return s
}



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
			ParseArray(val.([]interface{}),keys, buffy)
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

func ParseArray(anArray []interface{} ,keys *[]*string,buffy *bytes.Buffer) {
	for _, val := range anArray {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			parseMap(val.(map[string]interface{}),keys , buffy)
		case []interface{}:
			ParseArray(val.([]interface{}),keys ,buffy)
		default:
			buffy.WriteString(concreteVal.(string))
		}
	}
}
