import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { Formal } from "../../model/Formal";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useEditFormal(formalId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (formal: Formal) => {
      return api.put<void>(`admin/formals/${formalId}`, formal);
    },
    {
      onSuccess(_, option) {
        // TODO: update formal query with current data
        queryClient.invalidateQueries("admin/formals");
        queryClient.invalidateQueries("formals")
        toast({
          title: "Changes Saved",
          status: "success"
        })
        // const formals = queryClien
      },
    }
  );
}
