import React from "react";
import { BillExtra } from "./BillExtra";
import { Formal } from "./Formal";

export interface Bill {
    id: string;
    name: string;
    start: Date;
    end: Date;
    formals?: Formal[];
    extras?: BillExtra[];
}

export const BillContext = React.createContext<Bill>({} as Bill);