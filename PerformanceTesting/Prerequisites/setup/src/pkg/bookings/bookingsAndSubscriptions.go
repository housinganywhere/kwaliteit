package bookings

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/config"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/listings"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/users"
	"github.com/housinganywhere/kwaliteit/performance-testing/setup/src/pkg/utilities"
)

func GenerateListingsWithBookingRequest(count int, hostName string, exportFile string) {
	listingsList := listings.GenerateListingsWithUniqueLL(count, hostName, exportFile)

	bookingRequests := make([]exportBookingRequest, count)

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
			fmt.Println("The listing id of listing # " + strconv.Itoa(i+1) + "is " + listingsList[i].Id)
			listingId, _ := strconv.Atoi(listingsList[i].Id)
			go func(i int) {
				defer wg.Done()
				client := utilities.NewDataSeeder(hostName, exportFile)
				jar, err := cookiejar.New(nil)
				if err != nil {
					log.Fatalf("Got error while creating cookie jar %s", err.Error())
				}
				client.HttpClient.Jar = jar
				userDetails := users.RegisterUser(&client)
				doSubscription(&client, int(userDetails.Id), "firstName lastName")

				bookingDetails := sendBookingRequest(&client, listingId, int(userDetails.Id))
				bookingRequest := exportBookingRequest{
					LandlordUsername: listingsList[i].AdvertiserEmail,
					LandlordPassword: listingsList[i].AdvertiserPassword,
					LandlordId:       listingsList[i].AdvertiserId,
					LandlordUuid:     listingsList[i].Uuid,
					TenantId:         strings.Split(fmt.Sprintf("%f", userDetails.Id), ".")[0],
					TenantUuid:       userDetails.Uuid,
					TenantUsername:   userDetails.Username,
					TenantPassword:   userDetails.Password,
					BookingId:        bookingDetails["bookingId"],
					ListingId:        listingsList[i].Id,
					ListingUuid:      listingsList[i].Uuid,
					StartDate:        bookingDetails["startDate"],
					EndDate:          bookingDetails["endDate"],
				}
				bookingRequests[i+(50*j)] = bookingRequest
			}(i)
		}
		wg.Wait()
	}
	utilities.ExportData(exportBookingRequests{
		BookingRequests: bookingRequests,
	}, exportFile)
}

func doSubscription(client *utilities.DataSeeder, tenantId int, tenantName string) {

	stripeBody := url.Values{}
	stripeBody.Set("type", "card")
	stripeBody.Set("card[number]", config.StripeTestCardNumber)
	stripeBody.Set("card[cvc]", config.StripeTestCardCvc)
	stripeBody.Set("card[exp_month]", config.StripeTestCardExpiryMonth)
	stripeBody.Set("card[exp_year]", config.StripeTestCardExpiryYear)
	stripeBody.Set("key", config.StripeTestKey)

	headers := http.Header{}
	headers.Add("content-type", "application/x-www-form-urlencoded")
	headers.Add("Referer", config.StripeJs+"/")
	headers.Add("Origin", config.StripeJs)
	headers.Add("Host", config.StripeApi)
	haHost := client.Host
	client.Host = config.StripeApi

	res := utilities.Post(client, "/v1/payment_methods", stripeBody, &headers, http.StatusOK)
	stripePaymentMethodId := res["id"]

	subscriptionBody := &HASubscriptionPayload{
		StripePaymentMethodId: stripePaymentMethodId.(string),
		Currency:              "eur",
		PriceCodes:            []string{"tenant-subscription-month"},
		BillingCountry:        "NL",
		CustomerName:          tenantName,
	}
	client.Host = haHost
	utilities.Post(client, "/api/v2/billing/subscription", subscriptionBody, &http.Header{}, http.StatusOK)
}

