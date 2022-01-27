import { useQuery } from "react-query";
import { api } from "../config/api";
import { Formal } from "../model/Formal";

type FormalDto = Omit<Formal, 'saleStart' | 'saleEnd' | 'dateTime'> & {
  saleStart: string;
  saleEnd: string;
  dateTime: string;
}

function parseFormal(dto: FormalDto): Formal {
  return {
    ...dto,
    saleStart: new Date(dto.saleStart),
    saleEnd: new Date(dto.saleEnd),
    dateTime: new Date(dto.dateTime)
  }
}

export function useFormals() {
  return useQuery<Formal[]>('formals', async () => {
    const response = await api.get<FormalDto[]>("formals");
    return response.data.map(parseFormal);
  }, {
    staleTime: 60 * 1000 // 1 minute
  })
}