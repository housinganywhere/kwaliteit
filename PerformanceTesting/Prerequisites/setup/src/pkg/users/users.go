package users

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/go-faker/faker/v4"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/utilities"
)

// creates landlords and exports to json
func GenerateUsers(client *utilities.DataSeeder, count int) {
	var userDetails []exportUser

	for i := 0; i < count; i++ {
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatalf("Got error while creating cookie jar %s", err.Error())
		}
		client.HttpClient.Jar = jar
		landlordDetails := RegisterUser(client)
		VerifyUser(client, landlordDetails.Uuid, strings.Split(fmt.Sprintf("%f", landlordDetails.Id), ".")[0])

		user := exportUser{
			Username: landlordDetails.Username,
			Password: landlordDetails.Password,
			Id:       landlordDetails.Id,
			Uuid:     landlordDetails.Uuid,
		}

		userDetails = append(userDetails, user)
	}
	utilities.ExportData(userDetails, client.ExportLocation)
}

// registers a new user and returns export user object
func RegisterUser(client *utilities.DataSeeder) *exportUser {
	userDetails := map[string]string{
		"firstName": faker.FirstName(),
		"lastName":  faker.LastName(),
		"password":  faker.Password(),
		"email":     faker.Email(),
	}

	res := utilities.Post(client, "/api/v2/user", userDetails, &http.Header{}, http.StatusOK)

	eu := exportUser{
		Username: res["email"].(string),
		Password: userDetails["password"],
		Id:       res["id"].(float64),
		Uuid:     res["uuid"].(string),
	}
	return &eu
}

// verifies the newly created user
func VerifyUser(client *utilities.DataSeeder, guid string, id string) {

	phoneDetails := map[string]string{
		"countryCode": "NL",
		"phone":       "619377387",
	}

	utilities.Post(client, fmt.Sprintf("/api/v2/validation/phone?guid=%s", guid), phoneDetails, &http.Header{}, http.StatusOK)

	userDetails := map[string]string{
		"landlordType":     "student",
		"phoneCountryCode": "NL",
		"phone":            "619377387",
	}

	utilities.Put(client, fmt.Sprintf("/api/v2/user/%s", id), userDetails, &http.Header{}, http.StatusOK)

	verificationCodePayload := map[string]string{
		"type": "sms",
	}

	utilities.Post(client, fmt.Sprintf("/api/v2/user/%s/request-code?guid=%s", id, guid), verificationCodePayload, &http.Header{}, http.StatusNoContent)

	params := url.Values{}
	params.Add("guid", guid)
	params.Add("phone", "+31619377387")
	res := utilities.Get(client, fmt.Sprintf("/api/v2/tests/verification-code?%s", params.Encode()), http.StatusOK)
	code := res["code"].(string)

	verifyUserPayload := map[string]string{
		"code": code,
	}

	utilities.Post(client, fmt.Sprintf("/api/v2/user/%s/verify-code", id), verifyUserPayload, &http.Header{}, http.StatusOK)

}
