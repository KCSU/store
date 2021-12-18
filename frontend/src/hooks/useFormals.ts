import { useQuery } from "react-query";
import { api } from "../config/api";
import { Formal } from "../model/Formal";

export function useFormals() {
  return useQuery<Formal[]>('formals', async () => {
    const response = await api.get<Formal[]>("formals");
    const template: Partial<Formal> = {
      options: ["Normal", "Vegan", "Vegetarian", "Pescetarian"]
    };
    return response.data.map(f => Object.assign(f, template));
  }, {
    staleTime: 60 * 1000 // 1 minute
  })
}