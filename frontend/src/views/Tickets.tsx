import { Heading, SimpleGrid } from "@chakra-ui/react";
import { TicketOverview } from "../components/display/TicketOverview";
import { Card } from "../components/utility/Card";
import { useQueue } from "../hooks/useQueue";
import { useTickets } from "../hooks/useTickets";

export function Tickets() {
  const queue = useQueue();
  const tickets = useTickets();
  return (
    <>
      <Heading size="xl" mb={5}>
        My Tickets
      </Heading>
      {queue.length > 0 && (
        <>
          <Heading size="md" as="h3" mb={4}>
            Ticket Queue
          </Heading>
          <Card mb={5}>
            <SimpleGrid gap={2} minChildWidth="300px">
              {queue.map((t, i) => {
                return <TicketOverview ticket={t} key={i} queue/>;
              })}
            </SimpleGrid>
          </Card>
        </>
      )}
      <Heading size="md" as="h3" mb={4}>
        Upcoming Formals
      </Heading>
      <Card>
        <SimpleGrid gap={2} minChildWidth="300px">
          {tickets.map((t, i) => {
            return <TicketOverview ticket={t} key={i} />;
          })}
        </SimpleGrid>
      </Card>
    </>
  );
}
