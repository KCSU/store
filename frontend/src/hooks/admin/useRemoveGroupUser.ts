import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useRemoveGroupUser(id: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (email: string) => {
      return api.delete<void>(`admin/groups/${id}/users`, {
        params: {
          email,
        },
      });
    },
    {
      onSuccess(_, email) {
        queryClient.invalidateQueries("admin/groups");
        toast({
          title: `${email} removed from group`,
          status: "success",
        });
      },
    }
  );
}
