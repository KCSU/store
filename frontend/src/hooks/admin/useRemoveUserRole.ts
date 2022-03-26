import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export interface RemoveUserRoleDto {
  roleId: number;
  email: string;
}

export function useRemoveUserRole() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (dto: RemoveUserRoleDto) => {
      return api.delete<void>(`admin/roles/users`, {
        params: dto,
      });
    },
    {
      onSuccess() {
        queryClient.invalidateQueries("admin/userRoles");
        toast({
          title: "Successfully revoked role",
          status: 'success'
        });
      },
    }
  );
}
