package utilities

import (
	"io/ioutil"
	"log"
	"os"

	jsoniter "github.com/json-iterator/go"
)

// exports users into json file
func ExportData(jsonPayload any, location string) {
	//file, err := json.MarshalIndent(jsonPayload, "", " ")
	file, err := jsoniter.MarshalIndent(jsonPayload, "", " ")
	if err != nil {
		log.Fatal("error occured while indenting json data for exporting. The error is :- " + err.Error())
	}
	err = ioutil.WriteFile(location, file, 0644)
	if err != nil {
		log.Fatal("error occured while exporting data. The error is :- " + err.Error())
	}
}

// exists returns whether the given file or directory exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
