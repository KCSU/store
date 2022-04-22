import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export interface AddUserRoleDto {
  roleId: string;
  email: string;
}

export function useAddUserRole() {
  const queryClient = useQueryClient();
  return useCustomMutation(
    (dto: AddUserRoleDto) => {
      return api.post<void>(`admin/roles/users`, dto);
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/userRoles");
      },
    }
  );
}
