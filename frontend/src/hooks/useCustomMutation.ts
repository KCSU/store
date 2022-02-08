import { useToast } from "@chakra-ui/react";
import axios from "axios";
import { MutationFunction, useMutation, UseMutationOptions } from "react-query";

export function useCustomMutation<T, V = void>(
  mutationFn: MutationFunction<T, V>,
  options?: Omit<UseMutationOptions<T, unknown, V, unknown>, "mutationFn">
) {
  const toast = useToast();
  return useMutation(mutationFn, {
    ...options,
    onError(error, v, c) {
      if (axios.isAxiosError(error)) {
        toast({
          title: "Error",
          description: error.response?.data.message,
          status: "error",
        });
      }
      options?.onError?.(error, v, c)
    },
  });
}
