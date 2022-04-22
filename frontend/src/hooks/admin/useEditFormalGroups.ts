import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useEditFormalGroups(formalId: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (groups: string[]) => {
      return api.put<void>(`admin/formals/${formalId}/groups`, groups);
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/formals");
        queryClient.invalidateQueries("formals");
        toast({
          title: "Changes Saved",
          status: "success",
        });
      },
    }
  );
}
