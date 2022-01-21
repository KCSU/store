import { Box, Heading, SimpleGrid } from "@chakra-ui/react";
import { TicketOverview } from "../components/display/TicketOverview";
import { Card } from "../components/utility/Card";
import { useQueue } from "../hooks/useQueue";
import { useTickets } from "../hooks/useTickets";

export function Tickets() {
  const queue = useQueue();
  // TODO: loading indicators
  const {data: tickets, isLoading, isError} = useTickets();
  if (!tickets || isLoading || isError) {
    return <Box></Box>;
  }
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
            <SimpleGrid gap={2} templateColumns="repeat(auto-fill, minmax(350px, 1fr))">
              {queue.map((t, i) => {
                return <TicketOverview ticket={t} key={t.formal.id} queue/>;
              })}
            </SimpleGrid>
          </Card>
        </>
      )}
      <Heading size="md" as="h3" mb={4}>
        Upcoming Formals
      </Heading>
      {/* <Card> */}
        <SimpleGrid gap={2} templateColumns="repeat(auto-fill, minmax(350px, 1fr))">
          {tickets.map((t, i) => {
            return <TicketOverview ticket={t} key={t.formal.id} />;
  
          })}
        </SimpleGrid>
      {/* </Card> */}
    </>
  );
}
