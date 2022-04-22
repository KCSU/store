import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { TicketRequest } from "../../model/TicketRequest";
import { useCustomMutation } from "./useCustomMutation";

export function useAddTicket(formalId: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
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
      }
    }
  );
}
