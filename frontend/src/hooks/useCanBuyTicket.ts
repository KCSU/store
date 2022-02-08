import dayjs from "dayjs";
import { useContext, useMemo } from "react";
import { Formal } from "../model/Formal";
import { UserContext } from "../model/User";

export function useCanBuyTicket(formal: Formal): boolean {
    const user = useContext(UserContext);
    if (!user) {
        return false;
    }
    return useMemo(() => {
        let canBuy = dayjs(formal.saleEnd).isAfter(Date.now());
        canBuy &&= formal.groups.some(
            g => user.groups.some(h => h.id === g.id)
        )
        return canBuy;
    }, [user, formal])
}

export function useCanEditTicket(formal: Formal): boolean {
    let canBuy = dayjs(formal.saleEnd).isAfter(Date.now());
    return canBuy;
}