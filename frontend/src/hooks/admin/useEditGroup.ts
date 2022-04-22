import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { Group } from "../../model/Group";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useEditGroup(groupId: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (group: Group) => {
      return api.put<void>(`admin/groups/${groupId}`, {
        name: group.name,
        lookup: group.lookup,
        type: group.type,
      });
    },
    {
      onSuccess(_, option) {
        // TODO: update formal query with current data
        queryClient.invalidateQueries("admin/groups");
        queryClient.invalidateQueries("groups");
        toast({
          title: "Changes Saved",
          status: "success",
        });
        // const formals = queryClien
      },
    }
  );
}
