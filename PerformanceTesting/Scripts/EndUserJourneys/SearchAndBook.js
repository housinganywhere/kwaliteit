import { check } from "k6";
import http from "k6/http";
import { group } from "k6";
import { SharedArray } from "k6/data";
import { scenario } from "k6/execution";
import * as constants from "../Common/Constants.js";
import Utilities from "../Common/Utilities.js";
import ApplicationFlow from "../Common/ApplicationFlow.js";
import { randomString } from "https://jslib.k6.io/k6-utils/1.2.0/index.js";

const sharedListingIds = new SharedArray("publishedListings", function () {
  return JSON.parse(open("../../Data/Jsons/listings.json")).Listings;
});

export function searchListingAndRequestBookingJourney() {
  const utilities = new Utilities();
  let date = new Date();
  date.setMonth(date.getMonth() + 1);
  let startDate = utilities.formatDate(date);
  date.setMonth(date.getMonth() + 6);
  let endDate = utilities.formatDate(date);
  let tenantId;
  const appFlows = new ApplicationFlow();
  let customHeaders = Object.assign(
    {
      Origin: constants.baseUrl,
      "Content-Type": "application/json",
    },
    constants.haApiHeader
  );

  appFlows.visitHomePageAndSearchListing(startDate);

  const listingId = sharedListingIds[scenario.iterationInTest].Id;
  const landlordId = sharedListingIds[scenario.iterationInTest].AdvertiserId;
  const listingUuid = sharedListingIds[scenario.iterationInTest].Uuid;

  group("Open LDP and Click Contact Landlord", function () {
    utilities.get(
      `/room/${listingId}/nl/Rotterdam/zweedsestraat?startDate=${startDate}`,
      constants.haApiHeader,
      200,
      `/room/<listingId>/nl/Rotterdam/<location>?startDate=<startDate>`
    );

    utilities.get(
      `/api/v2/listing/${listingId}?expand=photos%2Cexclusions%2CbookablePeriods%2Cvideos%2Ccosts&lang=en`,
      constants.haApiHeader,
      200,
      `/api/v2/listing/<listingId>?expand=photos%2Cexclusions%2CbookablePeriods%2Cvideos%2Ccosts&lang=en`
    );

    let geonameSearchCity;

    geonameSearchCity = JSON.stringify({
      query: "Rotterdam, Netherlands",
      languages: ["en"],
    });
    utilities.post(
      `/api/v2/geonames/search-city`,
      customHeaders,
      geonameSearchCity,
      200,
      `/api/v2/geonames/search-city`
    );
    
    utilities.get(
      `/api/v2/user/${landlordId}`,
      constants.haApiHeader,
      200,
      `/api/v2/user/<userId>`
    );

    utilities.get(
      `/api/v2/user/${landlordId}/information?expand=listingsCount%2CsavedPaymentMethodsCount%2CadvertiserBookingsCount%2CisKYCDone`,
      constants.haApiHeader,
      200,
      `/api/v2/user/<userId>/information?expand=listingsCount%2CsavedPaymentMethodsCount%2CadvertiserBookingsCount%2CisKYCDone`
    );

    utilities.get(
      `/api/v2/conversations/user/me?listingId=${listingId}&limit=1`,
      constants.haApiHeader,
      401,
      `/api/v2/conversations/user/me?listingId=<listingId>&limit=1`
    );

    const rentalConditionsEstimateBody = JSON.stringify({
      source: { listingId: parseInt(listingId) },
      startDate: startDate,
      endDate: endDate,
    });
    utilities.post(
      `/api/v2/bookings/rental-conditions/estimate`,
      customHeaders,
      rentalConditionsEstimateBody,
      200,
      `/api/v2/bookings/rental-conditions/estimate`
    );
   
    utilities.get(
      `/api/v2/currency/EUR/rates`,
      constants.haApiHeader,
      200,
      `/api/v2/currency/EUR/rates`
    );
  });

  group("Purchase Subscription And Send Booking Request", function () {
    utilities.get(
      `/api/v2/billing/subscriptions/products/tenant`,
      constants.haApiHeader,
      200,
      `/api/v2/billing/subscriptions/products/tenant`
    );
    utilities.get(
      `/register?return_url=%2Fsubscription%2Ftenant%2Factivate%3Frecurrence_type%3Dmonthly%26return_url%3D%252Fmy%252Fbook%252Fcontact%252F823951%253FstartDate%253D${startDate}%2526endDate%253D${endDate}`,
      constants.haApiHeader,
      200,
      `/register?return_url=%2Fsubscription%2Ftenant%2Factivate%3Frecurrence_type%3Dmonthly%26return_url%3D%252Fmy%252Fbook%252Fcontact%252F823951%253FstartDate%253D<startDate>%2526endDate%253D<endDate>`
    );

    const regEmail = `loadTesting` + randomString(8) + `@ha.com`;
    const fName = `fname` + randomString(8);
    console.log("**********Email is :-*********" + regEmail);
    const registrationPayload = JSON.stringify({
      email: regEmail,
      firstName: fName,
      lastName: "lName",
      notificationLang: "en",
      password: "Housing@1234",
    });
    let res = utilities.post(
      `/api/v2/user`,
      customHeaders,
      registrationPayload,
      200,
      `/api/v2/user`
    );
    tenantId = res.json().id;
    
    utilities.get(
      `/subscription/tenant/activate?recurrence_type=monthly&return_url=%2Fmy%2Fbook%2Fcontact%2F${tenantId}%3FstartDate%3D${startDate}%26endDate%3D${endDate}`,
      constants.haApiHeader,
      200,
      `/subscription/tenant/activate?recurrence_type=monthly&return_url=%2Fmy%2Fbook%2Fcontact%2F<tenantId>%3FstartDate%3D<startDate>%26endDate%3D<endDate>`
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
      424,
      `/api/v2/user/me/realtime`
    );

    utilities.get(
      `/api/v2/conversations/counters/user/${tenantId}`,
      constants.haApiHeader,
      200,
      `/api/v2/conversations/counters/user/<userId>`
    );

    utilities.get(
      `/api/v2/listings/favorites?userId=${tenantId}`,
      constants.haApiHeader,
      200,
      `/api/v2/listings/favorites?userId=<userId>`
    );

    utilities.get(
      `/api/v2/billing/subscriptions/products/tenant`,
      constants.haApiHeader,
      200,
      `/api/v2/billing/subscriptions/products/tenant`
    );

    const editUserBody = JSON.stringify({
      isAliased: true,
    });
    utilities.put(
      `/api/v2/user/${tenantId}`,
      customHeaders,
      editUserBody,
      200,
      `/api/v2/user/<userId>`
    );

    const stripePayload = {
      type: "card",
      "card[number]": constants.testCardNumber,
      "card[cvc]": constants.testCardCvc,
      "card[exp_month]": parseInt(constants.testCardExpiryMonth),
      "card[exp_year]": parseInt(constants.testCardExpiryYear),
      key: constants.stripeKey,
    };
    res = http.post(
      constants.stripeApiUrl + `/v1/payment_methods`,
      stripePayload,
      {
        headers: {
          "user-agent":
            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
          Referer: "https://js.stripe.com/",
          Origin: "https://js.stripe.com",
          Host: "api.stripe.com",
          "content-type": "application/x-www-form-urlencoded",
        },
        tags: {
          name: constants.stripeApiUrl + `/v1/payment_methods`,
        },
      }
    );
    check(res, {
      "is status 200": (res) => res.status === 200,
    });
    let stripePaymentMethodId = res.json().id;

    const subscriptionPayload = JSON.stringify({
      stripePaymentMethodId: stripePaymentMethodId,
      currency: "eur",
      priceCodes: ["tenant-subscription-month"],
      billingCountry: "NL",
      customerName: fName + " lName",
    });
    utilities.post(
      `/api/v2/billing/subscription`,
      customHeaders,
      subscriptionPayload,
      200,
      `/api/v2/billing/subscription`
    );
    
    utilities.get(
      `/api/v2/user/${tenantId}/searches?limit=1`,
      constants.haApiHeader,
      200,
      `/api/v2/user/<userId>/searches?limit=1`
    );

    utilities.get(
      `/api/v2/user/me`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me`
    );

    utilities.get(
      `/api/v2/user/me/information?expand=advertiserBookingsCount%2CsavedPaymentMethodsCount%2ClistingsCount%2CsubscriptionFeatures`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/information?expand=advertiserBookingsCount%2CsavedPaymentMethodsCount%2ClistingsCount%2CsubscriptionFeatures`
    );

    utilities.get(
      `/api/v2/conversations/user/me?listingId=${listingId}&asTenant=true`,
      constants.haApiHeader,
      200,
      `/api/v2/conversations/user/me?listingId=<listingId>&asTenant=true`
    );

    const propsectiveLeadPayload = JSON.stringify({
      listingId: parseInt(listingId),
      startDate: startDate,
      endDate: endDate,
    });
    utilities.post(
      `/api/v2/conversations/prospective-lead`,
      customHeaders,
      propsectiveLeadPayload,
      204,
      `/api/v2/conversations/prospective-lead`
    );
    
    utilities.get(
      `/api/v2/listing/${listingId}?expand=lastMonthConversationCount%2Ccosts`,
      constants.haApiHeader,
      200,
      `/api/v2/listing/<listingId>?expand=lastMonthConversationCount%2Ccosts`
    );

    utilities.get(
      `/api/v2/listing/${listingId}/exclusions`,
      constants.haApiHeader,
      200,
      `/api/v2/listing/<listingId>/exclusions`
    );

    utilities.get(
      `/api/v2/listing/${listingId}/bookable-periods`,
      constants.haApiHeader,
      200,
      `/api/v2/listing/<listingId>/bookable-periods`
    );

    utilities.get(
      `/api/v2/listing/${listingId}/photos`,
      constants.haApiHeader,
      200,
      `/api/v2/listing/<listingId>/photos`
    );

    const rentalConditionsEstimateBody2 = JSON.stringify({
      source: { listingId: parseInt(listingId) },
      startDate: startDate,
      endDate: endDate,
      bookingLink: false,
    });
    res = utilities.post(
      `/api/v2/bookings/rental-conditions/estimate`,
      customHeaders,
      rentalConditionsEstimateBody2,
      200,
      `/api/v2/bookings/rental-conditions/estimate`
    );
    console.log("****The response body is :-*****")
    console.log(res.json());
    console.log(res);
    utilities.get(
      `/api/v2/conversations/last-message`,
      constants.haApiHeader,
      404,
      `/api/v2/conversations/last-message`
    );

    const qualifyUserBody = JSON.stringify({
      listingUuid: listingUuid,
    });
    utilities.post(
      `/api/v2/qualification/qualify-user`,
      customHeaders,
      qualifyUserBody,
      404,
      `/api/v2/qualification/qualify-user`
    );
    
    utilities.get(
      `/api/v2/user/${landlordId}`,
      customHeaders,
      200,
      `/api/v2/user/<userId>`
    );
    
    utilities.post(
      `/api/v2/conversations/prospective-lead`,
      customHeaders,
      propsectiveLeadPayload,
      204,
      `/api/v2/conversations/prospective-lead`
    );

    utilities.get(
      `/api/v2/user/me/information?expand=savedPayoutMethodsCount%2CadvertiserBookingsCount%2CisKYCDone`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/information?expand=savedPayoutMethodsCount%2CadvertiserBookingsCount%2CisKYCDone`
    );

    utilities.get(
      `/api/v2/user/${tenantId}/validate-billing-details`,
      constants.haApiHeader,
      200,
      `/api/v2/user/<userId>/validate-billing-details`
    );

    utilities.get(
      `/api/v2/user/${tenantId}/validate-billing-details`,
      constants.haApiHeader,
      200,
      `/api/v2/user/<userId>/validate-billing-details`
    );

    const userInfo = JSON.stringify({
      phone: "681154843",
      phoneCountryCode: "NL",
      birthDay: 22,
      birthMonth: 3,
      birthYear: 1987,
      gender: "male",
      occupation: "working professional",
    });
    utilities.put(
      `/api/v2/user/me`,
      customHeaders,
      userInfo,
      200,
      `/api/v2/user/me`
    );

    const conversationBody = JSON.stringify({
      adClickID: null,
      adProvider: null,
      messageText:
        "Hello! I'm interested in renting your accommodation. I believe I match your tenant preferences. \nPlease get back as soon as possible",
      startDate: startDate,
      endDate: endDate,
      bookingLink: false,
      listingId: parseInt(listingId),
    });
    res = utilities.post(
      `/api/v2/conversation`,
      customHeaders,
      conversationBody,
      200,
      `/api/v2/conversation`
    );
    const conversationId = res.json().id;
    const rentalConditionsId = res.json().rentalConditionsId;
    
    utilities.get(
      `/api/v2/conversation/${conversationId}?expand=messages%2Cadvertiser%2Ctenant%2CstateHistory%2Ctenancy%2Clisting%2CrentalConditions`,
      constants.haApiHeader,
      200,
      `/api/v2/conversation/<conversationId>?expand=messages%2Cadvertiser%2Ctenant%2CstateHistory%2Ctenancy%2Clisting%2CrentalConditions`
    );

    utilities.get(
      `/api/v2/user/me/information?expand=advertiserBookingsCount`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/information?expand=advertiserBookingsCount`
    );

    utilities.get(
      `/api/v2/listing/${listingId}?expand=photos%2Cexclusions%2CbookablePeriods`,
      constants.haApiHeader,
      200,
      `/api/v2/listing/<listingId>?expand=photos%2Cexclusions%2CbookablePeriods`
    );

    utilities.get(
      `/api/v2/qualification/documents-shared?ownerId=${tenantId}`,
      constants.haApiHeader,
      200,
      `/api/v2/qualification/documents-shared?ownerId=<userId>`
    );

    utilities.get(
      `/api/v2/user/${landlordId}/information?expand=isKYCDone`,
      constants.haApiHeader,
      200,
      `/api/v2/user/<userId>/information?expand=isKYCDone`
    );

    utilities.get(
      `/api/v2/qualification/requirements/by-listing/${listingUuid}`,
      constants.haApiHeader,
      404,
      `/api/v2/qualification/requirements/by-listing/<listingUuid>`
    );

    utilities.post(
      `/api/v2/qualification/qualify-user`,
      customHeaders,
      qualifyUserBody,
      404,
      `/api/v2/qualification/qualify-user`
    );

    const messageInfo = JSON.stringify({
      messageText:
        "Can you please provide more details around the apartment?\nHow far is metro?\nHow far is bus stand?\nHow many stairs does it have?",
      containsQuickReply: false,
      editedQuickReply: false,
    });
    utilities.put(
      `/api/v2/conversation/${conversationId}/response?expand=messages`,
      customHeaders,
      messageInfo,
      200,
      `/api/v2/conversation/<conversationId>/response?expand=messages`
    );
    
    utilities.get(
      `/my/talk/${conversationId}/book-request`,
      constants.haApiHeader,
      200,
      `/my/talk/<conversationId>/book-request`
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

    res = utilities.get(
      `/api/v2/conversation/${conversationId}?expand=advertiser%2Clisting%2CrentalConditions`,
      customHeaders,
      200,
      `/api/v2/conversation/<conversationId>?expand=advertiser%2Clisting%2CrentalConditions`
    );

    const price = res.json().listing.price;
    const landlordUuid = res.json().advertiser.uuid;
    const tenantUuid = res.json().rentalConditions.tenantUuid;
    const externalReferenceId = res.json().uuid;
    
    utilities.get(
      `/api/v2/listings/favorites?userId=${tenantId}`,
      customHeaders,
      200,
      `/api/v2/listings/favorites?userId=<userId>`
    );

    utilities.get(
      `/api/v2/conversations/counters/user/${tenantId}`,
      customHeaders,
      200,
      `/api/v2/conversations/counters/user/<userId>`
    );

    utilities.get(
      `/api/v2/user/me/information?expand=advertiserBookingsCount`,
      customHeaders,
      200,
      `/api/v2/user/me/information?expand=advertiserBookingsCount`
    );

    utilities.get(
      `/api/v2/listing/${listingId}/photos`,
      customHeaders,
      200,
      `/api/v2/listing/<listingId>/photos`
    );

    const settlementBody = JSON.stringify({
      externalReference: externalReferenceId,
      externalReferenceType: "booking",
      purchaseValue: {
        cents: price,
        currency: "EUR",
      },
      sellerAccountUuid: landlordUuid,
      buyerAccountUuid: tenantUuid,
    });
    res = utilities.post(
      `/api/v2/payment/settlement`,
      customHeaders,
      settlementBody,
      200,
      `/api/v2/payment/settlement`
    );

    const settlementReference = res.json().reference;
    
    utilities.get(
      `/api/v2/payment/settlement/${settlementReference}/pay-in`,
      constants.haApiHeader,
      200,
      `/api/v2/payment/settlement/<settlementReference>/pay-in`
    );

    utilities.get(
      `/api/v2/payment/method/pay-in?limit=100`,
      constants.haApiHeader,
      200,
      `/api/v2/payment/method/pay-in?limit=100`
    );

    utilities.get(
      `/api/v2/payment/method/pay-in-types`,
      constants.haApiHeader,
      200,
      `/api/v2/payment/method/pay-in-types`
    );

    const tenantInfo = JSON.stringify({
      firstName: fName,
      lastName: "lName",
      location: "Rotterdam, Netherlands",
    });
    utilities.put(
      `/api/v2/user/${tenantId}`,
      customHeaders,
      tenantInfo,
      200,
      `/api/v2/user/<userId>`
    );
    
    utilities.get(
      `/api/v2/conversation/${conversationId}`,
      constants.haApiHeader,
      200,
      `/api/v2/conversation/<conversationId>`
    );

    const payInInfo = JSON.stringify({
      methodType: "stripe-card",
      currency: "EUR",
    });
    res = utilities.post(
      `/api/v2/payment/method/pay-in`,
      customHeaders,
      payInInfo,
      200,
      `/api/v2/payment/method/pay-in`
    );
   
    const payInId = res.json().uuid;
    const customerUuid = res.json().customerUuid;
    const stripeToken = res.json().nextSteps[0].token1;
    const stripeIntent = stripeToken.split("_secret_")[0];

    const stripeBody = {
      "payment_method_data[type]": "card",
      "payment_method_data[card][number]": "4111111111111111",
      "payment_method_data[card][cvc]": "111",
      "payment_method_data[card][exp_month]": parseInt("03"),
      "payment_method_data[card][exp_year]": 23,
      expected_payment_method_type: "card",
      use_stripe_sdk: true,
      key: "pk_test_qH4b4EehCxTKC9SYtKHvPsPK",
      client_secret: stripeToken,
    };
    res = http.post(
      constants.stripeApiUrl + `/v1/setup_intents/${stripeIntent}/confirm`,
      stripeBody,
      {
        headers: {
          "user-agent":
            "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
          Referer: "https://js.stripe.com/",
          Origin: "https://js.stripe.com",
          Host: "api.stripe.com",
          "content-type": "application/x-www-form-urlencoded",
        },
        tags: {
          name: constants.baseUrl + `/v1/setup_intents/<stripeIntent>/confirm`,
        },
      }
    );
    check(res, {
      "is status 200": (res) => res.status === 200,
    });
    const payInToken = res.json().payment_method;

    const payInMethodInfo = JSON.stringify({
      token1: payInToken,
    });
    res = utilities.put(
      `/api/v2/payment/method/pay-in/${payInId}`,
      customHeaders,
      payInMethodInfo,
      200,
      `/api/v2/payment/method/pay-in/<payInId>`
    );
   
    const paymentMethodRef = res.json().uuid;

    const paymentMethodInfo = JSON.stringify({
      paymentMethodReference: paymentMethodRef,
      paymentType: "paymentengine",
    });
    res = utilities.put(
      `/api/v2/conversation/${conversationId}/booking-request`,
      customHeaders,
      paymentMethodInfo,
      200,
      `/api/v2/conversation/<conversationId>/booking-request`
    );

    const bookingRequestConfirmationId = res.json().uuid;
    
    utilities.get(
      `/my/talk/${conversationId}`,
      constants.haApiHeader,
      200,
      `/my/talk/<conversationId>`
    );

    utilities.get(
      `/api/v2/conversations/counters/user/${tenantId}`,
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
      `/api/v2/listings/favorites?userId=${tenantId}`,
      constants.haApiHeader,
      200,
      `/api/v2/listings/favorites?userId=<userId>`
    );

    utilities.get(
      `/api/v2/conversation/${conversationId}?expand=messages%2Cadvertiser%2Ctenant%2CstateHistory%2Ctenancy%2Clisting%2CrentalConditions`,
      constants.haApiHeader,
      200,
      `/api/v2/conversation/<conversationId>?expand=messages%2Cadvertiser%2Ctenant%2CstateHistory%2Ctenancy%2Clisting%2CrentalConditions`
    );

    utilities.get(
      `/api/v2/conversations/counters/user/${tenantId}`,
      constants.haApiHeader,
      200,
      `/api/v2/conversations/counters/user/<userId>`
    );

    utilities.get(
      `/api/v2/user/me/information?expand=advertiserBookingsCount`,
      constants.haApiHeader,
      200,
      `/api/v2/user/me/information?expand=advertiserBookingsCount`
    );

    utilities.get(
      `/api/v2/user/${landlordId}/information?expand=isKYCDone`,
      constants.haApiHeader,
      200,
      `/api/v2/user/<userId>/information?expand=isKYCDone`
    );

    utilities.get(
      `/api/v2/qualification/documents-shared?ownerId=${tenantId}`,
      constants.haApiHeader,
      200,
      `/api/v2/qualification/documents-shared?ownerId=<userId>`
    );

    utilities.get(
      `/api/v2/listing/${listingId}?expand=photos%2Cexclusions%2CbookablePeriods`,
      constants.haApiHeader,
      200,
      `/api/v2/listing/<listingId>?expand=photos%2Cexclusions%2CbookablePeriods`
    );

    utilities.get(
      `/api/v2/qualification/requirements/by-listing/${listingUuid}`,
      constants.haApiHeader,
      404,
      `/api/v2/qualification/requirements/by-listing/<listingUuid>`
    );

    utilities.post(
      `/api/v2/qualification/qualify-user`,
      customHeaders,
      qualifyUserBody,
      404,
      `/api/v2/qualification/qualify-user`
    );
  });
}
