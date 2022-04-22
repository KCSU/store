import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { Formal } from "../../model/Formal";
import { useCustomMutation } from "../mutations/useCustomMutation";

export type CreateFormalDto = Omit<Formal, "groups" | "id"> & { groups: string[] };

export function useCreateFormal() {
  const queryClient = useQueryClient();
  return useCustomMutation(async (f: CreateFormalDto) => {
    return api.post<void>("admin/formals", f);
  },
  {
    onSuccess() {
      queryClient.invalidateQueries("admin/formals");
      queryClient.invalidateQueries("formals");
    }
  });
}
