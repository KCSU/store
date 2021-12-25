import { useToast } from "@chakra-ui/react";
import axios from "axios";
import { useMutation, useQueryClient } from "react-query";
import { api } from "../config/api";
import { QueueRequest } from "../model/QueueRequest";

export function useBuyTicket() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation(
    async (qr: QueueRequest) => {
      // TODO: return type?
      return api.post<void>("tickets", qr);
    },
    {
      onSuccess() {
        // Alternatively, setQueryData?
        queryClient.invalidateQueries("tickets");
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
