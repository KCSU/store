import { useQuery } from "react-query";
import { api } from "../../config/api";
import { ScannedTicket } from "../../model/ScannedTicket";

export function useScanTicket(id: string) {
  return useQuery<ScannedTicket>(
    ["scan", {id}],
    async () => {
      const response = await api.get<ScannedTicket>(`scan/${id}`);
      return response.data;
    },
    {
      retry: false,
      staleTime: 10 * 1000 // 10 seconds
    }
  )
}