import { check } from "k6";
import http from "k6/http";
import { group } from "k6";
import * as constants from "./Constants.js";
import Utilities from "./Utilities.js";

//custom headers are used
let customHeaders = Object.assign(
  {
    Origin: constants.baseUrl,
    "Content-Type": "application/json",
  },
  constants.haApiHeader
);
const utilities = new Utilities();

export default class ApplicationFlow {
  visitHomePageAndLogin(username, password) {
    group("visit Home Page And Login", function () {
      utilities.get("/", constants.haApiHeader, 200, "/");
      utilities.get(
        "/api/v2/user/me/information?expand=impersonation",
        constants.haApiHeader,
        401,
        "/api/v2/user/me/information?expand=impersonation"
      );
      const loginPayload = JSON.stringify({
        email: username,
        password: password,
      });
      utilities.post(
        "/api/v2/action/login",
        customHeaders,
        loginPayload,
        200,
        "/api/v2/action/login"
      );
    });
  }

  visitHomePageAndSearchListing(startDate) {
    group("visit Home Page And Search", function () {
      utilities.get("/", constants.haApiHeader, 200, "/");
      let geonameSearchCity = JSON.stringify({
        query: "Rotterdam, Netherlands",
        languages: ["en"],
      });
      utilities.post(
        "/api/v2/geonames/search-city",
        customHeaders,
        geonameSearchCity,
        200,
        "/api/v2/geonames/search-city"
      );
      utilities.get(
        "/api/v2/user/me/information?expand=impersonation",
        constants.haApiHeader,
        401,
        "/api/v2/user/me/information?expand=impersonation"
      );
      utilities.get(
        `/s/Rotterdam--Netherlands?startDate=${startDate}`,
        constants.haApiHeader,
        200,
        `/s/Rotterdam--Netherlands?startDate=<startDate>`
      );
      utilities.get(
        "/api/v2/user/me/information?expand=impersonation",
        constants.haApiHeader,
        401,
        "/api/v2/user/me/information?expand=impersonation"
      );
    });
  }
}
