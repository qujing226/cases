import axios, { AxiosInstance } from "axios";
import router from "next/router";

const instance: AxiosInstance = axios.create({
    baseURL: "http://localhost:8080",
    withCredentials: true
});

export interface Result<T> {
    code: number;
    msg: string;
    data: T;
}

instance.interceptors.response.use(function (resp) {
    const newToken = resp.headers["x-jwt-token"];
    const newRefreshToken = resp.headers["x-refresh-token"];
    if (newToken) {
        localStorage.setItem("token", newToken);
    }
    if (newRefreshToken) {
        localStorage.setItem("refresh_token", newRefreshToken);
    }
    if (resp.status == 401) {
        window.location.href = "/users/login";
    }
    return resp;
}, (err) => {
    console.log(err);
    if (err.response.status == 401) {
        window.location.href = "/users/login";
    }
    return err;
});

instance.interceptors.request.use((req) => {
    const token = localStorage.getItem("token");
    req.headers.setAuthorization("Bearer " + token, true);
    return req;
}, (err) => {
    console.log(err);
});

export default instance;
