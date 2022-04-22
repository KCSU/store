import { TicketRequest } from "./TicketRequest";

export interface QueueRequest {
    formalId: string;
    ticket: TicketRequest;
    guestTickets: TicketRequest[];
}