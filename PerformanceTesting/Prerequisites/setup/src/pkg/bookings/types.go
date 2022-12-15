package bookings

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
	TenantPassword   string
	BookingId        string
	ListingId        string
	ListingUuid      string
	StartDate        string
	EndDate          string
}

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
