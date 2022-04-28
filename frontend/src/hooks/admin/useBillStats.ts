import { useQuery } from "react-query";
import { api } from "../../config/api";
import { BillStats } from "../../model/BillStats";

export function useBillStats(id: string) {
  return useQuery<BillStats>(
    ["admin/bills", { id }, "stats"],
    async () => {
      const response = await api.get<BillStats>(`admin/bills/${id}/stats`);
      return response.data;
    },
    {
      staleTime: 5 * 60 * 1000, // 5 minutes
    }
  );
}
