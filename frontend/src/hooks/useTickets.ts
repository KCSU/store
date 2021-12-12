import { Ticket } from "../model/Ticket";
import { useFormals } from "./useFormals";

const ticketData = [
  {
    formalId: 1,
    option: "Normal",
    guestOptions: ["Pescetarian", "Normal"],
  },
];

export function useTickets(): Ticket[] {
  const formals = useFormals();
  return ticketData.map((t) => ({
    formal: formals.find((f) => f.id === t.formalId)!,
    ticket: {
      option: t.option,
    },
    guestTickets: t.guestOptions.map((g) => ({
      option: g,
    })),
  }));
}
