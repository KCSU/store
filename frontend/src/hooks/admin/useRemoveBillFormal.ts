import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useRemoveBillFormal(billId: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (formalId: string) => {
      return api.delete<void>(`admin/bills/${billId}/formals/${formalId}`);
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/bills");
        toast({
          title: "Formal Removed from Bill",
          status: "success",
        });
      }
    }
  );
}