func sendBookingRequest(client *utilities.DataSeeder, listingId int, tenantId int) map[string]string {
	conversationBody := &ConversationPayload{
		MessageText: "Hello! I'm interested in renting your accommodation. I believe I match your tenant preferences. \nPlease get back as soon as possible",
		StartDate:   time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
		EndDate:     time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
		BookingLink: false,
		ListingId:   listingId,
	}
	res := utilities.Post(client, "/api/v2/conversation", conversationBody, &http.Header{}, http.StatusOK)
	conversationId := strings.Split(fmt.Sprintf("%f", res["id"]), ".")[0]
	params := url.Values{}
	params.Add("expand", "advertiser,listing,rentalConditions")
	res1 := utilities.Get(client, fmt.Sprintf("/api/v2/conversation/%s?%s", conversationId, params.Encode()), http.StatusOK)
	listing := res1["listing"].(map[string]interface{})
	price, _ := listing["price"].(float64)
	advertiserDetails := res1["advertiser"].(map[string]interface{})
	landlordUuid := advertiserDetails["uuid"].(string)
	rentalConditions := res1["rentalConditions"].(map[string]interface{})
	tenantUuid := rentalConditions["tenantUuid"].(string)
	externalReferenceId := res1["uuid"].(string)

	purchaseValue := &PurchaseValue{
		Cents:    int(price),
		Currency: "EUR",
	}

	paymentSettelementPayload := &PaymentSettelementPayload{
		ExternalReference:     externalReferenceId,
		ExternalReferenceType: "booking",
		PurchaseVal:           *purchaseValue,
		SellerAccountUuid:     landlordUuid,
		BuyerAccountUuid:      tenantUuid,
	}

	utilities.Post(client, "/api/v2/payment/settlement", paymentSettelementPayload, &http.Header{}, http.StatusOK)

	payInMethodPayload := &PayInMethodPayload{
		MethodType: "stripe-card",
		Currency:   "EUR",
	}
	res2 := utilities.Post(client, "/api/v2/payment/method/pay-in", payInMethodPayload, &http.Header{}, http.StatusOK)
	payInId := res2["uuid"].(string)
	nextSteps := res2["nextSteps"].([]interface{})
	firstStep, _ := nextSteps[0].(map[string]interface{})
	stripeToken := firstStep["token1"].(string)
	stripeIntent := strings.Split(stripeToken, "_secret_")[0]

	stripeBody := url.Values{}
	stripeBody.Set("payment_method_data[type]", "card")
	stripeBody.Set("payment_method_data[card][number]", config.StripeTestCardNumber)
	stripeBody.Set("payment_method_data[card][cvc]", config.StripeTestCardCvc)
	stripeBody.Set("payment_method_data[card][exp_month]", config.StripeTestCardExpiryMonth)
	stripeBody.Set("payment_method_data[card][exp_year]", config.StripeTestCardExpiryYear)
	stripeBody.Set("expected_payment_method_type", "card")
	stripeBody.Set("use_stripe_sdk", "true")
	stripeBody.Set("key", config.StripeTestKey)
	stripeBody.Set("client_secret", stripeToken)

	headers := http.Header{}
	headers.Add("content-type", "application/x-www-form-urlencoded")
	headers.Add("Referer", config.StripeJs+"/")
	headers.Add("Origin", config.StripeJs)
	headers.Add("Host", config.StripeApi)
	haHost := client.Host
	client.Host = config.StripeApi

	res3 := utilities.Post(client, fmt.Sprintf("/v1/setup_intents/%s/confirm", stripeIntent), stripeBody, &headers, http.StatusOK)
	payInToken := res3["payment_method"].(string)

	client.Host = haHost
	payInMethodInfo := map[string]string{
		"token1": payInToken,
	}
	res4 := utilities.Put(client, fmt.Sprintf("/api/v2/payment/method/pay-in/%s", payInId), payInMethodInfo, &http.Header{}, http.StatusOK)
	paymentMethodRef := res4["uuid"].(string)

	paymentMethodInfo := map[string]string{
		"paymentMethodReference": paymentMethodRef,
		"paymentType":            "paymentengine",
	}

	utilities.Put(client, fmt.Sprintf("/api/v2/conversation/%s/booking-request", conversationId), paymentMethodInfo, &http.Header{}, http.StatusOK)

	bookingDetails := map[string]string{}
	bookingDetails["bookingId"] = conversationId
	bookingDetails["startDate"] = conversationBody.StartDate
	bookingDetails["endDate"] = conversationBody.EndDate
	return bookingDetails
}
