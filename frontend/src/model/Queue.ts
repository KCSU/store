import { Bill } from "./Bill";
import { Formal } from "./Formal";
import { Ticket, FormalTicket } from "./Ticket";

export interface QueueTicket {
  formal: Formal;
  ticket: Ticket;
}
export interface Tickets {
  queue: (QueueTicket | FormalTicket)[];
  tickets: FormalTicket[];
  pastTickets: FormalTicket[];
  bills: {
    bill: Bill;
    tickets: FormalTicket[];
  }[];
}
