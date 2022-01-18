import { useToast } from "@chakra-ui/react";
import axios from "axios";
import { useMutation, useQueryClient } from "react-query";
import { api } from "../config/api";
import { User } from "../model/User";

export function useLogin() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation(
    (cred: string) => {
      let params = new URLSearchParams();
      params.set("credential", cred);
      return api.post<User>("auth/callback", params.toString());
    },
    {
      onSuccess(response) {
        queryClient.setQueryData<User>("authUser", response.data);
      },
      onError(error) {
        if (axios.isAxiosError(error)) {
          toast({
            title: "Error",
            description: error.response?.data.message,
            status: "error",
          });
        }
      },
    }
  );
}
