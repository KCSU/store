import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useEditTicket(ticketId: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (option: string) => {
      return api.put<void>(`admin/tickets/${ticketId}`, { option });
    },
    {
      onSuccess() {
        // FIXME: Specifically this formal? This is super broken
        // and wasteful.
        queryClient.invalidateQueries("admin/formals");
        toast({
          title: "Changes Saved",
          status: "success",
        });
      },
    }
  );
}
