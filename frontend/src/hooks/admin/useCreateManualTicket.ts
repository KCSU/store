import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../../config/api";
import { ManualTicket } from "../../model/ManualTicket";
import { useCustomMutation } from "../mutations/useCustomMutation";

export type CreateManualTicketDto = Omit<ManualTicket, "id">;

export function useCreateManualTicket() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    async (f: CreateManualTicketDto) => {
      return api.post<void>("admin/tickets/manual", f);
    },
    {
      onSuccess(_, f) {
        queryClient.invalidateQueries(["admin/formals", { id: f.formalId }]);
        toast({
          title: "Created ticket",
          status: "success",
        });
      },
    }
  );
}
