package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/bookings"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/listings"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/users"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/utilities"
)

func main() {
	var (
		noOfItems    int
		dataCategory string
		hostName     string
		exportFile   string
	)

	flag.IntVar(&noOfItems, "count", 5, "number of records of data to be created")
	flag.StringVar(&dataCategory, "dataCategory", "user", "category of data to be created")
	flag.StringVar(&hostName, "host", "https://stage.housinganywhere.com", "category of data to be created")
	flag.StringVar(&exportFile, "exportLocation", "./testdata.json", "location at which the generated test data is exported")

	_, err := url.ParseRequestURI(hostName)
	if err != nil {
		log.Fatalf("%s is not a valid url", hostName)
	}

	flag.Parse()

	datSeeder := utilities.NewDataSeeder(hostName, exportFile)
	if hostName == datSeeder.Config.ProductionUrl {
		log.Fatal("You cannot generate test data against prod environment")
	}
	switch dataCategory {
	case datSeeder.Config.DataCategoryUser:
		users.GenerateUsers(noOfItems, hostName, exportFile)
	case datSeeder.Config.DataCategoryListing:
		listings.CreateAndExportListings(&datSeeder, noOfItems)
	case datSeeder.Config.DataCategoryBookingRequest:
		bookings.GenerateListingsWithBookingRequest(&datSeeder, noOfItems)
	}
}
