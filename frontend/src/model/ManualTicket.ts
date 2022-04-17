export interface ManualTicket {
    id: number;
    name: string;
    justification: string;
    option: string;
    formalId: number;
    type: "complimentary" | "ents" | "standard" | "guest";
    billedTo?: string;
}