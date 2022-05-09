import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Alert,
  AlertIcon,
  Box,
  Button,
  Heading,
  Icon,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  SimpleGrid,
  Spinner,
  Text,
  useDisclosure,
} from "@chakra-ui/react";
import dayjs from "dayjs";
import { useEffect } from "react";
import { FaReceipt } from "react-icons/fa";
import { BillTicketsTable } from "../components/tickets/BillTicketsTable";
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
  const { queue, tickets, pastTickets, bills } = useProcessedTickets(data);
  const { isOpen, onOpen, onClose } = useDisclosure();
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
      {pastTickets.length > 0 && (
        <>
          <Heading size="md" as="h3" mt={5} mb={1}>
            Current Bookings
          </Heading>
          <Text mb={3} fontSize="sm">
            Tickets for formal events that have already happened but which have
            not yet been billed.
          </Text>
          <SimpleGrid gap={2} templateColumns={templateColumns}>
            {pastTickets.map((t, i) => {
              return (
                <Box overflowX="auto" key={t.ticket.id}>
                  <TicketInfoCard ticket={t} />
                </Box>
              );
            })}
          </SimpleGrid>
        </>
      )}
      <Button
        colorScheme="brand"
        mt={4}
        onClick={onOpen}
        leftIcon={<Icon as={FaReceipt} />}
      >
        View Previous Bills
      </Button>
      <Modal isOpen={isOpen} onClose={onClose} size="4xl">
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Previous Bills</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            {bills.length === 0 && (
              <Alert status="info">
                <AlertIcon />
                No previous bills found.
              </Alert>
            )}
            <Accordion allowMultiple mb={5}>
              {bills.map(({ bill, tickets }) => (
                <AccordionItem key={bill.id}>
                  <AccordionButton>
                    <Box flex="1" textAlign="left">
                      {bill.name} ({dayjs(bill.start).format("ll")}
                      &ndash;
                      {dayjs(bill.end).format("ll")})
                    </Box>
                    <AccordionIcon />
                  </AccordionButton>
                  <AccordionPanel>
                    <BillTicketsTable tickets={tickets} />
                  </AccordionPanel>
                </AccordionItem>
              ))}
            </Accordion>
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
}
