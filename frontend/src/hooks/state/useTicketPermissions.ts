import dayjs from "dayjs";
import { useContext, useMemo } from "react";
import { Formal } from "../../model/Formal";
import { UserContext } from "../../model/User";

export interface TicketPermissions {
    canBuy: boolean;
    isSaleEnded: boolean;
    isInGroup: boolean;
    hasTicket: boolean;
}

export function useTicketPermissions(formal: Formal): TicketPermissions {
    const user = useContext(UserContext);
    return useMemo(() => {
        if (!user) {
            return {
                canBuy: false,
                isSaleEnded: false,
                isInGroup: false,
                hasTicket: false
            };
        }
        const isSaleEnded = !dayjs(formal.saleEnd).isAfter(Date.now());
        // TODO: always get groups!!
        const isInGroup = formal.groups?.some(
            g => user.groups.some(h => h.id === g.id)
        ) ?? true;
        const hasTicket = (formal.myTickets?.length ?? 0) > 0;
        return {
            canBuy: isInGroup && !isSaleEnded && !hasTicket,
            isSaleEnded,
            isInGroup,
            hasTicket
        }
    }, [user, formal])
}

export function useCanEditTicket(formal: Formal): boolean {
    let canBuy = dayjs(formal.saleEnd).isAfter(Date.now());
    return canBuy;
}