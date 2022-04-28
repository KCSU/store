import dayjs from "dayjs";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { Bill } from "../../model/Bill";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useCreateBill() {
  const queryClient = useQueryClient();
  return useCustomMutation(
    async (b: Bill) => {
      const dto = {
        name: b.name,
        start: dayjs(b.start).format("YYYY-MM-DD"),
        end: dayjs(b.end).format("YYYY-MM-DD"),
      };
      return api.post<void>("admin/bills", dto);
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/bills");
      },
    }
  );
}