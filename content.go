package awclip

import (
	"io/ioutil"
	"log"
	"os"
	"github.com/megaproaktiv/awclip/cache"

)

func ReadContent(id *string) (*string, error) {
	file, err := os.Open(*cache.GetLocationData(id))
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

func ReadContentUpdate(id *string) (*string, error) {

	content, err := ReadContent(id)
	if err != nil {
		return nil, err
	}

	metadata, _ := cache.ReadMetaData(id)
	err = cache.UpdateMetaData(metadata)
	if err != nil {
		log.Print(err)
	}
	return content, nil
}

func WriteContent(id *string, content *string) error {
	location := cache.GetLocationData(id)
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
