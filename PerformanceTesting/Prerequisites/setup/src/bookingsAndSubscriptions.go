package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type StripeSubscriptionPayload struct {
	Type         string `json:"type"`
	CardNumber   string `json:"card[number]"`
	CardCvc      string `json:"card[cvc]"`
	CardExpMonth int    `json:"card[exp_month]"`
	CardExpYear  int    `json:"card[exp_year]"`
	Key          string `json:"key"`
}

type HASubscriptionPayload struct {
	StripePaymentMethodId string   `json:"stripePaymentMethodId"`
	Currency              string   `json:"currency"`
	PriceCodes            []string `json:"priceCodes"`
	BillingCountry        string   `json:"billingCountry"`
	CustomerName          string   `json:"customerName"`
}

type ConversationPayload struct {
	MessageText string `json:"messageText"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	BookingLink bool   `json:"bookingLink"`
	ListingId   int    `json:"listingId"`
}

type PaymentSettelementPayload struct {
	ExternalReference     string        `json:"externalReference"`
	ExternalReferenceType string        `json:"externalReferenceType"`
	PurchaseVal           PurchaseValue `json:"purchaseValue"`
	SellerAccountUuid     string        `json:"sellerAccountUuid"`
	BuyerAccountUuid      string        `json:"buyerAccountUuid"`
}

type PayInMethodPayload struct {
	MethodType string `json:"methodType"`
	Currency   string `json:"currency"`
}

type PurchaseValue struct {
	Cents    int    `json:"cents"`
	Currency string `json:"currency"`
}

func doSubscription(tenantId int, tenantName string, client *http.Client) {

	stripeBody := url.Values{}
	stripeBody.Set("type", "card")
	stripeBody.Set("card[number]", STRIPE_TEST_CARD_NUMBER)
	stripeBody.Set("card[cvc]", STRIPE_TEST_CARD_CVC)
	stripeBody.Set("card[exp_month]", STRIPE_TEST_CARD_EXPIRY_MONTH)
	stripeBody.Set("card[exp_year]", STRIPE_TEST_CARD_EXPIRY_YEAR)
	stripeBody.Set("key", STRIPE_TEST_KEY)
	encodedData := stripeBody.Encode()

	req, err := http.NewRequest(http.MethodPost, "https://api.stripe.com/v1/payment_methods", strings.NewReader(encodedData))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "https://js.stripe.com/")
	req.Header.Add("Origin", "https://js.stripe.com")
	req.Header.Add("Host", "api.stripe.com")
	if err != nil {
		log.Fatal(err)
	}
	resp, _ := client.Do(req)
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	stripePaymentMethodId := res["id"]

	subscriptionBody := &HASubscriptionPayload{
		StripePaymentMethodId: stripePaymentMethodId.(string),
		Currency:              "eur",
		PriceCodes:            []string{"tenant-subscription-month"},
		BillingCountry:        "NL",
		CustomerName:          tenantName,
	}
	subscriptionBuffer := new(bytes.Buffer)
	json.NewEncoder(subscriptionBuffer).Encode(subscriptionBody)
	req1, err := http.NewRequest(http.MethodPost, host+"/api/v2/billing/subscription", subscriptionBuffer)
	req1.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	resp1, _ := client.Do(req1)
	var res1 map[string]interface{}
	json.NewDecoder(resp1.Body).Decode(&res1)

}

func sendBookingRequest(listingId int, tenantId int, client *http.Client) map[string]string {
	conversationBody := &ConversationPayload{
		MessageText: "Hello! I'm interested in renting your accommodation. I believe I match your tenant preferences. \nPlease get back as soon as possible",
		StartDate:   time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
		EndDate:     time.Now().AddDate(0, 6, 0).Format("2006-01-02"),
		BookingLink: false,
		ListingId:   listingId,
	}

	conversationBodyBuffer := new(bytes.Buffer)
	json.NewEncoder(conversationBodyBuffer).Encode(conversationBody)
	req, err := http.NewRequest(http.MethodPost, host+"/api/v2/conversation", conversationBodyBuffer)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	resp, _ := client.Do(req)
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println("The conversation response is")
	fmt.Println(res)
	conversationId := strings.Split(fmt.Sprintf("%f", res["id"]), ".")[0]

	resp2, err := client.Get(host + "/api/v2/conversation/" + conversationId + "?expand=advertiser%2Clisting%2CrentalConditions")
	if err != nil {
		log.Fatal(err)
	}
	if resp2.StatusCode != 200 {
		log.Fatal("/api/v2/conversation/ request returned error")
	}
	var res2 map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&res2)

	listing := res2["listing"].(map[string]interface{})
	price, _ := listing["price"].(float64)
	advertiserDetails := res2["advertiser"].(map[string]interface{})
	landlordUuid := advertiserDetails["uuid"].(string)
	rentalConditions := res2["rentalConditions"].(map[string]interface{})
	tenantUuid := rentalConditions["tenantUuid"].(string)
	externalReferenceId := res2["uuid"].(string)

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

	paymentSettelementBuffer := new(bytes.Buffer)
	json.NewEncoder(paymentSettelementBuffer).Encode(paymentSettelementPayload)
	req3, err := http.NewRequest(http.MethodPost, host+"/api/v2/payment/settlement", paymentSettelementBuffer)
	req3.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	resp3, _ := client.Do(req3)
	var res3 map[string]interface{}
	json.NewDecoder(resp3.Body).Decode(&res3)

	payInMethodPayload := &PayInMethodPayload{
		MethodType: "stripe-card",
		Currency:   "EUR",
	}
	payInMethodBuffer := new(bytes.Buffer)
	json.NewEncoder(payInMethodBuffer).Encode(payInMethodPayload)
	req4, err := http.NewRequest(http.MethodPost, host+"/api/v2/payment/method/pay-in", payInMethodBuffer)
	req4.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	resp4, _ := client.Do(req4)
	var res4 map[string]interface{}
	json.NewDecoder(resp4.Body).Decode(&res4)
	payInId := res4["uuid"].(string)
	//customerUuid := res3["customerUuid"]
	nextSteps := res4["nextSteps"].([]interface{})
	firstStep, _ := nextSteps[0].(map[string]interface{})
	stripeToken := firstStep["token1"].(string)
	stripeIntent := strings.Split(stripeToken, "_secret_")[0]

	stripeBody := url.Values{}
	stripeBody.Set("payment_method_data[type]", "card")
	stripeBody.Set("payment_method_data[card][number]", STRIPE_TEST_CARD_NUMBER)
	stripeBody.Set("payment_method_data[card][cvc]", STRIPE_TEST_CARD_CVC)
	stripeBody.Set("payment_method_data[card][exp_month]", STRIPE_TEST_CARD_EXPIRY_MONTH)
	stripeBody.Set("payment_method_data[card][exp_year]", STRIPE_TEST_CARD_EXPIRY_YEAR)
	stripeBody.Set("expected_payment_method_type", "card")
	stripeBody.Set("use_stripe_sdk", "true")
	stripeBody.Set("key", STRIPE_TEST_KEY)
	stripeBody.Set("client_secret", stripeToken)
	encodedData := stripeBody.Encode()

	req5, err := http.NewRequest(http.MethodPost, "https://api.stripe.com/v1/setup_intents/"+stripeIntent+"/confirm", strings.NewReader(encodedData))
	req5.Header.Add("content-type", "application/x-www-form-urlencoded")
	req5.Header.Add("Referer", "https://js.stripe.com/")
	req5.Header.Add("Origin", "https://js.stripe.com")
	req5.Header.Add("Host", "api.stripe.com")
	if err != nil {
		log.Fatal(err)
	}
	resp5, _ := client.Do(req5)
	var res5 map[string]interface{}
	json.NewDecoder(resp5.Body).Decode(&res5)
	payInToken := res5["payment_method"].(string)

	payInMethodInfo := map[string]string{
		"token1": payInToken,
	}
	payInMethodInfoJson, err := json.Marshal(payInMethodInfo)
	if err != nil {
		log.Fatal(err)
	}

	req6, err := http.NewRequest(http.MethodPut, host+"/api/v2/payment/method/pay-in/"+payInId, bytes.NewBuffer(payInMethodInfoJson))
	req6.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	resp6, _ := client.Do(req6)
	var res6 map[string]interface{}
	json.NewDecoder(resp6.Body).Decode(&res6)
	paymentMethodRef := res6["uuid"].(string)

	paymentMethodInfo := map[string]string{
		"paymentMethodReference": paymentMethodRef,
		"paymentType":            "paymentengine",
	}
	paymentMethodInfoJson, err := json.Marshal(paymentMethodInfo)
	if err != nil {
		log.Fatal(err)
	}

	req7, err := http.NewRequest(http.MethodPut, host+"/api/v2/conversation/"+conversationId+"/booking-request", bytes.NewBuffer(paymentMethodInfoJson))
	req7.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	resp7, _ := client.Do(req7)
	var res7 map[string]interface{}
	json.NewDecoder(resp7.Body).Decode(&res7)

	bookingDetails := map[string]string{}
	bookingDetails["bookingId"] = conversationId
	bookingDetails["startDate"] = conversationBody.StartDate
	bookingDetails["endDate"] = conversationBody.EndDate
	return bookingDetails
}
