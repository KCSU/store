import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useAddBillFormals(billId: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (formalIds: string[]) => {
      return api.post<void>(`admin/bills/${billId}/formals`, {
        formalIds
      });
    },
    {
      onSuccess(_, fs) {
        const plural = fs.length === 1 ? "" : "s";
        queryClient.invalidateQueries("admin/bills");
        toast({
          title: `Formal${plural} Added to Bill`,
          status: "success"
        })
      }
    }
  );
}