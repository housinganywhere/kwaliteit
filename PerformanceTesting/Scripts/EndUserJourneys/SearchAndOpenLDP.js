import { check } from "k6";
import http from "k6/http";
import { group } from "k6";
import { randomItem } from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { SharedArray } from "k6/data";
import * as constants from "../Common/Constants.js";
import ApplicationFlow from '../Common/ApplicationFlow.js';
import Utilities from "../Common/Utilities.js";

const listings = new SharedArray("listingList", function () {
  return JSON.parse(open("../TestData/listings.json")).Listings;
});

export function searchListingAndOpenLDPJourney() {
  let listing = randomItem(listings);
  let listingId = listing.Id;
  let listingUuid = listing.Uuid;
  let advertiserId = listing.AdvertiserId;
  let date = new Date();
  const utilities = new Utilities();
  date.setMonth(date.getMonth() + 1);
  let startDate = date.toLocaleDateString("sv");
  const appFlows = new ApplicationFlow();
  let customHeaders = Object.assign({
    Origin: constants.baseUrl,
    "Content-Type": "application/json"
  }, constants.haApiHeader);

  appFlows.visitHomePageAndSearchListing(startDate);

  group("open listing details page", function () {
    utilities.get(
      `/room/${listingId}/nl/Rotterdam/spoorsingel?startDate=${startDate}`,
      constants.haApiHeader,
      200,
      `/room/<listingId>/nl/Rotterdam/<location>?startDate=<startDate>`
    );

    utilities.get(
      `/api/v2/user/me/information?expand=impersonation`,
      constants.haApiHeader,
      401,
      `/api/v2/user/me/information?expand=impersonation`
    );

    utilities.get(
      `/api/v2/qualification/requirements/by-listing/${listingUuid}`,
      constants.haApiHeader,
      404,
      `/api/v2/qualification/requirements/by-listing/<listingUuid>`
    );

    utilities.get(
      `/api/v2/currency/EUR/rates`,
      constants.haApiHeader,
      200,
      `/api/v2/currency/EUR/rates`
    );

    utilities.get(
      `/api/v2/qualification/requirements/by-listing/${listingUuid}`,
      constants.haApiHeader,
      404,
      `/api/v2/qualification/requirements/by-listing/<listingUuid>`
    );

    utilities.get(
      `/api/v2/listing/${listingId}?expand=photos%2Cexclusions%2CbookablePeriods%2Cvideos%2Ccosts&lang=en`,
      constants.haApiHeader,
      200,
      `/api/v2/listing/<listingId>?expand=photos%2Cexclusions%2CbookablePeriods%2Cvideos%2Ccosts&lang=en`
    );
    
    utilities.get(
      `/api/v2/user/${advertiserId}`,
      constants.haApiHeader,
      200,
      `/api/v2/user/<userId>`
    );

    utilities.get(
      `/api/v2/user/${advertiserId}/information?expand=listingsCount%2CsavedPaymentMethodsCount%2CadvertiserBookingsCount%2CisKYCDone`,
      constants.haApiHeader,
      200,
      `/api/v2/user/<userId>/information?expand=listingsCount%2CsavedPaymentMethodsCount%2CadvertiserBookingsCount%2CisKYCDone`
    );
    
    const rentalConditionsEstimateBody = JSON.stringify({
      source: { listingId: parseInt(listingId) },
      startDate: "2022-10-01",
      endDate: "2023-10-31",
    });
    utilities.post(
      `/api/v2/bookings/rental-conditions/estimate`,
      customHeaders,
      rentalConditionsEstimateBody,
      200,
      `/api/v2/bookings/rental-conditions/estimate`
    );
    
    utilities.get(
      `/api/v2/conversations/user/me?listingId=${listingId}&limit=1`,
      constants.haApiHeader,
      401,
      `/api/v2/conversations/user/me?listingId=<listingId>&limit=1`
    );

    utilities.get(
      `/api/v2/billing/subscriptions/products/tenant`,
      constants.haApiHeader,
      200,
      `/api/v2/billing/subscriptions/products/tenant`
    );
  });
}
