import { Ticket } from "../model/Ticket";
import { useTickets } from "./useTickets";

// TODO: make this work properly
// this needs to notify on queue success/failure,
// and all sorts of other fun stuff
export function useQueue(): Ticket[] {
    // we will remove this later
    const ticket = useTickets()[1];
    return [ticket];
}