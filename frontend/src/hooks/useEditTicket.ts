import { useToast } from "@chakra-ui/react";
import axios from "axios";
import { useMutation, useQueryClient } from "react-query";
import { api } from "../config/api";
import { FormalTicket } from "../model/Ticket";

export function useEditTicket(ticketId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation(
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
