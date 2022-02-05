import axios from "axios";
import { useQuery } from "react-query";
import { api } from "../config/api";
import { User } from "../model/User";

export function useAuthUser() {
    return useQuery<User | undefined>('authUser', async () => {
        try {
            const r = await api.get<User>("oauth/user");
            return r.data;
        } catch (err) {
            if (axios.isAxiosError(err) && err.response?.status === 401) {
                return undefined;
            }
            throw err;
        }
    }, {
        staleTime: Infinity
        // this should be manually invalidated in certain cases
    })
}