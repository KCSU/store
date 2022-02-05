import { useMutation, useQueryClient } from "react-query";
import { api } from "../config/api";

export function useLogout() {
  const queryClient = useQueryClient();
  return useMutation(
    () => {
      return api.post<void>("oauth/logout");
    },
    {
      onSuccess() {
        queryClient.setQueryData("authUser", undefined);
      },
      onError() {
        queryClient.setQueryData("authUser", undefined);
      },
    }
  );
}
