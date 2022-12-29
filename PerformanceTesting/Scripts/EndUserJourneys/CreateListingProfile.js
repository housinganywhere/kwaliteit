import { check } from "k6";
import http from "k6/http";
import { group } from "k6";
import { SharedArray } from "k6/data";
import ApplicationFlow from "../Common/ApplicationFlow.js";
import * as constants from "../Common/Constants.js";
import Utilities from "../Common/Utilities.js";

const binaryFile = open("../TestData/room.png", "b");
const users = new SharedArray("userCredentials", function () {
  return JSON.parse(open("../TestData/users.json")).Users;
});

export function createListingJourney() {
  const user = users[Math.floor(Math.random() * users.length)];
  const appFlows = new ApplicationFlow();
  const utilities = new Utilities();
  let customHeaders = Object.assign(
    {
      Origin: constants.baseUrl,
      "Content-Type": "application/json",
    },
    constants.haApiHeader
  );
  
  appFlows.visitHomePageAndLogin(user.Username, user.Password);

  group("Open List Your Place Page", function () {
    utilities.get(
      `/list-your-place`,
      constants.haApiHeader,
      200,
      `/list-your-place`
    );
    utilities.get(
      `/api/v2/user/me/information?expand=impersonation`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/information?expand=impersonation`
    );
    utilities.get(
      `/api/v2/listings/configuration?expand=countries`,
      constants.haApiHeader,
      200,
      `/api/v2/listings/configuration?expand=countries`
    );
  });

  //to be used in other groups as well
  //hence declared outside of the group
  let listingId;
  let advertiserId;
  let listingUuid;
  let listingPath;
  let listingIcsCalendarPath;

  group("Create Listing Step 1", function () {
    const createListingStep1Payload = JSON.stringify(
      constants.listingBasicInfo
    );
    let res = utilities.post(
      "/api/v2/listing",
      customHeaders,
      createListingStep1Payload,
      200,
      "/api/v2/listing"
    );
    listingId = res.json().id;
    advertiserId = res.json().advertiserId;
    listingUuid = res.json().uuid;
    listingPath = res.json().listingPath;
    listingIcsCalendarPath = res.json().icsCalendarPath;

    let date = new Date();
    let exclusionEndDate = utilities.formatDate(
      new Date(date.getFullYear(), date.getMonth() + 1, 0)
    );
    const exclusionPayload = JSON.stringify({
      dates: [{ from: "0001-01-01", to: exclusionEndDate }],
      adjustDates: true,
    });
    res = utilities.post(
      `/api/v2/listing/${listingId}/exclusions`,
      customHeaders,
      exclusionPayload,
      200,
      `/api/v2/listing/<listingId>/exclusions`
    );

    utilities.get(
      `/my/listings/${listingId}/edit-draft`,
      constants.haApiHeader,
      200,
      `/my/listings/<listingId>/edit-draft`
    );

    utilities.get(
      `/api/v2/user/me/information?expand=impersonation`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/information?expand=impersonation`
    );
  });

  group("Create Listing Step 2", function () {
    utilities.get(
      `/api/v2/listings/configuration?expand=countries`,
      constants.haApiHeader,
      200,
      `/api/v2/listings/configuration?expand=countries`
    );
    utilities.get(
      `/api/v2/user/me/realtime`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/realtime`
    );
    utilities.get(
      `/api/v2/listing/${listingId}?expand=videos%2CbookablePeriods%2Ccosts`,
      constants.haApiHeader,
      200,
      `/api/v2/listing/<listingId>?expand=videos%2CbookablePeriods%2Ccosts`
    );
    utilities.get(
      `/api/v2/qualification/requirements/by-listing/${listingUuid}`,
      constants.haApiHeader,
      404,
      `/api/v2/listing/<listingUuid>?expand=videos%2CbookablePeriods%2Ccosts`
    );
    utilities.get(
      `/api/v2/qualification/requirements/defaults`,
      constants.haApiHeader,
      404,
      `/api/v2/qualification/requirements/defaults`
    );

    const createListingStep2Payload = JSON.stringify({
      facilities: constants.facilities,
      freePlaces: 2,
      description: "This is a test description out here. Please go over it",
    });
    let res = utilities.put(
      `/api/v2/listing/${listingId}?expand=videos%2CbookablePeriods%2Ccosts`,
      customHeaders,
      createListingStep2Payload,
      200,
      `/api/v2/listing/<listingId>?expand=videos%2CbookablePeriods%2Ccosts`
    );

    utilities.get(
      `/api/v2/qualification/requirements/by-listing/${listingUuid}`,
      constants.haApiHeader,
      404,
      `/api/v2/qualification/requirements/by-listing/<listingUuid>`
    );
  });

  group("Create Listing Step 3", function () {
    const createListingStep3Payload = JSON.stringify({
      facilities: Object.assign(constants.facilities, {
        allergy_friendly: "yes",
        balcony_terrace: "private",
        basement: "private",
        bathroom: "private",
        garden: "private",
        wheelchair_accessible: "yes",
      }),
    });
    utilities.put(
      `/api/v2/listing/${listingId}?expand=videos%2CbookablePeriods%2Ccosts`,
      customHeaders,
      createListingStep3Payload,
      200,
      `/api/v2/listing/<listingId>?expand=videos%2CbookablePeriods%2Ccosts`
    );
    utilities.get(
      `/api/v2/qualification/requirements/by-listing/${listingUuid}`,
      constants.haApiHeader,
      404,
      `/api/v2/qualification/requirements/by-listing/<listingUuid>`
    );
  });

  group("Create Listing Step 4", function () {
    const createListingStep4Payload = JSON.stringify({
      facilities: Object.assign(constants.facilities, {
        ac: "no",
        allergy_friendly: "yes",
        balcony_terrace: "private",
        basement: "private",
        bathroom: "private",
        bed: "yes",
        closet: "yes",
        desk: "yes",
        dishwasher: "yes",
        dryer: "yes",
        flooring: "laminate",
        garden: "private",
        heating: "central",
        kitchen: "private",
        kitchenware: "private",
        living_room: "private",
        lroom_furniture: "yes",
        parking: "private",
        toilet: "private",
        tv: "yes",
        washing_machine: "yes",
        wheelchair_accessible: "yes",
        wifi: "yes",
      }),
    });
    utilities.put(
      `/api/v2/listing/${listingId}?expand=videos%2CbookablePeriods%2Ccosts`,
      customHeaders,
      createListingStep4Payload,
      200,
      `/api/v2/listing/<listingId>?expand=videos%2CbookablePeriods%2Ccosts`
    );

    utilities.get(
      `/api/v2/qualification/requirements/by-listing/${listingUuid}`,
      constants.haApiHeader,
      404,
      `/api/v2/qualification/requirements/by-listing/<listingUuid>`
    );

    utilities.get(
      `/api/v2/user/me/kyc-info`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/kyc-info`
    );
  });

  group("Create Listing Step 5", function () {
    const createListingStep5Payload = JSON.stringify(constants.listingCosts);
    utilities.put(
      `/api/v2/listing/${listingId}?expand=videos%2CbookablePeriods%2Ccosts`,
      customHeaders,
      createListingStep5Payload,
      200,
      `/api/v2/listing/<listingId>?expand=videos%2CbookablePeriods%2Ccosts`
    );

    utilities.get(
      `/api/v2/qualification/requirements/by-listing/${listingUuid}`,
      constants.haApiHeader,
      404,
      `/api/v2/qualification/requirements/by-listing/<listingUuid>`
    );
    utilities.get(
      `/api/v2/user/me/kyc-info`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/kyc-info`
    );
  });

  group("Create Listing Step 6", function () {
    const createListingStep6Payload = JSON.stringify({
      preferredGender: "",
      minAge: 26,
      maxAge: 49,
      facilities: Object.assign(constants.facilities, {
        ac: "no",
        allergy_friendly: "yes",
        animal_allowed: "discussable",
        balcony_terrace: "private",
        basement: "private",
        bathroom: "private",
        bed: "yes",
        closet: "yes",
        desk: "yes",
        dishwasher: "yes",
        dryer: "yes",
        electricity_included: "yes",
        flooring: "laminate",
        garden: "private",
        gas_cost_included: "no",
        heating: "central",
        internet_included: "no",
        kitchen: "private",
        kitchenware: "private",
        living_room: "private",
        lroom_furniture: "yes",
        parking: "private",
        play_music: "no",
        registration_possible: "yes",
        smoking_allowed: "no",
        tenant_status: "working professionals only",
        toilet: "private",
        tv: "yes",
        washing_machine: "yes",
        water_cost_included: "yes",
        wheelchair_accessible: "yes",
        wifi: "yes",
      }),
      couplesAllowed: "yes",
      requiredDocumentsAgreement: true,
    });
    utilities.put(
      `/api/v2/listing/${listingId}?expand=videos%2CbookablePeriods%2Ccosts`,
      customHeaders,
      createListingStep6Payload,
      200,
      `/api/v2/listing/<listingId>?expand=videos%2CbookablePeriods%2Ccosts`
    );

    utilities.get(
      `/api/v2/qualification/requirements/by-listing/${listingUuid}`,
      constants.haApiHeader,
      404,
      `/api/v2/qualification/requirements/by-listing/<listingUuid>`
    );
  });

  group("Create Listing Step 7 - Upload Photos", function () {
    let res = utilities.get(
      `/api/v2/listing/${listingId}/photos`,
      constants.haApiHeader,
      200,
      `/api/v2/listing/<listingId>/photos`
    );

    let storageUrl;
    let photoId;
    let photoGuid;

    res = utilities.post(
      `/api/v2/listing/${listingId}/photos/request`,
      customHeaders,
      JSON.stringify({}),
      200,
      `/api/v2/listing/<listingId>/photos/request`
    );
    storageUrl = res.json().url;
    photoId = res.json().photoId;
    photoGuid = res.json().photoUuid;
    console.log("photo id is :- "+ photoId)
    console.log("storage url is :- " + storageUrl)
    
    res = http.put(storageUrl, binaryFile, {
      headers: {
        "user-agent":
          "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
        Referer: constants.baseUrl,
        "x-ha-lang": "en",
        "Content-Type": "image/png",
      },
      tags: {
        name: storageUrl,
      },
    });
    check(res, {
      "is status 200": (res) => res.status === 200,
    });

    utilities.put(
      `/api/v2/listing/${listingId}/photo/${photoId}/confirm`,
      customHeaders,
      JSON.stringify({}),
      200,
      `/api/v2/listing/<listingId>/photo/<photoId>/confirm`
    );
  });

  group("Publish Listing", function () {
    utilities.put(
      `/api/v2/listing/${listingId}/publish`,
      customHeaders,
      JSON.stringify({}),
      204,
      `/api/v2/listing/<listingId>/publish`
    );
    utilities.get(
      `/my/listings/${listingId}/success`,
      constants.haApiHeader,
      200,
      `/my/listings/<listingId>/success`
    );
    utilities.get(
      `/api/v2/user/me/information?expand=impersonation`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/information?expand=impersonation`
    );
    utilities.get(
      `/api/v2/user/me/realtime`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/realtime`
    );
    utilities.get(
      `/api/v2/conversations/counters/user/${advertiserId}`,
      constants.haApiHeader,
      200,
      `/api/v2/conversations/counters/user/<userId>`
    );
    utilities.get(
      `/api/v2/user/me/kyc-info`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/kyc-info`
    );

    utilities.get(
      `/api/v2/user/me/information?expand=savedPayoutMethodsCount%2CadvertiserBookingsCount%2CisKYCDone`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/information?expand=savedPayoutMethodsCount%2CadvertiserBookingsCount%2CisKYCDone`
    );

    let res = http.get(
      `https://moneta-public-api.stage.housinganywhere.com/v1/user/${advertiserId}/validate-billing-data`,
      {
        headers: constants.haApiHeader,
        tags: {
          name: `https://moneta-public-api.stage.housinganywhere.com/v1/user/<advertiserId>/validate-billing-data`,
        },
      }
    );
    check(res, {
      "is status 200": (res) => res.status === 200,
    });
  });
}
