import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useEditRole(roleId: number) {
    const queryClient = useQueryClient();
    const toast = useToast();
    return useCustomMutation(
        (name: string) => {
            return api.put<void>(`admin/roles/${roleId}`, {name});
        },
        {
            onSuccess() {
                queryClient.invalidateQueries("admin/roles");
                queryClient.invalidateQueries("admin/userRoles");
                toast({
                    title: "Changes Saved",
                    status: "success"
                })
            }
        }
    )
}