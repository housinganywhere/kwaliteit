package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func createListing(client *http.Client, listingDetailsBytesArray []byte) (map[string]string, error) {
	path := "/api/v2/listing"

	resp, err := client.Post(fmt.Sprintf("%s%s", host, path), "application/json",
		bytes.NewBuffer(listingDetailsBytesArray))

	if err != nil {
		return nil, fmt.Errorf("error: %v, path: %s, statusCode: %d", err, path, resp.StatusCode)
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("error: %v, path: %s, statusCode: %d", string(body), path, resp.StatusCode)
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	listingId := strings.Split(res["listingPath"].(string), "/")[2]

	path = fmt.Sprintf("%s/%s%s", path, listingId, "/publish")

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s%s", host, path), bytes.NewBuffer([]byte{}))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	resp1, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var res1 map[string]interface{}
	json.NewDecoder(resp1.Body).Decode(&res1)

	if resp1.StatusCode != 204 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("error: %v, path: %s, statusCode: %d", string(body), path, resp.StatusCode)

	}
	listingDetails := map[string]string{}
	listingDetails["listingId"] = listingId
	listingDetails["listingUuid"] = res["uuid"].(string)
	listingDetails["advertiserId"] = strings.Split(fmt.Sprintf("%f", res["advertiserId"]), ".")[0]

	return listingDetails, nil
}
