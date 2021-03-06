import { useQuery } from "react-query";
import { api } from "../../config/api";
import { Group } from "../../model/Group";

export function useGroups() {
    return useQuery<Group[]>('admin/groups', async () => {
        const response = await api.get<Group[]>("admin/groups");
        return response.data;
    }, {
        staleTime: 5 * 60 * 1000,
    })
}