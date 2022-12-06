package utilities

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// exports users into json file
func ExportData(jsonPayload any, location string) {
	file, err := json.MarshalIndent(jsonPayload, "", " ")
	if err != nil {
		log.Fatal("error occured while indenting json data for exporting. The error is :- " + err.Error())
	}
	err = ioutil.WriteFile(location, file, 0644)
	if err != nil {
		log.Fatal("error occured while exporting data. The error is :- " + err.Error())
	}
}
