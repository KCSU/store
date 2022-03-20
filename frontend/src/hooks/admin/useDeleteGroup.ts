import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useDeleteGroup(id: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    () => {
      return api.delete<void>(`admin/groups/${id}`);
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/groups");
        toast({
          title: "Group deleted successfully",
          status: "success",
        });
      },
    }
  );
}
