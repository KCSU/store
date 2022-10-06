import { Formal } from "../model/Formal";

export function getBuyText(formal: Formal): string {
  if (formal.firstSaleStart > new Date()) {
    return "Buy Tickets";
  } else if (
    formal.secondSaleStart > new Date() &&
    formal.ticketsRemaining <= formal.secondSaleTickets
  ) {
    return "Buy Tickets";
  } else if (
    formal.guestTicketsRemaining === 0 &&
    formal.ticketsRemaining === 0
  ) {
    return "Join Waiting List";
  }
  return "Buy Tickets";
}
