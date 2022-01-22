import { Formal } from "./Formal";
import { TicketRequest } from "./TicketRequest";

// TODO: fix this
export interface FormalTicket {
    formal: Formal;
    ticket: Ticket;
    guestTickets: Ticket[];
}

export interface Ticket {
    id: number;
    isGuest: boolean;
    isQueue: boolean;
    option: string;
    formalId: number;
    userId: number;
}