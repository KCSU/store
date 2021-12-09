import { Formal } from "../model/Formal";

export function getBuyText(formal: Formal): string {
  if (formal.saleStart > new Date()) {
    return "Join Queue";
  } else if (
    formal.guestTicketsRemaining === 0 &&
    formal.ticketsRemaining === 0
  ) {
    return "Join Waiting List";
  }
  return "Buy Tickets";
}
