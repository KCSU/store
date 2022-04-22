import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useAddGroupUser(id: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (email: string) => {
      return api.post<void>(`admin/groups/${id}/users`, {
        userEmail: email,
      });
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/groups");
        toast({
          title: "User Added to Group",
          status: "success"
        })
      },
    }
  );
}
