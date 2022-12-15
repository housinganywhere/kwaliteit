# What all things can it setup for me?

Currently this setup can create following data :-

- Set of verified users
- Set of listings in Rotterdam
- Set of Listings with Booking Requests sent by tenants

# Pre-requisites

It's a golang console program(main package). So, you need to have golang installed on your machine or wherever you are planning to execute it. One can follow [this](https://go.dev/doc/install) link to have it installed and setup the machine

# What command to run and where will the data get created?

1. Navigate to the /Prerequisites/setup/src folder
2. Run following command
```
go run ./ -dataCategory user -count 3 -host https://stage.housinganywhere.com -exportLocation ../../../Data/Jsons/Users.json

```
Following is the significance of each of the parameters in the above command :-
* -dataCategory :- Indicates what type of data to be created. In the above mentioned example `-dataCategory user` indicates that we need to create users
The other dataCategories available are :-
    *  listing :- creates a listing
    *  BookingRequest :- creates a booking request
* -count :- Indicates quantity of data to be created. In the above mentioned example `-count 3` indicates 3 users need to be created
* -host :- Indicates against which environment data needs to be created. In the above example `-host https://stage.housinganywhere.com` indicates the data needs to be created against staging environment.
If not specified, by default it executes against staging environment
* -exportLocation :- Indicates at which location the json file of data needs to be stored. One can provide a relative path in here. In the above example `../../../Data/Jsons/Users.json` indicates the data needs to be stored inside /Data/Jsons/Users.json located at the root level of the PerformanceTesting folder

# Shape of the sample json data files created

**Users.json**
```
{
 "Users": [
  {
   "Username": "vcjbyrt@akzquba.info",
   "Password": "******",
   "Id": 711542,
   "Uuid": "78c33941-2a80-4987-b451-c43143cb37dd"
  },
  {
   "Username": "tnjyzpc@lusgbcd.com",
   "Password": "******",
   "Id": 711543,
   "Uuid": "15c73565-2e84-44e8-b529-0b4ec2966c51"
  },
  {
   "Username": "nxivrvy@enasrum.org",
   "Password": "******",
   "Id": 711544,
   "Uuid": "da27f931-09af-4f19-ac05-57d45d45f71f"
  }
 ]
}

```

listings.json
```
{
 "Listings": [
  {
   "Id": "2323483",
   "Uuid": "b0be1ee1-df9e-416f-bdbd-d09f27524de8",
   "AdvertiserId": "711549",
   "AdvertiserEmail": "uoluaiu@juynduo.org",
   "AdvertiserPassword": "******"
  },
  {
   "Id": "2323484",
   "Uuid": "0ede55e7-e72a-4583-b6f7-af7a6577ab9d",
   "AdvertiserId": "711550",
   "AdvertiserEmail": "uyigecg@ektsxnw.com",
   "AdvertiserPassword": "******"
  },
  {
   "Id": "2323485",
   "Uuid": "1fcfeebe-c2be-4caa-8d64-19664600b7e1",
   "AdvertiserId": "711551",
   "AdvertiserEmail": "qyqrqyj@wygqxrj.info",
   "AdvertiserPassword": "******"
  }
 ]
}

```

**BookingRequests.json**

```
{
 "BookingRequests": [
  {
   "LandlordUsername": "kbbptev@fnwchsg.com",
   "LandlordPassword": "******",
   "LandlordId": "711849",
   "LandlordUuid": "bab234e7-7dcf-4562-b41f-1e84c67b0948",
   "TenantId": "711850",
   "TenantUuid": "8171352f-8e5c-4f1f-a48d-aba1203c9c6c",
   "TenantUsername": "xqvuevt@uelhbof.biz",
   "TenantPassword": "****",
   "BookingId": "920686",
   "ListingId": "2377046",
   "ListingUuid": "bab234e7-7dcf-4562-b41f-1e84c67b0948",
   "StartDate": "2023-01-15",
   "EndDate": "2023-06-15"
  },
  {
   "LandlordUsername": "iffopht@wqixjmk.info",
   "LandlordPassword": "******",
   "LandlordId": "711847",
   "LandlordUuid": "b6662303-b4fb-49c3-b423-470c65d6f7cf",
   "TenantId": "711852",
   "TenantUuid": "974be1ac-1d2d-4933-9f26-bbd7b4bf99f8",
   "TenantUsername": "axisnab@tyxlsad.net",
   "TenantPassword": "******",
   "BookingId": "920688",
   "ListingId": "2377047",
   "ListingUuid": "b6662303-b4fb-49c3-b423-470c65d6f7cf",
   "StartDate": "2023-01-15",
   "EndDate": "2023-06-15"
  },
  {
   "LandlordUsername": "fmpqxio@xbmpwck.org",
   "LandlordPassword": "******",
   "LandlordId": "711848",
   "LandlordUuid": "5bfddd50-0614-43e9-9fd5-b0a2eb317cf7",
   "TenantId": "711851",
   "TenantUuid": "52149bff-1d34-4747-950a-5d7a6188628f",
   "TenantUsername": "uaedwyn@clotcfr.org",
   "TenantPassword": "******",
   "BookingId": "920687",
   "ListingId": "2377048",
   "ListingUuid": "5bfddd50-0614-43e9-9fd5-b0a2eb317cf7",
   "StartDate": "2023-01-15",
   "EndDate": "2023-06-15"
  }
 ]
}

```

