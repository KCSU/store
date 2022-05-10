import { useMutation, useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useLoginStatus } from "../state/useLoginStatus";

export function useLogout() {
  const queryClient = useQueryClient();
  const [, setLoginStatus] = useLoginStatus();
  return useMutation(
    () => {
      return api.post<void>("oauth/logout");
    },
    {
      onSuccess() {
        queryClient.setQueryData("authUser", undefined);
        setLoginStatus(new Date(0));
      },
      onError() {
        queryClient.setQueryData("authUser", undefined);
        setLoginStatus(new Date(0));
      },
    }
  );
}
