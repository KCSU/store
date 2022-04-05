import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { useCustomMutation } from "../mutations/useCustomMutation";

export function useCancelTicket(id: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    () => {
      return api.delete<void>(`admin/tickets/${id}`);
    },
    {
      onSuccess() {
        // TODO: Invalidate based on ticket formal id?
        queryClient.invalidateQueries("admin/formals");
        toast({
          title: "Ticket cancelled successfully",
          status: "success",
        });
      },
    }
  );
}
