import { Formal } from "./Formal";
import { TicketRequest } from "./TicketRequest";

// TODO: fix this
export interface FormalTicket {
    formal: Formal;
    ticket: Ticket;
    guestTickets: Ticket[];
}

export interface Ticket {
    id: string;
    isGuest: boolean;
    isQueue: boolean;
    option: string;
    formalId: string;
    userId: string;
}

export type AdminTicket = Ticket & {
    userName: string;
    userEmail: string;
}