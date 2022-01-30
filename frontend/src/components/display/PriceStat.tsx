import { Stat, StatLabel, StatNumber, StatHelpText } from "@chakra-ui/react";
import { formatMoney } from "../../helpers/formatMoney";
import { Formal } from "../../model/Formal";

export interface PriceStatProps {
  formal: Formal;
  guestTickets: any[];
}

export function PriceStat({ formal, guestTickets }: PriceStatProps) {
  return (
    <Stat textAlign="center">
      <StatLabel>Overall ticket price:</StatLabel>
      <StatNumber>
        {formatMoney(formal.price + formal.guestPrice * guestTickets.length)}
      </StatNumber>
      <StatHelpText>Will be added to college bill</StatHelpText>
    </Stat>
  );
}
