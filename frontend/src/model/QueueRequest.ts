import { TicketRequest } from "./TicketRequest";

export interface QueueRequest {
    formalId: number;
    ticket: TicketRequest;
    guestTickets: TicketRequest[];
}