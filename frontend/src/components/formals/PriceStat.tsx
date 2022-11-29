import { Stat, StatLabel, StatNumber, StatHelpText } from "@chakra-ui/react";
import { useMemo } from "react";
import { formatMoney } from "../../helpers/formatMoney";
import { Formal } from "../../model/Formal";
import { Ticket } from "../../model/Ticket";
import { TicketRequest } from "../../model/TicketRequest";

export interface PriceStatProps {
  formal: Formal;
  isQueue?: boolean;
  guestTickets: Ticket[] | TicketRequest[];
}

export function PriceStat({ formal, isQueue = true, guestTickets }: PriceStatProps) {
  const total = useMemo(() => {
    let sum = formal.price;
    for (const t of guestTickets) {
      if (isQueue || ("isQueue" in t && !t.isQueue)) {
        sum += formal.guestPrice;
      }
    }
    return formatMoney(sum);
  }, [formal, guestTickets]);
  return (
    <Stat textAlign="center">
      <StatLabel>Overall ticket price:</StatLabel>
      <StatNumber>
        {total}
      </StatNumber>
      <StatHelpText>Will be added to college bill</StatHelpText>
    </Stat>
  );
}
