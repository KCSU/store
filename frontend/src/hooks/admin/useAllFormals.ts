import { useQuery } from "react-query";
import { api } from "../../config/api";
import { FormalDto, parseFormal } from "../../helpers/parseFormal";
import { Formal } from "../../model/Formal";

export function useAllFormals() {
  return useQuery<Formal[]>(
    "formals",
    async () => {
      const response = await api.get<FormalDto[]>("admin/formals");
      return response.data.map(parseFormal);
    },
    {
      staleTime: 5 * 60 * 1000, // 5 minutes
    }
  );
}
