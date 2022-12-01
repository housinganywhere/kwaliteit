package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"sync"
)

const MAX_CONCURRENCY = 50

// will be replaced by whatever values are passed from commandline
var host = "https://localhost:8080"
var exportLocation = "./testdata.json"

type exportBookingRequests struct {
	BookingRequests []exportBookingRequest
}

type exportBookingRequest struct {
	LandlordUsername string
	LandlordPassword string
	LandlordId       string
	LandlordUuid     string
	TenantId         string
	TenantUuid       string
	TenantUsername   string
	BookingId        string
	ListingId        string
	ListingUuid      string
	StartDate        string
	EndDate          string
}

type streets struct {
	Streets []street
}

type street struct {
	Name    string
	City    string
	Country string
}

type exportListings struct {
	Listings []exportListing
}

type exportListing struct {
	Id                 string
	Uuid               string
	AdvertiserId       string
	AdvertiserEmail    string
	AdvertiserPassword string
}

// exports users into json file
func exportUsersList(userDetails []exportUser) {
	usersList := exportUsers{
		Users: userDetails,
	}

	file, _ := json.MarshalIndent(usersList, "", " ")
	err := ioutil.WriteFile(exportLocation+"/Users_test.json", file, 0644)

	if err != nil {
		panic(err)
	}
}

// creates landlords and exports to json
func generateUsers(count int) {
	var userDetails []exportUser

	for i := 0; i < count; i++ {

		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatalf("Got error while creating cookie jar %s", err.Error())
		}
		client := http.Client{
			Jar: jar,
		}

		landlordDetails := registerUser(&client)
		verifyUser(landlordDetails.Uuid, strings.Split(fmt.Sprintf("%f", landlordDetails.Id), ".")[0], &client)

		user := exportUser{
			Username: landlordDetails.Username,
			Password: landlordDetails.Password,
			Id:       landlordDetails.Id,
			Uuid:     landlordDetails.Uuid,
		}

		userDetails = append(userDetails, user)
	}
	exportUsersList(userDetails)
}

func generateListings(count int) {
	listingsList := createMultipleListings(count)
	doExportListings(listingsList)
}

// creates listings and exports to json
func createMultipleListings(count int) []exportListing {
	maxConcurrency := MAX_CONCURRENCY

	if count < maxConcurrency {
		maxConcurrency = count - 1
	}

	listingDetailsJson, err := os.Open("../data/ListingDetails.json")
	if err != nil {
		fmt.Println(err)
	}
	defer listingDetailsJson.Close()

	listingDetailsInBytes, _ := ioutil.ReadAll(listingDetailsJson)

	streetsDataJson, err := os.Open("../data/Streets.json")
	if err != nil {
		fmt.Println(err)
	}

	defer streetsDataJson.Close()
	streetsDataBytesValue, _ := ioutil.ReadAll(streetsDataJson)
	var streetList streets

	json.Unmarshal(streetsDataBytesValue, &streetList)

	listingsToExport := make([]exportListing, count+1)

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	client := http.Client{
		Jar: jar,
	}

	landlordDetails := registerUser(&client)
	verifyUser(landlordDetails.Uuid, strings.Split(fmt.Sprintf("%f", landlordDetails.Id), ".")[0], &client)

	var wg sync.WaitGroup

	for i := 0; i < count; i++ {

		wg.Add(1)

		go func(idx int, listingToExport []exportListing) {
			fmt.Println(idx)
			defer wg.Done()

			var listingDetailsMap map[string]interface{}
			noOfStreets := len(streetList.Streets)
			selectedStreet := idx % noOfStreets

			json.Unmarshal([]byte(listingDetailsInBytes), &listingDetailsMap)
			listingDetailsMap["address"] = streetList.Streets[selectedStreet].Name + ", " + streetList.Streets[selectedStreet].City
			listingDetailsMap["street"] = streetList.Streets[selectedStreet].Name
			listingDetailsBytesArray, _ := json.Marshal(listingDetailsMap)

			listingDetails, err := createListing(&client, listingDetailsBytesArray)

			if err != nil {
				fmt.Println(err)
			}

			listing := exportListing{
				Id:                 listingDetails["listingId"],
				Uuid:               listingDetails["listingUuid"],
				AdvertiserId:       listingDetails["advertiserId"],
				AdvertiserEmail:    landlordDetails.Username,
				AdvertiserPassword: landlordDetails.Password,
			}

			listingsToExport[i] = listing
		}(i, listingsToExport)

		if i%maxConcurrency == 0 && i != 0 {
			wg.Wait()
		}
	}
	return listingsToExport
}

// exports users into json file
func doExportListings(listingsToExport []exportListing) {
	listingsList := exportListings{
		Listings: listingsToExport,
	}
	file, _ := json.MarshalIndent(listingsList, "", " ")
	//_ = ioutil.WriteFile("listings.json", file, 0644)
	_ = ioutil.WriteFile(exportLocation+"/Listings_test.json", file, 0644)
}

func generateListingsWithBookingRequest(count int) {
	var bookingRequests []exportBookingRequest

	listingsList := createMultipleListings(count)

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	client := http.Client{
		Jar: jar,
	}

	userDetails := registerUser(&client)

	for i := 0; i < count; i++ {
		doSubscription(int(userDetails.Id), "firstName lastName", &client)
		listingId, _ := strconv.Atoi(listingsList[i].Id)
		bookingDetails := sendBookingRequest(listingId, int(userDetails.Id), &client)
		bookingRequest := exportBookingRequest{
			LandlordUsername: listingsList[i].AdvertiserEmail,
			LandlordPassword: listingsList[i].AdvertiserPassword,
			LandlordId:       listingsList[i].AdvertiserId,
			LandlordUuid:     listingsList[i].Uuid,
			TenantId:         strings.Split(fmt.Sprintf("%f", userDetails.Id), ".")[0],
			TenantUuid:       userDetails.Uuid,
			TenantUsername:   userDetails.Username,
			BookingId:        bookingDetails["bookingId"],
			ListingId:        listingsList[i].Id,
			ListingUuid:      listingsList[i].Uuid,
			StartDate:        bookingDetails["startDate"],
			EndDate:          bookingDetails["endDate"],
		}
		bookingRequests = append(bookingRequests, bookingRequest)
	}

	doExportBookingRequests(bookingRequests)

}

// exports users into json file
func doExportBookingRequests(bookingRequests []exportBookingRequest) {
	bookingRequestsList := exportBookingRequests{
		BookingRequests: bookingRequests,
	}

	file, _ := json.MarshalIndent(bookingRequestsList, "", " ")
	//_ = ioutil.WriteFile("listings.json", file, 0644)
	_ = ioutil.WriteFile(exportLocation, file, 0644)
}

func main() {
	dataCategory := flag.String("dataCategory", DC_USER, "category of data to be created")
	noOfItems := flag.Int("count", 5, "number of records of data to be created")
	hostName := flag.String("host", "https://stage.housinganywhere.com", "host against which the data needs to be created")
	exportfile := flag.String("exportLocation", "./testdata.json", "location at which the generated test data is exported")

	flag.Parse()
	exportLocation = *exportfile
	host = *hostName

	switch *dataCategory {
	case DC_USER:
		generateUsers(*noOfItems)
	case DC_LISTING:
		generateListings(*noOfItems)
	case DC_BOOKING_REQUESTS:
		generateListingsWithBookingRequest(*noOfItems)
	}
}
