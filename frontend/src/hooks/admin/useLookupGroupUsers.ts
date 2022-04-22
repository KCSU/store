import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useLookupGroupUsers(id: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    () => {
      return api.post<void>(`admin/groups/${id}/users/lookup`);
    },
    {
      onSuccess() {
        queryClient.invalidateQueries(["admin/groups", {id}]);
        toast({
          title: "Group users synced successfully",
          status: "success",
        });
      },
    }
  );
}
