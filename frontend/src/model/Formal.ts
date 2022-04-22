import { Group } from "./Group";
import { ManualTicket } from "./ManualTicket";
import { AdminTicket } from "./Ticket";

export interface Formal {
    id: string;
    name: string;
    menu: string; // TODO:
    price: number;
    guestPrice: number;
    // options: string[]; // TODO:
    guestLimit: number; // TODO:
    tickets: number;
    ticketsRemaining: number;
    guestTickets: number;
    guestTicketsRemaining: number;
    saleStart: Date; // TODO:
    saleEnd: Date;
    dateTime: Date;
    // guestList: 
    // hidden?
    groups?: Group[];
    ticketSales?: AdminTicket[];
    manualTickets?: ManualTicket[];
}