import { useMemo } from "react";
import { useQuery } from "react-query";
import { api } from "../../config/api";
import { parseBill } from "../../helpers/parseBill";
import { FormalDto, parseFormal } from "../../helpers/parseFormal";
import { Bill } from "../../model/Bill";
import { Tickets } from "../../model/Queue";
import { FormalTicket, Ticket } from "../../model/Ticket";

export function useProcessedTickets(tickets: FormalTicket[]): Tickets {
  return useMemo(() => {
    let result: Tickets = {
      queue: [],
      tickets: [],
      pastTickets: [],
      bills: [],
    }
    const bills = new Map<string, { bill: Bill; tickets: FormalTicket[] }>();
    for (let ticket of tickets) {
      const isFormalClosed = ticket.formal.saleEnd < new Date();
      const isFormalCompleted = ticket.formal.dateTime < new Date();
      // Is the ticket overall in queue?
      if (ticket.ticket.isQueue && !isFormalClosed) {
        result.queue.push(ticket);
      } else {
        // Find non-queue guest tickets
        let nonQueue: Ticket[] = [];
        for (let guest of ticket.guestTickets) {
          if (guest.isQueue && !isFormalClosed) {
            result.queue.push({
              formal: ticket.formal,
              ticket: guest
            });
          } else {
            nonQueue.push(guest);
          }
        }
        const ticketCollection = {
          formal: ticket.formal,
          ticket: ticket.ticket,
          guestTickets: nonQueue
        };
        // Check if formal has already been billed
        if (ticket.formal.bill) {
          if (bills.has(ticket.formal.bill.id)) {
            bills.get(ticket.formal.bill.id)?.tickets.push(ticketCollection);
          } else {
            bills.set(ticket.formal.bill.id, {
              bill: parseBill(ticket.formal.bill),
              tickets: [ticketCollection]
            });
          }
        } else if (isFormalCompleted) {
          result.pastTickets.push(ticketCollection);
        } else {
          result.tickets.push(ticketCollection);
        }
      }
    }
    result.bills = Array.from(bills.values()).sort(
      (a, b) => b.bill.start.getTime() - a.bill.start.getTime()
    );
    return result;
  }, [tickets]);
}

type FormalTicketDto = Omit<FormalTicket, "formal"> & {
  formal: FormalDto;
}

export function useTickets() {
  return useQuery<FormalTicket[]>("tickets", async () => {
    const {data} = await api.get<FormalTicketDto[]>("tickets");
    return data.map(t => ({
      ...t, formal: parseFormal(t.formal)
    }));
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