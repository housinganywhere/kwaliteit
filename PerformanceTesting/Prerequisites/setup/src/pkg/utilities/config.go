package utilities

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	StripeTestCardNumber       string `json:"stripeTestCardNumber"`
	StripeTestCardExpiryMonth  string `json:"stripeTestCardExpiryMonth"`
	StripeTestCardExpiryYear   string `json:"stripeTestCardExpiryYear"`
	StripeTestCardCvc          string `json:"stripeTestCardCvc"`
	StripeTestKey              string `json:"stripeTestKey"`
	DataCategoryListing        string `json:"dataCategoryListing"`
	DataCategoryUser           string `json:"dataCategoryUser"`
	DataCategoryBookingRequest string `json:"dataCategoryBookingRequest"`
	ProductionUrl              string `json:"productionUrl"`
	StripeApi                  string `json:"stripeApi"`
	StripeJs                   string `json:"stripeJs"`
}

func loadConfiguration() *config {
	var setupConfig *config
	configFile, err := os.Open("../data/config.json")
	if err != nil {
		log.Fatalf("Error occurred while reading the config file:- %s", err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&setupConfig)
	return setupConfig
}
