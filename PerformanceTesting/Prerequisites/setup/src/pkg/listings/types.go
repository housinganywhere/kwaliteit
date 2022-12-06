package listings

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
