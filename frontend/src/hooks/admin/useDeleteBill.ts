import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useDeleteBill(id: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    () => {
      return api.delete<void>(`admin/bills/${id}`);
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/bills");
        toast({
          title: "Bill deleted successfully",
          status: "success",
        });
      },
    }
  );
}