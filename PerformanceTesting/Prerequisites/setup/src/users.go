package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-faker/faker/v4"
)

type exportUsers struct {
	Users []exportUser
}

type exportUser struct {
	Username string
	Password string
	Id       float64
	Uuid     string
}

type UserDetails struct {
	firstName string `faker:"first_name"`
	lastName  string `faker:"last_name"`
	password  string `faker:"password"`
	email     string `faker:"email"`
}

// registers a new user and returns export user object
func registerUser(client *http.Client) *exportUser {
	userDetails := map[string]string{
		"firstName": faker.FirstName(),
		"lastName":  faker.LastName(),
		"password":  faker.Password(),
		"email":     faker.Email(),
	}
	json_data, err := json.Marshal(userDetails)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Post(host+"/api/v2/user", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}
	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	p := exportUser{
		Username: res["email"].(string),
		Password: userDetails["password"],
		Id:       res["id"].(float64),
		Uuid:     res["uuid"].(string),
	}
	// user is not exportd after creation in listings so added it here
	// exportUsersList([]exportUser{p})
	return &p
}

// verifies the newly created user
func verifyUser(guid string, id string, client *http.Client) {

	phoneDetails := map[string]string{
		"countryCode": "NL",
		"phone":       "619377387",
	}
	json_data, err := json.Marshal(phoneDetails)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Post(host+"/api/v2/validation/phone?guid="+guid, "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		log.Fatal("/validate/phone request returned error")
	}

	userDetails := map[string]string{
		"landlordType":     "student",
		"phoneCountryCode": "NL",
		"phone":            "619377387",
	}
	json_data1, err := json.Marshal(userDetails)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPut, host+"/api/v2/user/"+id, bytes.NewBuffer(json_data1))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	resp1, _ := client.Do(req)
	var res map[string]interface{}
	json.NewDecoder(resp1.Body).Decode(&res)
	if resp1.StatusCode != 200 {
		var res map[string]interface{}
		json.NewDecoder(resp1.Body).Decode(&res)
		log.Fatal("/api/v2/user/ request returned error")
	}

	verificationCodePayload := map[string]string{
		"type": "sms",
	}
	json_data2, err := json.Marshal(verificationCodePayload)
	if err != nil {
		log.Fatal(err)
	}

	resp2, err := client.Post(host+"/api/v2/user/"+id+"/request-code?guid="+guid, "application/json",
		bytes.NewBuffer(json_data2))
	if err != nil {
		log.Fatal(err)
	}
	if resp2.StatusCode != 204 {
		log.Fatal("/user/<userid>/request-code request returned error", err, resp2.StatusCode)
	}
	var res2 map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&res2)

	resp3, err := client.Get(host + "/api/v2/tests/verification-code?guid=" + guid + "&phone=%2B31619377387")
	//resp3, err := client.Get(host+"/api/v2/tests/verification-code?guid=a3c38b11-33b2-11ed-938d-d902b977c503&phone=%2B31619377387")

	if err != nil {
		log.Fatal(err)
	}
	if resp3.StatusCode != 200 {
		var res map[string]interface{}
		json.NewDecoder(resp3.Body).Decode(&res)
		log.Fatal("/tests/verification-code request returned error")
	}
	var res3 map[string]interface{}
	json.NewDecoder(resp3.Body).Decode(&res3)
	code := res3["code"].(string)

	verifyUserPayload := map[string]string{
		"code": code,
	}
	json_data3, err := json.Marshal(verifyUserPayload)
	if err != nil {
		log.Fatal(err)
	}
	resp4, err := client.Post(host+"/api/v2/user/"+id+"/verify-code", "application/json",
		bytes.NewBuffer(json_data3))
	if err != nil {
		log.Fatal(err)
	}
	if resp4.StatusCode != 200 {
		log.Fatal("/user/verify-code request returned error")
	}

}
