import { useQuery } from "react-query";
import { api } from "../config/api";
import { Ticket } from "../model/Ticket";

export function useTickets() {
  return useQuery<Ticket[]>("tickets", async () => {
    const {data: tickets} = await api.get<Ticket[]>("tickets");
    return tickets;
  }, {
    staleTime: 60 * 1000 // 1 minute
  });
}