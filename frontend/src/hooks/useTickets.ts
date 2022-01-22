import { useQuery } from "react-query";
import { api } from "../config/api";
import { Formal } from "../model/Formal";
import { FormalTicket, Ticket } from "../model/Ticket";


export interface QueueTicket {
  formal: Formal;
  ticket: Ticket;
}
export interface Tickets {
  queue: (QueueTicket | FormalTicket)[];
  tickets: FormalTicket[];
}

function processTickets(tickets: FormalTicket[]): Tickets {
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
}

export function useTickets() {
  return useQuery<Tickets>("tickets", async () => {
    const {data: tickets} = await api.get<FormalTicket[]>("tickets");
    return processTickets(tickets);
  }, {
    staleTime: 60 * 1000 // 1 minute
  });
}