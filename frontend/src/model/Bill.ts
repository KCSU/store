import { Formal } from "./Formal";

export interface Bill {
    id: string;
    name: string;
    start: Date;
    end: Date;
    formals?: Formal[];
}