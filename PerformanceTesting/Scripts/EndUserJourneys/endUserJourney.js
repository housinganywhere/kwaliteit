import { sleep } from "k6";
import { createListingJourney } from "./CreateListingProfile.js";
import { searchListingAndOpenLDPJourney } from "./SearchAndOpenLDP.js";
import { searchListingAndRequestBookingJourney } from "./SearchAndBook.js";
import { landlordSendsMessageAcceptsBooking } from "./LandlordSendsMessageAcceptBooking.js";

let myOptions = JSON.parse(open(__ENV.LOAD_TYPE));

  // export const options = {
  //   "scenarios": {
  //     "createListing": {
  //       "exec": "createListinge2e",
  //       "executor": "per-vu-iterations",
  //       "vus": 1,
  //       "iterations": 1,
  //       "maxDuration": "5m"
  //     },
  //     "searchListingsAndOpenLDP": {
  //       "exec": "searchAndOpenLDPe2e",
  //       "executor": "per-vu-iterations",
  //       "vus": 1,
  //       "iterations": 1,
  //       "maxDuration": "5m"
  //     },
  //     "searchListingsAndSendBookingRequest": {
  //       "exec": "searchListingAndSendBookingRequest",
  //       "executor": "per-vu-iterations",
  //       "vus": 1,
  //       "iterations": 1,
  //       "maxDuration": "5m"
  //     },
  //     "landlordSendsMessangeAndActsOnBookingRequest": {
  //       "exec": "landlordSendsMessangeAndActsOnBookingRequest",
  //       "executor": "per-vu-iterations",
  //       "vus": 1,
  //       "iterations": 1,
  //       "maxDuration": "5m"
  //     }
  //   }
  // };

export const options = {
  scenarios: myOptions.scenarios
};

export function createListinge2e() {
  createListingJourney();
  sleep(15);
}

export function searchAndOpenLDPe2e() {
  searchListingAndOpenLDPJourney();
}

export function searchListingAndSendBookingRequest() {
  searchListingAndRequestBookingJourney();
}

export function landlordSendsMessangeAndActsOnBookingRequest() {
  landlordSendsMessageAcceptsBooking();
}
