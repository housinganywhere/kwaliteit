package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/bookings"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/config"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/listings"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/users"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/utilities"
)

func main() {
	var (
		noOfItems      int
		dataCategory   string
		hostName       string
		exportLocation string
	)

	flag.IntVar(&noOfItems, "count", 5, "number of records of data to be created")
	flag.StringVar(&dataCategory, "dataCategory", "user", "which category of data to be created, options available are:- listing, user & bookingRequests. Defaults to user")
	flag.StringVar(&hostName, "host", "https://stage.housinganywhere.com", "environment on which the data needs to be created, e.g. https://stage.housinganywhere.com")
	flag.StringVar(&exportLocation, "exportLocation", ".", "location at which the generated test data is to be exported exported. Should be relative or absolute path to the directory where the data file needs to be stored")

	_, err := url.ParseRequestURI(hostName)
	if err != nil {
		log.Fatalf("%s is not a valid url", hostName)
	}

	flag.Parse()

	isLocationPresent, err := utilities.PathExists(exportLocation)
	if err != nil {
		log.Fatalf("The export location provided %s has some errors %s", exportLocation, err)
	}
	if !isLocationPresent {
		log.Fatalf("The export location provided %s does not exist", exportLocation)
	}

	if hostName == config.ProductionUrl {
		log.Fatal("You cannot generate test data against prod environment")
	}
	switch dataCategory {
	case config.DataCategoryUser:
		users.GenerateUsers(noOfItems, hostName, exportLocation+"/users.json")
	case config.DataCategoryListing:
		listings.CreateAndExportListings(noOfItems, hostName, exportLocation+"/listings.json")
	case config.DataCategoryBookingRequest:
		if noOfItems > 15 {
			log.Fatal("You cannot request for more than 15 listings with booking requests at a time. This is to respect the rate limiting in place for stripe test mode")
		} else {
			bookings.GenerateListingsWithBookingRequest(noOfItems, hostName, exportLocation+"/bookingRequests.json")
		}
	default:
		log.Printf("Unsupported dataCategory option %s. Supported options are listing, users, bookingRequest", dataCategory)
	}
}
