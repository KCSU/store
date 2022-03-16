import { useQuery } from "react-query";
import { api } from "../../config/api";
import { FormalDto, parseFormal } from "../../helpers/parseFormal";
import { Formal } from "../../model/Formal";

export function useFormal(id: number) {
  return useQuery<Formal>(
    ["admin/formals", {id}],
    async () => {
      const response = await api.get<FormalDto>(`admin/formals/${id}`);
      return parseFormal(response.data);
    },
    {
      staleTime: 5 * 60 * 1000, // 5 minutes
    }
  );
}
