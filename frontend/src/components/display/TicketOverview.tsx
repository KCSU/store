import {
  Button,
  Heading,
  HStack,
  Table,
  Tbody,
  Td,
  Text,
  Tfoot,
  Th,
  Thead,
  Tr,
  useColorModeValue,
} from "@chakra-ui/react";
import { FaEdit, FaTrashAlt } from "react-icons/fa";
import { formatMoney } from "../../helpers/formatMoney";
import { Ticket } from "../../model/Ticket";
import { Card } from "../utility/Card";

interface TicketOverviewProps {
  ticket: Ticket;
}

// TODO: make this responsive!
export function TicketOverview({ ticket }: TicketOverviewProps) {
  const price =
    ticket.formal.price + ticket.formal.guestPrice * ticket.guestTickets.length;
    // const borderColor = useColorModeValue("gray.300", "gray.600")
  return (
    <Card borderWidth="1px" boxShadow="none" borderRadius="md" p={3}>
      <Heading size="md" as="h4" mb={2}>
        {ticket.formal.title}
      </Heading>
      <Table size="sm" mb={2}>
        <Thead>
          <Tr>
            <Th>Type</Th>
            <Th>Meal Option</Th>
            <Th isNumeric>Price</Th>
          </Tr>
        </Thead>
        <Tbody>
          <Tr>
            <Td>King's Ticket</Td>
            <Td>{ticket.ticket.option}</Td>
            <Td isNumeric>{formatMoney(ticket.formal.price)}</Td>
          </Tr>
          {ticket.guestTickets.map((gt, j) => {
            return (
              <Tr key={j}>
                <Td>Guest Ticket</Td>
                <Td>{gt.option}</Td>
                <Td isNumeric>{formatMoney(ticket.formal.guestPrice)}</Td>
              </Tr>
            );
          })}
        </Tbody>
        <Tfoot fontWeight="bold">
          <Tr>
            <Td border="none">
              <Text fontSize="md">Total</Text>
            </Td>
            <Td border="none"></Td>
            <Td isNumeric border="none">
              <Text fontSize="md">{formatMoney(price)}</Text>
            </Td>
          </Tr>
        </Tfoot>
      </Table>
      <HStack justifyContent="flex-end">
        <Button size="sm" variant="outline" leftIcon={<FaEdit />}>
          Edit
        </Button>
        <Button size="sm" colorScheme="red" leftIcon={<FaTrashAlt />}>
          Cancel Ticket
        </Button>
      </HStack>
    </Card>
  );
}
