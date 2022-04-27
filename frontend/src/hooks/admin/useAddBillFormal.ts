import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useAddBillFormal(billId: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (formalId: string) => {
      return api.post<void>(`admin/bills/${billId}/formals`, {
        formalId,
      });
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/bills");
        toast({
          title: "Formal Added to Bill",
          status: "success"
        })
      }
    }
  );
}