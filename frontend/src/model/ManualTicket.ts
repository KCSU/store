export type ManualTicketType = "complimentary" | "ents" | "standard" | "guest";

export interface ManualTicket {
    id: string;
    name: string;
    justification: string;
    option: string;
    formalId: string;
    type: ManualTicketType;
    email: string;
}