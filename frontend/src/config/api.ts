import axios from "axios"

const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL
});

if (import.meta.env.DEV) {
    api.defaults.withCredentials = true;
}

export {api}