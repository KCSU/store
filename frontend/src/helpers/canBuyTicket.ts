import dayjs from "dayjs";
import { Formal } from "../model/Formal";

export function canBuyTicket(formal: Formal): boolean {
    // TODO: what if user has ticket already?
    let canBuy = dayjs(formal.saleEnd).isAfter(Date.now());
    return canBuy;
}

export function canEditTicket(formal: Formal): boolean {
    let canBuy = dayjs(formal.saleEnd).isAfter(Date.now());
    return canBuy;
}