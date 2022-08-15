import { Formal } from "../model/Formal";

export type FormalDto = Omit<Formal, "firstSaleStart" | "secondSaleStart" | "saleEnd" | "dateTime"> & {
  firstSaleStart: string;
  secondSaleStart: string;
  saleEnd: string;
  dateTime: string;
};

export function parseFormal(dto: FormalDto): Formal {
  return {
    ...dto,
    firstSaleStart: new Date(dto.firstSaleStart),
    secondSaleStart: new Date(dto.secondSaleStart),
    saleEnd: new Date(dto.saleEnd),
    dateTime: new Date(dto.dateTime),
  };
}
