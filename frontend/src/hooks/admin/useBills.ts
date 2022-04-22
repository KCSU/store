import { useQuery } from "react-query";
import { api } from "../../config/api";
import { BillDto, parseBill } from "../../helpers/parseBill";
import { Bill } from "../../model/Bill";

export function useBills() {
  return useQuery<Bill[]>(
    "admin/bills",
    async () => {
      const response = await api.get<BillDto[]>("admin/bills");
      return response.data.map(parseBill);
    },
    {
      staleTime: 5 * 60 * 1000,
    }
  );
}
