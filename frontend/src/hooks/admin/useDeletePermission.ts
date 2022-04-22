import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useDeletePermission(id: string) {
  const queryClient = useQueryClient();
  return useCustomMutation(
    () => {
      return api.delete<void>(`admin/permissions/${id}`);
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/roles");
      },
    }
  );
}
