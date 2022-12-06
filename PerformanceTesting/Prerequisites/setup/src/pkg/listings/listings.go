package listings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/users"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/utilities"
)

func CreateAndExportListings(client *utilities.DataSeeder, count int) {
	listingsToExport := GenerateListingsWithUniqueLL(client, count)
	utilities.ExportData(exportListings{
		Listings: listingsToExport,
	},
		client.ExportLocation)
}

func GenerateListingsWithUniqueLL(client *utilities.DataSeeder, count int) []exportListing {
	listingDetailsJson, err := os.Open("../data/listingDetails.json")
	if err != nil {
		log.Fatal("error occurred file opening the listing details json file:- " + err.Error())
	}
	defer listingDetailsJson.Close()

	listingDetailsInBytes, err := ioutil.ReadAll(listingDetailsJson)
	if err != nil {
		log.Fatal("error occurred file reading the listing details json file:- " + err.Error())
	}

	var listingDetailsMap map[string]interface{}
	json.Unmarshal([]byte(listingDetailsInBytes), &listingDetailsMap)

	streetsDataJson, err := os.Open("../data/streets.json")
	if err != nil {
		log.Fatal("error occurred file opening the streets json file:- " + err.Error())
	}
	defer streetsDataJson.Close()

	streetsDataBytesValue, err := ioutil.ReadAll(streetsDataJson)
	if err != nil {
		log.Fatal("error occurred file reading the streets json file:- " + err.Error())
	}

	var streetList streets
	json.Unmarshal(streetsDataBytesValue, &streetList)

	var listingsToExport []exportListing

	for i := 0; i < count; i++ {
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatalf("Got error while creating cookie jar %s", err.Error())
		}
		client.HttpClient.Jar = jar
		user := users.RegisterUser(client)
		users.VerifyUser(client, user.Uuid, strings.Split(fmt.Sprintf("%f", user.Id), ".")[0])

		listingDetailsMap["address"] = streetList.Streets[i].Name + ", " + streetList.Streets[i].City
		listingDetailsMap["street"] = streetList.Streets[i].Name
		if err != nil {
			log.Fatalf("Error occurred while marshaling listingDetailsMap into Json:- %s", err.Error())
		}

		listing := CreateListing(client, listingDetailsMap)

		listingsToExport = append(listingsToExport, exportListing{
			Id:                 listing["listingId"],
			Uuid:               listing["listingUuid"],
			AdvertiserId:       listing["advertiserId"],
			AdvertiserEmail:    user.Username,
			AdvertiserPassword: user.Password,
		})

	}
	return listingsToExport
}

func CreateListing(client *utilities.DataSeeder, listingDetails map[string]interface{}) map[string]string {
	res := utilities.Post(client, "/api/v2/listing", listingDetails, &http.Header{}, http.StatusOK)
	listingId := strings.Split(res["listingPath"].(string), "/")[2]

	utilities.Put(client, fmt.Sprintf("/api/v2/listing/%s/publish", listingId), map[string]string{}, &http.Header{}, http.StatusNoContent)

	ld := map[string]string{}
	ld["listingId"] = listingId
	ld["listingUuid"] = res["uuid"].(string)
	ld["advertiserId"] = strings.Split(fmt.Sprintf("%f", res["advertiserId"]), ".")[0]

	return ld
}
