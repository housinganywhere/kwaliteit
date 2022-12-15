package listings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/users"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/utilities"
)

func CreateAndExportListings(count int, hostName string, exportFile string) {
	listingsToExport := GenerateListingsWithUniqueLL(count, hostName, exportFile)
	utilities.ExportData(exportListings{
		Listings: listingsToExport,
	},
		exportFile)
}

func GenerateListingsWithUniqueLL(count int, hostName string, exportFile string) []*exportListing {
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
	streetCount := len(streetList.Streets)
	listingsToExport := make([]*exportListing, count)

	//go-routine batches -- each batch of 50
	batchCount := (count / 50) + 1

	for j := 0; j < batchCount; j++ {
		var it int
		var wg sync.WaitGroup
		if j == (batchCount - 1) {
			it = count % 50
		} else {
			it = 50
		}
		wg.Add(it)
		for i := 0; i < it; i++ {
			go func(i int) {
				defer wg.Done()
				client := utilities.NewDataSeeder(hostName, exportFile)
				jar, err := cookiejar.New(nil)
				if err != nil {
					log.Fatalf("Got error while creating cookie jar %s", err.Error())
				}
				client.HttpClient.Jar = jar
				user := users.RegisterUser(&client)
				users.VerifyUser(&client, user.Uuid, strings.Split(fmt.Sprintf("%f", user.Id), ".")[0])
				rand.Seed(time.Now().UnixNano())
				streetIndex := rand.Intn(streetCount)
				listingDetailsMap["address"] = streetList.Streets[streetIndex].Name + ", " + streetList.Streets[i].City
				listingDetailsMap["street"] = streetList.Streets[streetIndex].Name
				if err != nil {
					log.Fatalf("Error occurred while marshaling listingDetailsMap into Json:- %s", err.Error())
				}

				listing := CreateListing(&client, listingDetailsMap)

				listingsToExport[i+(50*j)] = &exportListing{
					Id:                 listing["listingId"],
					Uuid:               listing["listingUuid"],
					AdvertiserId:       listing["advertiserId"],
					AdvertiserEmail:    user.Username,
					AdvertiserPassword: user.Password,
				}
			}(i)
		}
		wg.Wait()
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
