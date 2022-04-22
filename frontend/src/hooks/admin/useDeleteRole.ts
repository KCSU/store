import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useDeleteRole(id: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    () => {
      return api.delete<void>(`admin/roles/${id}`);
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/roles");
        queryClient.invalidateQueries("admin/userRoles");
        toast({
          title: "Role deleted successfully",
          status: "success",
        });
      },
    }
  );
}
