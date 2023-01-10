# How to setup k6 performance test tool?

[k6](https://k6.io/) is the tool we are using to do performance/load testing at HA.
Installing k6 on MacOS is pretty straightforward. One can do it via homebrew :-
`brew install k6`
There are lot many other options as well one can use for different variants of Linux and Windows which can be found [here](https://k6.io/docs/get-started/installation/)

# What are the user journeys covered in k6 performance tests?

We currently cover 4 critical end user journeys, which would be subjected to a pre-configured load while executing the performance tests. Following are those user journeys:-

1. A probable tenant searching for a property and scanning listing details
    * Opens site
    * Enters Location, duration of stay and clicks on Search
    * Opens one of the listing and view its details
    * Clicks on Contact Landlord
    * Taken to the subscription page
Location:- /Scripts/EndUserJourneys/SearchAndOpenLDP.js

2. A probable tenant liking a listing and sending booking request
    * Opens site
    * Enters Location, duration of stay and clicks on Search
    * Opens one of the listing and view its details
    * Clicks on Contact Landlord
    * Taken to the subscription page
    * Registers
    * Purchases a subscription
    * Sends couple of messages to landlords
    * Sending a booking Request
Location:- /Scripts/EndUserJourneys/SearchAndBook.js

3. A landlord accepting or rejecting booking request from tenant
    * Landlord logins
    * Opens up booking request from the tenant
    * Views tenant's profile
    * Sends couple of messages to the tenant as a part of conversations
    * Accepts or Rejects the booking request
Location:- /Scripts/EndUserJourneys/LandlordSendsMessageActsOnBooking.js

4. A landlord creating a new listing
    * Landlord logins
    * Completes all the steps of listing creation flow and creates a listing
    * Publishes the listing
Location:- /Scripts/EndUserJourneys/CreateListingProfile.js


# Do these tests open a browser instance and execute like e2e tests?

No, these tests do not launch browser and execute. Each of the end user journey contains all the platform API requests + some web requests required to perform the journey successfully. Currently, it does not contain the .js static resources requests. 
So, for now we can monitor the performance of the backend (platform or other services involved) servers more accurately compared to tokamak.
Plan is to soon include those static resources calls as well so that we can load the tokamak servers as well during these tests


# How can I execute the tests?

Following steps need to be followed for executing the k6 performance tests :-
1. Open the root folder of this project, i.e. PerformanceTesting
2. First generate the test data required to run the performance tests. The test data generator is written in Golang and can be found in /Prerequisites/setup folder. The details around how to setup golang can be found here 
    * cd /Prerequisites/setup/src in the terminal and run below three commands

    * `go run ./ -dataCategory listing -count 1 -host https://stage.housinganywhere.com -exportLocation ../../../Scripts/TestData`
    * `go run ./ -dataCategory user -count 1 -host https://stage.housinganywhere.com -exportLocation ../../../Scripts/TestData`
    * `go run ./ -dataCategory bookingRequest -count 1 -host https://stage.housinganywhere.com -exportLocation ../../../Scripts/TestData`<br><br>
    * It will setup required test data of the domains:- user, listing and bookingRequest against the stage environment. The important fields of the test data of each domain gets exported into Scripts/TestData folder and can be found inside the respective .json file.
    *Note:- Make sure to change the value of parameters `count` and `host` depending upon the number of records of each category of test data you want to generate and the environment against you want to run the tests respectively. For simply trying out k6 scripts the above mentioned values should be fine.
3. Navigate to the root directory, i.e. PerformanceTesting
4. Fire the command :-
```
k6 run --out json=results.json --env LOAD_TYPE=Scripts/TestData/minLoad.config.json Scripts/EndUserJourneys/endUserJourney.js

```
Lets see what each of the parameter here indicates:-
* --out :- exports the metrics data of each and every point during the execution to results.json file.
* --env :- Set any environment variables to be accessed in the k6 script. e.g. here we are passing an enviornment variable LOAD_TYPE, the value of which points to the  minLoad.config.json file containing the load pattern to be used during the performance test.<br>
As the name suggests `minLoad.config.json` contains minimum load for each of the scenarios, i.e. 1 Virtual User performing each scenario. There are other config.json files present at the same location with higher load patterns, i.e.
     * `baseLoad.config.json` :- contains basic level of load pattern, i.e. pattern on the basis of last year average/peak
     * `highLoad.config.json` :- contains estimated load pattern for next peak season
        Note:- Right now these files wont have accurate numbers. Once we figure out the accurate or close to accurate load numbers, we can put in there. 
3. At the end of execution you will be able to see a summary metrics of the performance tests in the terminal like below:-
    ```
    Create Listing Step 3

       ✓ status is 200
       ✓ status is 404

     █ Create Listing Step 4

       ✓ status is 200
       ✓ status is 404

     █ Create Listing Step 5

       ✓ status is 200
       ✓ status is 404

     █ Accept or Reject a tenant booking request

       ✗ status is 200
        ↳  85% — ✓ 6 / ✗ 1
       ✓ status is 404
       ✓ is status 200

     █ Create Listing Step 6

       ✓ status is 200
       ✓ status is 404

     █ Create Listing Step 7 - Upload Photos

       ✓ status is 200
       ✓ is status 200

     █ Publish Listing

       ✓ status is 204
       ✓ status is 200
       ✓ is status 200

     checks.........................: 98.27% ✓ 171      ✗ 3  
     data_received..................: 3.7 MB 78 kB/s
     data_sent......................: 273 kB 5.7 kB/s
     group_duration.................: avg=6.58s    min=643.95ms med=4.46s    max=33.43s   p(90)=11.82s  p(95)=14.3s   
     http_req_blocked...............: avg=3.72ms   min=0s       med=1µs      max=151.19ms p(90)=2µs     p(95)=3.34µs  
     http_req_connecting............: avg=892.06µs min=0s       med=0s       max=52.06ms  p(90)=0s      p(95)=0s      
     http_req_duration..............: avg=714.92ms min=19.54ms  med=167.86ms max=12.76s   p(90)=1.67s   p(95)=3.84s   
       { expected_response:true }...: avg=724.94ms min=19.54ms  med=175.87ms max=12.76s   p(90)=1.67s   p(95)=3.84s   
     http_req_failed................: 1.72%  ✓ 3        ✗ 171
     http_req_receiving.............: avg=77.28ms  min=40µs     med=181.5µs  max=2.9s     p(90)=12.19ms p(95)=445.76ms
     http_req_sending...............: avg=1.01ms   min=29µs     med=106µs    max=155.35ms p(90)=202.7µs p(95)=283.84µs
     http_req_tls_handshaking.......: avg=2.12ms   min=0s       med=0s       max=116.54ms p(90)=0s      p(95)=0s      
     http_req_waiting...............: avg=636.63ms min=19.23ms  med=146.3ms  max=12.76s   p(90)=1.6s    p(95)=2.83s   
     http_reqs......................: 174    3.638125/s
     iteration_duration.............: avg=35.03s   min=14.7s    med=38.8s    max=47.82s   p(90)=47.42s  p(95)=47.62s  
     iterations.....................: 4      0.083635/s
     vus............................: 1      min=1      max=4
     vus_max........................: 4      min=4      max=4

    ```
    Some of the above metrics in the summary are self-explanatory. However, a short explanation around each metric can be found [here](https://k6.io/docs/using-k6/metrics/)
    
    `Scripts/EndUserJourneys/endUserJourney.js` is the entry point file for executing the performance tests. This file basically contains the method(s) for each of the end user journeys

    # How can I emit the metrics to prometheus?

    The above summary of metrics is good to get a high level idea about the performance of the system. But, in a real world scenario one would like to get more insights around how the system behaved(performed) at different points of time during the performance testing span. Such insights can only be obtained if we are able to store the metrics into a time-series DB like prometheus. On top of prometheus, we can set up graphs using grafana or google cloud monitoring for visualization. <br>
    k6 supports emiting these performance metrics to all the famous TimeStamp Dbs including prometheus.<br>
    For that we need to use the [official k6 prometheus extension](https://github.com/grafana/xk6-output-prometheus-remote). To build k6 binary with `Prometheus remote write output extension` use:
    ```
    xk6 build --with github.com/grafana/xk6-output-prometheus-remote@latest
    ```
    Use the new k6 binary to run k6 scripts, e.g.
    ```
    K6_PROMETHEUS_REMOTE_URL=http://localhost:9090/api/v1/write ./k6 run --env TEST_TYPE=../../Data/Jsons/minLoad.config.json Scripts/EndUserJourneys/endUserJourney.js -o output-prometheus-remote
    ```

     # How are the performance tests executed via CI / CD?
     -- COMING SOON --

