import axios from "axios";
import dayjs from "dayjs";
import { useQuery } from "react-query";
import { api } from "../../config/api";
import { User } from "../../model/User";
import { useLoginStatus } from "../state/useLoginStatus";

export function useAuthUser() {
    const [lastLogin, setLastLogin] = useLoginStatus();
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
        // this should be manually invalidated in certain cases
        staleTime: Infinity,
        onSuccess(u) {
            if (u === undefined) {
                return;
            }
            if (!lastLogin || lastLogin < dayjs().subtract(1, "day").toDate()) {
                setLastLogin(new Date());
            }
        }
    })
}