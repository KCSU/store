import { Formal } from "../model/Formal";

export type FormalDto = Omit<Formal, "saleStart" | "saleEnd" | "dateTime"> & {
  saleStart: string;
  saleEnd: string;
  dateTime: string;
};

export function parseFormal(dto: FormalDto): Formal {
  return {
    ...dto,
    saleStart: new Date(dto.saleStart),
    saleEnd: new Date(dto.saleEnd),
    dateTime: new Date(dto.dateTime),
  };
}
