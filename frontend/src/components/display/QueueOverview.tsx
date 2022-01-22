import { HStack, Heading, Badge, Text, Button } from "@chakra-ui/react";
import { FaEdit, FaTrash, FaTrashAlt } from "react-icons/fa";
import { formatMoney } from "../../helpers/formatMoney";
import { useDateTime } from "../../hooks/useDateTime";
import { QueueTicket } from "../../hooks/useTickets";
import { Card } from "../utility/Card";

interface QueueOverviewProps {
  ticket: QueueTicket;
}

export function QueueOverview({ ticket }: QueueOverviewProps) {
  const datetime = useDateTime(ticket.formal.dateTime);
  return (
    <Card p={3} borderRadius="md">
      <HStack>
        <Heading size="md" as="h4">
          {ticket.formal.name}
        </Heading>
        <Badge colorScheme="teal">Guest</Badge>
        <Badge colorScheme="brand">In Queue</Badge>
      </HStack>
      <Text as="b">{datetime}</Text>
      <Text fontSize="sm">Meal Option: {ticket.ticket.option}</Text>
      <HStack justifyContent="flex-end">
        <Text as="b">{formatMoney(ticket.formal.price)}</Text>
        <Button size="xs" variant="outline" leftIcon={<FaEdit />}>
          Edit
        </Button>
        <Button size="xs" colorScheme="red" leftIcon={<FaTrashAlt />}>
          Cancel
        </Button>
        {/* <CancelTicketButton formalId={ticket.formal.id} isQueue={queue} /> */}
      </HStack>
    </Card>
  );
}
