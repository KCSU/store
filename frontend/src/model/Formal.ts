import React from "react";
import { Group } from "./Group";
import { ManualTicket } from "./ManualTicket";
import { AdminTicket, Ticket } from "./Ticket";

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
    hasGuestList: boolean;
    isVisible: boolean;
    // guestList: 
    // hidden?
    myTickets?: Ticket[];
    billId?: string;
    groups?: Group[];
    ticketSales?: AdminTicket[];
    manualTickets?: ManualTicket[];
}

export const FormalContext = React.createContext<Formal>({} as Formal);