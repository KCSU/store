import { Box, Heading, SimpleGrid, Spinner } from "@chakra-ui/react";
import { QueueTicketInfoCard } from "../components/tickets/QueueTicketInfoCard";
import { TicketInfoCard } from "../components/tickets/TicketInfoCard";
import { useProcessedTickets, useTickets } from "../hooks/queries/useTickets";
import { FormalTicket } from "../model/Ticket";

export function TicketsView() {
  // TODO: loading indicators
  const { data, isLoading, isError } = useTickets();
  if (!data || isLoading || isError) {
    return <Box></Box>;
  }
  return <TicketsContent data={data} />;
}

interface TicketsContentProps {
  data: FormalTicket[];
}

function TicketsContent({ data }: TicketsContentProps) {
  const { queue, tickets } = useProcessedTickets(data);
  const templateColumns = {
    sm: "repeat(auto-fill, minmax(400px, 1fr))",
    base: "1fr",
  };
  return (
    <>
      <Heading size="xl" mb={5}>
        My Tickets
      </Heading>
      {queue.length > 0 && (
        <>
          <Heading size="md" as="h3" mb={4}>
            Ticket Queue <Spinner size="sm" speed="1s" ml={3} />
          </Heading>
          <SimpleGrid
            gap={2}
            templateColumns={templateColumns}
            mb={5}
            autoFlow="dense"
          >
            {queue.map((t, i) => {
              if ("guestTickets" in t) {
                return (
                  <Box gridRow="span 2" key={t.ticket.id} overflowX="auto">
                    <TicketInfoCard ticket={t} queue />
                  </Box>
                );
              }
              return (
                <Box key={`${t.formal.id}.${t.ticket.id}`}>
                  <QueueTicketInfoCard ticket={t} />
                </Box>
              );
            })}
          </SimpleGrid>
        </>
      )}
      {tickets.length > 0 && (
        <>
          <Heading size="md" as="h3" mb={4}>
            Upcoming Formals
          </Heading>
          <SimpleGrid gap={2} templateColumns={templateColumns}>
            {tickets.map((t, i) => {
              return (
                <Box overflowX="auto" key={t.ticket.id}>
                  <TicketInfoCard ticket={t} />
                </Box>
              );
            })}
          </SimpleGrid>
        </>
      )}
    </>
  );
}
