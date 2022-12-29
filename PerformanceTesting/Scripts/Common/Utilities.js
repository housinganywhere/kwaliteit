import { check } from 'k6';
import http from 'k6/http';
import * as constants from './Constants.js';

export default class Utilities {
    formatDate(date) {
        var d = new Date(date),
            month = '' + (d.getMonth() + 1),
            day = '' + d.getDate(),
            year = d.getFullYear();
    
        if (month.length < 2) 
            month = '0' + month;
        if (day.length < 2) 
            day = '0' + day;
    
        return [year, month, day].join('-');
    }

    get(urlPath, headers, expectedCode, tag){
        const statuscheckText = `status is ${expectedCode}`;
        const expectedStatusCallBack = http.expectedStatuses(expectedCode);
        let res = http.get(constants.baseUrl + urlPath, {
            headers: headers,
            tags: {
                name: constants.baseUrl + tag,
            },
            responseCallback: expectedStatusCallBack
        });
        check(res, {
            [statuscheckText]: (res) => res.status === expectedCode,
          });
         return res;
    }

    post(urlPath, headers, body, expectedCode, tag){
        const statuscheckText = `status is ${expectedCode}`;
        const expectedStatusCallBack = http.expectedStatuses(expectedCode);
        let res = http.post(constants.baseUrl + urlPath, body, {
            headers: headers,
            tags: {
                name: constants.baseUrl + tag,
            },
            responseCallback: expectedStatusCallBack
        });
        check(res, {
            [statuscheckText]: (res) => res.status === expectedCode,
          });
        return res;
    }

    put(urlPath, headers, body, expectedCode, tag){
        const statuscheckText = `status is ${expectedCode}`;
        const expectedStatusCallBack = http.expectedStatuses(expectedCode);
        let res = http.put(constants.baseUrl + urlPath, body, {
            headers: headers,
            tags: {
                name: constants.baseUrl + tag,
            },
            responseCallback: expectedStatusCallBack
        });
        check(res, {
            [statuscheckText]: (res) => res.status === expectedCode,
          });
          return res;
    }
}