import React from "react";
import { BillDto } from "../helpers/parseBill";
import { Group } from "./Group";
import { ManualTicket } from "./ManualTicket";
import { AdminTicket, Ticket } from "./Ticket";

export interface Formal {
    id: string;
    name: string;
    menu: string;
    price: number;
    guestPrice: number;
    guestLimit: number;
    firstSaleTickets: number;
    ticketsRemaining: number;
    firstSaleGuestTickets: number;
    guestTicketsRemaining: number;
    firstSaleStart: Date;
    secondSaleStart: Date;
    secondSaleTickets: number;
    secondSaleGuestTickets: number;
    saleEnd: Date;
    dateTime: Date;
    hasGuestList: boolean;
    isVisible: boolean;
    // guestList: 
    // hidden?
    myTickets?: Ticket[];
    billId?: string;
    bill?: BillDto;
    groups?: Group[];
    ticketSales?: AdminTicket[];
    manualTickets?: ManualTicket[];
    queueLength?: number;
}

export const FormalContext = React.createContext<Formal>({} as Formal);