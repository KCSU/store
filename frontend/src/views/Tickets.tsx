import {
  Heading,
  SimpleGrid,
} from "@chakra-ui/react";
import { TicketOverview } from "../components/display/TicketOverview";
import { Card } from "../components/utility/Card";
import { useTickets } from "../hooks/useTickets";

export function Tickets() {
  const tickets = useTickets();
  return (
    <Card>
      <Heading size="lg" mb={5}>
        My Tickets
      </Heading>
      <SimpleGrid>
        {tickets.map((t, i) => {
          return <TicketOverview ticket={t} key={i}/>;
        })}
      </SimpleGrid>
    </Card>
  );
}
