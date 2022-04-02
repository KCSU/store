import { Group } from "./Group";
import { AdminTicket } from "./Ticket";

export interface Formal {
    id: number;
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
}