import http from "k6/http";
import { check, sleep } from "k6";


// jalankan:
// k6 run --env BASE_URL="http://localhost:8080" tests/load-test.js


const BASE_URL = __ENV.BASE_URL || "http://localhost:8080";


export const options = {

    stages: [

        {
            duration: "30s",
            target: 10000
        },

        {
            duration: "1m",
            target: 10000
        },

        {
            duration: "30s",
            target: 0
        }

    ],

    thresholds: {

        http_req_duration:[
            "p(95)<500"
        ],

        http_req_failed:[
            "rate<0.01"
        ]

    }

};



export function setup() {


    let tokens = [];


    // buat token sebanyak user test
    for(let i = 0; i < 10000; i++){


        let login = http.get(
            `${BASE_URL}/login`
        );


        if(login.status !== 200){

            console.error(
                "LOGIN FAILED:",
                login.status
            );

            continue;
        }


        let data = login.json();


        tokens.push(
            data.token
        );


    }


    console.log(
        "TOTAL TOKENS:",
        tokens.length
    );


    return {
        tokens
    };

}



export default function(data){


    let index = (__VU - 1) % data.tokens.length;


    let token = data.tokens[index];


    let res = http.get(
        `${BASE_URL}/protected/me`,
        {

            headers:{

                "Authorization":
                    `Bearer ${token}`

            }

        }
    );



    check(res,{

        "authorized":
            (r)=>r.status===200

    });



    sleep(1);

}