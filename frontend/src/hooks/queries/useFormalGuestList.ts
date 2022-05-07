import { useQuery } from "react-query";
import { api } from "../../config/api";
import { FormalGuest } from "../../model/FormalGuest";

export function useFormalGuestList(id: string) {
  return useQuery<FormalGuest[]>(
    ["formals", { id }, "guests"],
    async () => {
      const response = await api.get<FormalGuest[]>(`formals/${id}/guests`);
      return response.data;
    },
    {
      staleTime: 5 * 60 * 1000, // 5 minutes
    }
  );
}
