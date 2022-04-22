import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { ManualTicket } from "../../model/ManualTicket";
import { useCustomMutation } from "../mutations/useCustomMutation";

export type EditManualTicketDto = Omit<ManualTicket, "formalId" | "id">;

export function useEditManualTicket(ticketId: string) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (dto: EditManualTicketDto) => {
      return api.put<void>(`admin/tickets/manual/${ticketId}`, dto);
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
