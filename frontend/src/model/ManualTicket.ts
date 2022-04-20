export type ManualTicketType = "complimentary" | "ents" | "standard" | "guest";

export interface ManualTicket {
    id: number;
    name: string;
    justification: string;
    option: string;
    formalId: number;
    type: ManualTicketType;
    email: string;
}