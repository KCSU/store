import { Box, Heading, SimpleGrid, Spinner } from "@chakra-ui/react";
import { QueueOverview } from "../components/display/QueueOverview";
import { TicketOverview } from "../components/display/TicketOverview";
import { useTickets } from "../hooks/useTickets";

export function Tickets() {
  const templateColumns = {
    sm: "repeat(auto-fill, minmax(400px, 1fr))",
    base: "1fr",
  };
  // TODO: loading indicators
  const { data: tickets, isLoading, isError } = useTickets();
  if (!tickets || isLoading || isError) {
    return <Box></Box>;
  }
  return (
    <>
      <Heading size="xl" mb={5}>
        My Tickets
      </Heading>
      {tickets.queue.length > 0 && (
        <>
          <Heading size="md" as="h3" mb={4}>
            Ticket Queue <Spinner size="sm" speed="1s" ml={3}/>
          </Heading>
          <SimpleGrid gap={2} templateColumns={templateColumns} mb={5} autoFlow="dense">
            {tickets.queue.map((t, i) => {
              if ("guestTickets" in t) {
                return <Box gridRow="span 2" key={t.ticket.id} overflowX="auto">
                  <TicketOverview ticket={t} queue />
                </Box>;
              }
              return <Box key={`${t.formal.id}.${t.ticket.id}`}>
                <QueueOverview
                  ticket={t}
                />
              </Box>;
            })}
          </SimpleGrid>
        </>
      )}
      {tickets.tickets.length > 0 && (
        <>
          <Heading size="md" as="h3" mb={4}>
            Upcoming Formals
          </Heading>
          <SimpleGrid gap={2} templateColumns={templateColumns}>
            {tickets.tickets.map((t, i) => {
              return <Box key={t.ticket.id}>
                <TicketOverview ticket={t} />
              </Box>;
            })}
          </SimpleGrid>
        </>
      )}
    </>
  );
}
