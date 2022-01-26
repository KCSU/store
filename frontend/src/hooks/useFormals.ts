import { useQuery } from "react-query";
import { api } from "../config/api";
import { Formal } from "../model/Formal";

export function useFormals() {
  return useQuery<Formal[]>('formals', async () => {
    const response = await api.get<Formal[]>("formals");
    return response.data;
  }, {
    staleTime: 60 * 1000 // 1 minute
  })
}