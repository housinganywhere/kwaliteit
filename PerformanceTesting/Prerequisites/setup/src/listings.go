package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func createListing(client *http.Client, listingDetailsBytesArray []byte) map[string]string {
	resp, err := client.Post(host+"/api/v2/listing", "application/json",
		bytes.NewBuffer(listingDetailsBytesArray))
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("/api/v2/listing request returned error")
	}
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	listingId := strings.Split(res["listingPath"].(string), "/")[2]

	req, err := http.NewRequest(http.MethodPut, host+"/api/v2/listing/"+listingId+"/publish", bytes.NewBuffer([]byte{}))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	resp1, _ := client.Do(req)

	var res1 map[string]interface{}
	json.NewDecoder(resp1.Body).Decode(&res1)
	if resp1.StatusCode != 204 {
		log.Fatal("/api/v2/user/ request returned error")
	}
	listingDetails := map[string]string{}
	listingDetails["listingId"] = listingId
	listingDetails["listingUuid"] = res["uuid"].(string)
	listingDetails["advertiserId"] = strings.Split(fmt.Sprintf("%f", res["advertiserId"]), ".")[0]

	return listingDetails
}
