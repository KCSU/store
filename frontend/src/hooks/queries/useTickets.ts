import { useMemo } from "react";
import { useQuery } from "react-query";
import { api } from "../../config/api";
import { Tickets } from "../../model/Queue";
import { FormalTicket, Ticket } from "../../model/Ticket";

export function useProcessedTickets(tickets: FormalTicket[]): Tickets {
  return useMemo(() => {
    let result: Tickets = {
      queue: [],
      tickets: []
    }
    for (let ticket of tickets) {
      // Is the ticket overall in queue?
      if (ticket.ticket.isQueue) {
        result.queue.push(ticket);
      } else {
        // Find non-queue guest tickets
        let nonQueue: Ticket[] = [];
        for (let guest of ticket.guestTickets) {
          if (guest.isQueue) {
            result.queue.push({
              formal: ticket.formal,
              ticket: guest
            });
          } else {
            nonQueue.push(guest);
          }
        }
        result.tickets.push({
          formal: ticket.formal,
          ticket: ticket.ticket,
          guestTickets: nonQueue
        });
      }
    }
    return result;
  }, [tickets]);
}

export function useTickets() {
  return useQuery<FormalTicket[]>("tickets", async () => {
    const {data} = await api.get<FormalTicket[]>("tickets");
    return data;
  }, {
    // TODO: change this
    staleTime: 5 * 60 * 1000, // 1 minute,
    refetchInterval(data) {
      if (!data) {
        return false;
      }
      if (data.some(t => t.ticket.isQueue || t.guestTickets.some(gt => gt.isQueue))) {
        return 60 * 1000;
      }
      return false;
    }
  });
}