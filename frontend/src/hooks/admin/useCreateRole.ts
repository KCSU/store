import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useCreateRole() {
  const queryClient = useQueryClient();
  return useCustomMutation(async (name: string) => {
    return api.post<void>("admin/roles", {name});
  },
  {
    onSuccess() {
      queryClient.invalidateQueries("admin/roles");
    }
  });
}
