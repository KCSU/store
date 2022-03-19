import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { Group } from "../../model/Group";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useCreateGroup() {
  const queryClient = useQueryClient();
  return useCustomMutation(async (g: Group) => {
    const dto = {
        name: g.name,
        type: g.type,
        lookup: g.lookup
    };
    return api.post<void>("admin/groups", dto);
  },
  {
    onSuccess() {
      queryClient.invalidateQueries("admin/groups");
      queryClient.invalidateQueries("groups");
    }
  });
}
