import { Bill } from "../model/Bill";

export type BillDto = Omit<Bill, "start" | "end"> & {
  start: string;
  end: string;
};

export function parseBill(dto: BillDto): Bill {
  return {
    ...dto,
    start: new Date(dto.start),
    end: new Date(dto.end),
  };
}