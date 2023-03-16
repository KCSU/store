import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useRemoveBillExtra(billId: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (extraId: string) => {
      return api.delete<void>(`admin/bills/${billId}/extras/${extraId}`);
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/bills");
        toast({
          title: "Extra Charge Removed from Bill",
          status: "success",
        });
      }
    }
  );
}