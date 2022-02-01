import { useToast } from "@chakra-ui/react";
import axios from "axios";
import { useMutation, useQueryClient } from "react-query";
import { api } from "../config/api";
import { TicketRequest } from "../model/TicketRequest";

export function useAddTicket(formalId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation(
    async (tr: TicketRequest) => {
      // TODO: return type?
      return api.post<void>(`formals/${formalId}/tickets`, tr);
    },
    {
      onSuccess() {
        // Alternatively, setQueryData?
        queryClient.invalidateQueries("tickets");
        toast({
            title: "Success",
            description: "Your ticket has been added to the queue.",
            status: "success",
        });
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
