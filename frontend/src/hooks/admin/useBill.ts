import { useQuery } from "react-query";
import { api } from "../../config/api";
import { BillDto, parseBill } from "../../helpers/parseBill";
import { Bill } from "../../model/Bill";

export function useBill(id: string) {
  return useQuery<Bill>(
    ["admin/bills", { id }],
    async () => {
      const response = await api.get<BillDto>(`admin/bills/${id}`);
      return parseBill(response.data);
    },
    {
      staleTime: 5 * 60 * 1000, // 5 minutes
    }
  );
}