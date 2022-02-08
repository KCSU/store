import { useToast } from "@chakra-ui/react";
import { useQueryClient } from "react-query";
import { api } from "../config/api";
import { FormalTicket } from "../model/Ticket";
import { useCustomMutation } from "./useCustomMutation";

export function useEditTicket(ticketId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useCustomMutation(
    (option: string) => {
      return api.put<void>(`tickets/${ticketId}`, { option });
    },
    {
      onSuccess(_, option) {
        const tickets = queryClient.getQueryData<FormalTicket[]>("tickets");
        if (tickets) {
          const newTix = tickets.map((t) => updateTicket(t, ticketId, option));
          queryClient.setQueryData("tickets", newTix);
        }
        toast({
          title: "Changes Saved",
          status: "success"
        })
      }
    }
  );
}

function updateTicket(
  t: FormalTicket,
  ticketId: number,
  option: string
): FormalTicket {
  if (t.ticket.id === ticketId) {
    return {
      ...t,
      ticket: {
        ...t.ticket,
        option,
      },
    };
  }
  let guestTickets = t.guestTickets.map((gt) => {
    if (gt.id === ticketId) {
      return { ...gt, option };
    } else {
      return gt;
    }
  });
  return {
    ...t,
    guestTickets,
  };
}
