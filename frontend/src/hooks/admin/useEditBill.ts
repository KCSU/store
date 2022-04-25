import { useToast } from "@chakra-ui/react";
import dayjs from "dayjs";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { Bill } from "../../model/Bill";
import { useCustomMutation } from "../mutations/useCustomMutation";

export type EditBillDto = Pick<Bill, "name" | "start" | "end">

export function useEditBill(billId: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (bill: EditBillDto) => {
      return api.put<void>(`admin/bills/${billId}`, {
        name: bill.name,
        start: dayjs(bill.start).format("YYYY-MM-DD"),
        end: dayjs(bill.end).format("YYYY-MM-DD"),
      });
    },
    {
      onSuccess(_, option) {
        queryClient.invalidateQueries(["admin/bills", {id: billId}]);
        toast({
          title: "Changes Saved",
          status: "success",
        });
      }
    }
  );
}
