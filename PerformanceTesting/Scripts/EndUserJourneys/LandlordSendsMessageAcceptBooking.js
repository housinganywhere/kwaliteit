import { check } from "k6";
import http from "k6/http";
import { group } from "k6";
import { SharedArray } from "k6/data";
import ApplicationFlow from "../Common/ApplicationFlow.js";
import * as constants from "../Common/Constants.js";
import Utilities from "../Common/Utilities.js";

const bookingDetails = new SharedArray("bookingDetails", function () {
  return JSON.parse(
    open("../TestData/BookingRequests.json")
  ).BookingRequests;
});

export function landlordSendsMessageAcceptsBooking() {
  const utilities = new Utilities();
  const booking = bookingDetails[Math.floor(Math.random() * bookingDetails.length)];
  const appFlows = new ApplicationFlow();
  let customHeaders = Object.assign(
    {
      Origin: constants.baseUrl,
      "Content-Type": "application/json",
    },
    constants.haApiHeader
  );

  appFlows.visitHomePageAndLogin(booking.LandlordUsername, booking.LandlordPassword);

  group("Open the LDP from Inbox", function () {
    utilities.get(
      `/my/inbox`,
      constants.haApiHeader,
      200,
      `/my/inbox`
    );
    
    utilities.get(
      `/api/v2/user/me/information?expand=impersonation`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/information?expand=impersonation`
    );

  utilities.get(
      `/dist/manifest.webmanifest`,
      constants.haApiHeader,
      200,
      `/dist/manifest.webmanifest`
    );

  utilities.get(
      `/api/v2/user/me/realtime`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/realtime`
    );

  utilities.get(
      `/api/v2/conversations/counters/user/${booking.LandlordId}`,
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
      `/api/v2/conversations/user/me?expand=messages%2Clisting%2Cadvertiser%2Ctenant%2CrentalConditions&limit=10&order=messages-desc&state=inbox`,
      constants.haApiHeader,
      200,
      `/api/v2/conversations/user/me?expand=messages%2Clisting%2Cadvertiser%2Ctenant%2CrentalConditions&limit=10&order=messages-desc&state=inbox`
    );

  utilities.get(
      `/api/v2/user/me/kyc`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/kyc`
    );

  utilities.get(
      `/api/v2/user/me/realtime`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/realtime`
    );
  });

  group("Send a conversation message", function () {
    utilities.get(
      `/api/v2/user/me/information?expand=savedPayoutMethodsCount%2CadvertiserBookingsCount%2CisKYCDone`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/information?expand=savedPayoutMethodsCount%2CadvertiserBookingsCount%2CisKYCDone`
    );

    utilities.get(
      `/api/v2/user/${booking.LandlordId}/validate-billing-details`,
      constants.haApiHeader,
      200,
      `/api/v2/user/<userId>/validate-billing-details`
    );

    utilities.get(
      `/api/v2/listing/${booking.ListingId}?expand=photos%2Cexclusions%2CbookablePeriods`,
      constants.haApiHeader,
      200,
      `/api/v2/listing/<listingId>?expand=photos%2Cexclusions%2CbookablePeriods`
    );

    utilities.get(
      `/api/v2/conversation/${booking.BookingId}/payout-information`,
      constants.haApiHeader,
      200,
      `/api/v2/conversation/<bookingId>/payout-information`
    );

    utilities.get(
      `/api/v2/qualification/documents-shared?ownerId=${booking.LandlordId}`,
      constants.haApiHeader,
      200,
      `/api/v2/qualification/documents-shared?ownerId=<userId>`
    );

    utilities.get(
      `/api/v2/qualification/requirements/by-listing/${booking.ListingUuid}`,
      constants.haApiHeader,
      404,
      `/api/v2/qualification/requirements/by-listing/<listingUuid>`
    );
    
    const qualifyUserBody = JSON.stringify({
      listingUuid: booking.ListingUuid,
    });
    utilities.post(
      `/api/v2/qualification/qualify-user`,
      customHeaders,
      qualifyUserBody,
      404,
      `/api/v2/qualification/qualify-user`
    );

    utilities.put(
      `/api/v2/conversation/${booking.BookingId}/read`,
      customHeaders,
      JSON.stringify({}),
      204,
      `/api/v2/conversation/<bookingId>/read`
    );

    const rentalConditionsEstimateBody = JSON.stringify({
      source: { conversationId: parseInt(booking.BookingId) },
      startDate: booking.StartDate + "T00:00:00Z",
      endDate: booking.EndDate + "T00:00:00Z",
    });
    utilities.put(
      `/api/v2/bookings/rental-conditions/estimate`,
      customHeaders,
      rentalConditionsEstimateBody,
      204,
      `/api/v2/bookings/rental-conditions/estimate`
    );
    
    utilities.get(
      `/api/v2/conversations/counters/user/${booking.LandlordId}`,
      constants.haApiHeader,
      200,
      `/api/v2/conversations/counters/user/<userId>`
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
      `/api/v2/conversations/counters/user/${booking.LandlordId}`,
      constants.haApiHeader,
      200,
      `/api/v2/conversations/counters/user/<userId>`
    );

    const sendMessageBody = JSON.stringify({
      messageText:
        "Hey Tenant,\nThanks for your booking request\nI am accepting you \nHope you have a good stay",
      containsQuickReply: false,
      editedQuickReply: false,
    });
    utilities.put(
      `/api/v2/conversation/${booking.BookingId}/response?expand=messages`,
      customHeaders,
      sendMessageBody,
      200,
      `/api/v2/conversation/<bookingId>/response?expand=messages`
    );
  });

  group("Accept or Reject a tenant booking request", function () {
    utilities.put(
      `/api/v2/conversation/${booking.BookingId}/charge-tenant`,
      Object.assign(customHeaders, {'Referer': constants.baseUrl+"/my/talk/"+booking.BookingId}),
      JSON.stringify({}),
      200,
      `/api/v2/conversation/<bookingId>/charge-tenant`
    );
    
    let res = utilities.get(
      `/api/v2/conversation/${booking.BookingId}?expand=messages%2Cadvertiser%2Ctenant%2CstateHistory%2Ctenancy%2Clisting%2CrentalConditions`,
      constants.haApiHeader,
      200,
      `/api/v2/conversation/<bookingId>?expand=messages%2Cadvertiser%2Ctenant%2CstateHistory%2Ctenancy%2Clisting%2CrentalConditions`
    );
    const tenancyUuid = res.json().tenancy.uuid;
    
    utilities.get(
      `/api/v2/conversation/${booking.BookingId}?expand=messages%2Cadvertiser%2Ctenant%2CstateHistory%2Ctenancy%2CrentalConditions`,
      constants.haApiHeader,
      200,
      `/api/v2/conversation/<bookingId>?expand=messages%2Cadvertiser%2Ctenant%2CstateHistory%2Ctenancy%2CrentalConditions`
    );
    
    utilities.put(
      `/api/v2/conversation/${booking.BookingId}/read`,
      customHeaders,
      JSON.stringify({}),
      204,
      `/api/v2/conversation/<bookingId>/read`
    );

    utilities.get(
      `/api/v2/user/${booking.LandlordId}/conversation/${booking.BookingId}/commission-invoice`,
      constants.haApiHeader,
      200,
      `/api/v2/user/<userId>/conversation/<bookingId>/commission-invoice`
    );

    const qualifyUserBody = JSON.stringify({
      listingUuid: booking.ListingUuid,
    });
    utilities.post(
      `/api/v2/qualification/qualify-user`,
      customHeaders,
      qualifyUserBody,
      404,
      `/api/v2/qualification/qualify-user`
    );

    res = http.get(
      `https://osiris.s.${constants.host}/api/v1/payments/overview?tenancy=${tenancyUuid}`,
      {
        headers: constants.haApiHeader,
        tags: {
          name: `https://osiris.s.${constants.host}/api/v1/payments/overview?tenancy=<tenancyId>`,
        },
      }
    );
    check(res, {
      "is status 200": (res) => res.status === 200,
    });

    utilities.get(
      `/api/v2/conversations/counters/user/${booking.LandlordId}`,
      constants.haApiHeader,
      200,
      `/api/v2/conversations/counters/user/<userId>`
    );

    utilities.get(
      `/api/v2/listing/${booking.ListingId}/photos`,
      constants.haApiHeader,
      200,
      `/api/v2/listing/<listingId>/photos`
    );
  });
}
