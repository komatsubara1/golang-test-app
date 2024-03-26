import http from "k6/http";
import { check } from "k6";

export default function() {
    // ユーザー作成
    let body = JSON.stringify({
        user_name: "test user",
    });
    let headers = {headers: {
        "Content-Type": "application/json",
    }};
    let res = http.post("http://localhost:8080/create", body, headers);

    let j = JSON.parse(res.body);

    check(res, {
        "status is 200": (r) => r.status === 200,
    });

    body = JSON.stringify({
        user_id: j.user_id,
    })
    headers = {headers: {
        "Content-Type": "application/json",
    }};
    res = http.post("http://localhost:8080/login", body, headers);

    j = JSON.parse(res.body);

    let authorization = res.headers.Authorization;
}