import React from "react";
import { Formal } from "./Formal";

export interface Bill {
    id: string;
    name: string;
    start: Date;
    end: Date;
    formals?: Formal[];
}

export const BillContext = React.createContext<Bill>({} as Bill);