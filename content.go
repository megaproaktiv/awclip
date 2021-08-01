package awclip

import (
	"io/ioutil"
	"log"
	"os"
)

func ReadContent(id *string) (*string, error) {
	file, err := os.Open(*GetLocationData(id))
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
	content := string(data)

	return &content, nil
}

func WriteContent(id *string, content *string) error {
	location := GetLocationData(id)
	file, err := os.OpenFile(*location, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Panicf("failed open file: %s, %s", *location, err)
		return err
	}
	defer file.Close()
	file.WriteString(*content)
	if err != nil {
		log.Panicf("failed write file: %s, %s", *location, err)
		return err
	}

	return nil
}
