import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { Permission } from "../../model/Permission";
import { useCustomMutation } from "../mutations/useCustomMutation";

export type CreatePermissionDto = Omit<Permission, "id"> & {
    roleId: string
}

export function useCreatePermission() {
  const queryClient = useQueryClient();
  return useCustomMutation(async (p: CreatePermissionDto) => {
    return api.post<void>("admin/permissions", p);
  },
  {
    onSuccess() {
      queryClient.invalidateQueries("admin/roles");
    }
  });
}
