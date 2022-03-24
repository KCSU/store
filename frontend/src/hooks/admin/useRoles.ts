import { useQuery } from "react-query";
import { api } from "../../config/api";
import { Role } from "../../model/Role";


export function useRoles() {
    return useQuery<Role[]>('admin/roles', async () => {
        const response = await api.get<Role[]>("admin/roles");
        return response.data;
    }, {
        staleTime: 5 * 60 * 1000,
    })
}