package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
)

type DataSeeder struct {
	Host           string
	ExportLocation string
	HttpClient     *http.Client
}

func NewDataSeeder(host string, exportLocation string) DataSeeder {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	return DataSeeder{
		Host:           strings.TrimRight(host, "/"),
		ExportLocation: exportLocation,
		HttpClient: &http.Client{
			Jar: jar,
		},
	}
}

func Get(client *DataSeeder, path string, expectedStatus int) map[string]interface{} {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", client.Host, path), nil)
	if err != nil {

		log.Fatalf("error occurred while framing request for path:- %s. The error is:- %s", path, err.Error())
	}
	resp, err := client.HttpClient.Do(req)
	if err != nil {
		rq, _ := httputil.DumpRequest(req, true)
		rs, _ := httputil.DumpResponse(resp, true)
		fmt.Println("The request was :- " + string(rq))
		fmt.Println("The response was :- " + string(rs))
		log.Fatalf("The call to %s returned an error. The error is:- %s", path, err)
	}
	if resp.StatusCode != expectedStatus {
		rq, _ := httputil.DumpRequest(req, true)
		rs, _ := httputil.DumpResponse(resp, true)
		fmt.Println("The request was :- " + string(rq))
		fmt.Println("The response was :- " + string(rs))
		log.Fatalf("The call to %s returned an incorrect status code. Expected code was %d but actual is %d", path, expectedStatus, resp.StatusCode)
	}
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	return res
}

func Post(client *DataSeeder, path string, payload any, headers *http.Header, expectedStatus int) map[string]interface{} {
	resp := makeUpdateCall(http.MethodPost, client, path, payload, headers, expectedStatus)
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	return res
}

func Put(client *DataSeeder, path string, payload any, headers *http.Header, expectedStatus int) map[string]interface{} {
	resp := makeUpdateCall(http.MethodPut, client, path, payload, headers, expectedStatus)
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	return res
}

func makeUpdateCall(callMethod string, client *DataSeeder, path string, payload any, headers *http.Header, expectedStatus int) *http.Response {
	if headers.Get("Content-Type") == "" {
		headers.Add("Content-Type", "application/json")
	}

	var formBody io.Reader
	var res *http.Response

	if headers.Get("Content-Type") == "application/json" {
		requestBodyBytes, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("Error occurred while marshalling payload to json %s", payload)
		}
		formBody = bytes.NewBuffer(requestBodyBytes)
	}
	if headers.Get("Content-Type") == "application/x-www-form-urlencoded" {
		urlEncodedBody, ok := payload.(url.Values)
		if ok {
			formBody = strings.NewReader(urlEncodedBody.Encode())
		} else {
			log.Fatal("The payload was not in correct form url encoded format")
		}
	}

	req, err := http.NewRequest(callMethod, fmt.Sprintf("%s%s", client.Host, path), formBody)
	if err != nil {
		log.Fatalf("error occurred while framing request for path:- %s. The error is:- %s", path, err)
	}
	req.Header = *headers
	res, err = client.HttpClient.Do(req)
	if err != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(req.Body)
		rs, _ := httputil.DumpResponse(res, true)
		fmt.Println("The request was :- " + buf.String())
		fmt.Println("The response was :- " + string(rs))
		log.Fatalf("The call to %s%s returned an error. The error is:- %s", client.Host, path, err)
	}
	if res.StatusCode != expectedStatus {
		buf := new(bytes.Buffer)
		buf.ReadFrom(formBody)
		rs, _ := httputil.DumpResponse(res, true)
		fmt.Println("The request was :- " + buf.String())
		fmt.Println("The response was :- " + string(rs))
		log.Fatalf("The call to %s%s returned an incorrect status code. Expected code was %d but actual is %d", client.Host, path, expectedStatus, res.StatusCode)

	}

	return res
}
