import { Formal } from "./Formal";
import { TicketRequest } from "./TicketRequest";

// TODO: fix this
export interface Ticket {
    formal: Formal;
    ticket: TicketRequest;
    guestTickets: TicketRequest[];
}