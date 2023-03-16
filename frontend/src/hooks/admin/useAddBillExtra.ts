import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

interface AddBillExtraDto {
  description: string;
  amount: number;
}

export function useAddBillExtra(billId: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (extra: AddBillExtraDto) => {
      return api.post<void>(`admin/bills/${billId}/extras`, extra);
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/bills");
        toast({
          title: "Extra Charge Added to Bill",
          status: "success"
        })
      }
    }
  );
}