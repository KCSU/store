import { useQuery } from "react-query";
import { api } from "../../config/api";
import { UserRole } from "../../model/UserRole";


export function useUserRoles() {
    return useQuery<UserRole[]>('admin/userRoles', async () => {
        const response = await api.get<UserRole[]>("admin/roles/users");
        return response.data;
    }, {
        staleTime: 5 * 60 * 1000,
    })
}