import http from "k6/http";
export const baseUrl = "https://stage.housinganywhere.com";
export const host = "stage.housinganywhere.com";
export const stripeApiUrl = "https://api.stripe.com";
export const stripeKey = "pk_test_qH4b4EehCxTKC9SYtKHvPsPK";
export const testCardNumber = "4111111111111111";
export const testCardCvc = "111";
export const testCardExpiryMonth = "03";
export const testCardExpiryYear = "23";
export const expect401Code = http.expectedStatuses(401);
export const expect404Code = http.expectedStatuses(404);
export const haApiHeader = {
  "user-agent":
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
  Referer: baseUrl,
  "x-ha-lang": "en",
};
export const listingBasicInfo = {
  address: "Landscheiding 101, 3045 NK Rotterdam, Netherlands",
  currency: "EUR",
  isCharitable: false,
  kind: 1,
  maximumStay: 24,
  minimumStayMonths: 1,
  price: 120000,
  type: 3,
};
export const facilities = {
  ac: null,
  allergy_friendly: null,
  animal_allowed: null,
  balcony_terrace: null,
  basement: null,
  bathroom: null,
  bed: null,
  bedroom_count: "2",
  bedroom_furnished: "yes",
  bedroom_size: null,
  cleaning_common: null,
  cleaning_private: null,
  closet: null,
  desk: null,
  dishwasher: null,
  dryer: null,
  electricity_included: null,
  flooring: null,
  garden: null,
  gas_cost_included: null,
  heating: null,
  housemates_gender: null,
  internet: null,
  internet_included: null,
  kitchen: null,
  kitchenware: null,
  living_room: null,
  lock: "yes",
  lroom_furniture: null,
  parking: null,
  play_music: null,
  registration_possible: null,
  smoking_allowed: null,
  tenant_status: "any",
  toilet: null,
  total_size: "60",
  tv: null,
  washing_machine: null,
  water_cost_included: null,
  wheelchair_accessible: null,
  wifi: null,
};

export const listingCosts = {
  contractType: "daily",
  cancellationPolicy: "strict",
  pricingType:"flat",
  pricingValues: {flat: 110000},
  costs: {'electricity-bill': {id: "electricity-bill",
                               payableBy: "included-in-rent",
                               payableAt: "monthly",
                               refundable: false,
                               required: true,
                               isEstimated: false,
                               value: 0},
          'water-bill': {id: "water-bill",
                        payableBy: "included-in-rent",
                        payableAt: "monthly",
                        refundable: false,
                        required: true,
                        isEstimated: false,
                        value: 0},
          'gas-bill': {id: "gas-bill",
                       payableBy: "tenant",
                       payableAt: "monthly",
                       refundable: false,
                       required: true,
                       isEstimated: true,
                       value: 2000},
          'internet-bill': {id: "internet-bill",
                            payableBy: "tenant",
                            payableAt: "monthly",
                            refundable: false,
                            required: true,
                            isEstimated: true,
                            value: 4600},
          'administration-fee': { id: "administration-fee",
                                  payableBy: "tenant",
                                  payableAt: "move-in",
                                  refundable: false,
                                  required: true,
                                  isEstimated: false,
                                  value: 24700},
          'cleaning-service': {id: "cleaning-service",
                               payableBy: "tenant",
                               payableAt: "move-out",
                               refundable: false,
                               required: false,
                               isEstimated: false,
                               value: 9700},
          'security-deposit': {id: "security-deposit",
                               payableBy: "tenant",
                               payableAt: "move-in",
                               refundable: true,
                               required: true,
                               isEstimated: false,
                               value: 97700}}}