import { useQuery } from "react-query";
import { api } from "../../config/api";
import { Group } from "../../model/Group";

export function useGroup(id: string) {
  return useQuery<Group>(
    ["admin/groups", { id }],
    async () => {
      const response = await api.get<Group>(`admin/groups/${id}`);
      return response.data;
    },
    {
      staleTime: 5 * 60 * 1000, // 5 minutes
    }
  );
}
