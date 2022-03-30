import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useDeleteFormal(id: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    () => {
      return api.delete<void>(`admin/formals/${id}`);
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("formals");
        queryClient.invalidateQueries("admin/formals");
        toast({
          title: "Formal deleted successfully",
          status: "success",
        });
      },
    }
  );
}
