import http from "k6/http";
import { check, group, sleep } from "k6";

export const options = {
  stages: [
    { duration: "30s", target: 20 }, // simulate ramp-up of traffic from 1 to 20 users over 30s.
    { duration: "1m", target: 100 }, // stay at 100 users for 1 minute
    { duration: "10s", target: 0 }, // ramp-down to 0 users
  ],
  thresholds: {
    http_req_duration: ["p(99)<1500"], // 99% of requests must completed below 1.5s
  },
};

const BASE_URL = "http://localhost:8080/distance";

export default () => {
  check(http.get(BASE_URL), {
    "status is 200": (r) => r.status == 200,
  });
};
