import { useQuery } from "react-query";
import { api } from "../../config/api";
import { Bill } from "../../model/Bill";

type BillDto = Omit<Bill, "start" | "end"> & {
  start: string;
  end: string;
};

export function useBills() {
  return useQuery<Bill[]>(
    "admin/bills",
    async () => {
      const response = await api.get<BillDto[]>("admin/bills");
      return response.data.map((dto) => ({
        ...dto,
        start: new Date(dto.start),
        end: new Date(dto.end),
      }));
    },
    {
      staleTime: 5 * 60 * 1000,
    }
  );
}
