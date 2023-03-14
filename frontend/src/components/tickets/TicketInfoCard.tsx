import {
  Badge,
  Box,
  Button,
  Heading,
  HStack,
  Icon,
  IconButton,
  Table,
  Tbody,
  Td,
  Text,
  Tfoot,
  Th,
  Thead,
  Tooltip,
  Tr,
} from "@chakra-ui/react";
import { FaEdit, FaQrcode } from "react-icons/fa";
import { Link } from "react-router-dom";
import { formatMoney } from "../../helpers/formatMoney";
import { useCanEditTicket } from "../../hooks/state/useTicketPermissions";
import { useDateTime } from "../../hooks/state/useDateTime";
import { FormalTicket } from "../../model/Ticket";
import { Card } from "../utility/Card";
import { CancelTicketButton } from "./CancelTicketButton";
import { QrCodeModal } from "./QrCodeModal";

interface TicketInfoCardProps {
  ticket: FormalTicket;
  queue?: boolean;
}

// TODO: make this responsive!
export function TicketInfoCard({ ticket, queue = false }: TicketInfoCardProps) {
  const price =
    ticket.formal.price + ticket.formal.guestPrice * ticket.guestTickets.length;
  const canEdit = useCanEditTicket(ticket.formal);
  // const borderColor = useColorModeValue("gray.300", "gray.600")
  const datetime = useDateTime(ticket.formal.dateTime);
  return (
    <Card
      p={3}
      borderRadius="md"
      //borderWidth="1px" boxShadow="none" borderRadius="md" p={3}
    >
      <HStack>
        <Heading size="md" as="h4">
          {ticket.formal.name}
        </Heading>
        {queue ? (
          <Badge colorScheme="brand">In Queue</Badge>
        ) : (
          <Badge colorScheme="green">Confirmed</Badge>
        )}
      </HStack>
      <Text as="b">{datetime}</Text>
      <Box overflowX="auto" whiteSpace="nowrap">
        <Table size="sm" my={2} minW={{ base: "400px", sm: "0" }}>
          <Thead>
            <Tr>
              <Th>Type</Th>
              <Th>Meal Option</Th>
              <Th isNumeric minW={20} pl={0}>
                Price
              </Th>
            </Tr>
          </Thead>
          <Tbody>
            <Tr>
              <Td>King's Ticket</Td>

              <Td
                maxW={8}
                textOverflow="ellipsis"
                overflow="hidden"
                whiteSpace="nowrap"
              >
                <Tooltip label={ticket.ticket.option}>
                  {ticket.ticket.option}
                </Tooltip>
              </Td>
              <Td isNumeric minW={20} pl={0}>
                {formatMoney(ticket.formal.price)}
              </Td>
            </Tr>
            {ticket.guestTickets.map((gt, j) => {
              return (
                <Tr key={j}>
                  <Td>Guest Ticket</Td>
                  <Td
                    maxW={8}
                    textOverflow="ellipsis"
                    overflow="hidden"
                    whiteSpace="nowrap"
                  >
                    <Tooltip label={gt.option}>{gt.option}</Tooltip>
                  </Td>
                  <Td isNumeric minW={20} pl={0}>
                    {formatMoney(ticket.formal.guestPrice)}
                  </Td>
                </Tr>
              );
            })}
          </Tbody>
          <Tfoot fontWeight="bold">
            <Tr>
              <Td border="none">
                <Text fontSize="md">Total</Text>
              </Td>
              <Td isNumeric border="none" colSpan={2}>
                <Text fontSize="md">{formatMoney(price)}</Text>
              </Td>
            </Tr>
          </Tfoot>
        </Table>
      </Box>
      <HStack justifyContent="flex-end">
        {canEdit && (
          <>
            <Button
              size="sm"
              variant="outline"
              leftIcon={<FaEdit />}
              as={Link}
              to={`/tickets/${ticket.ticket.id}`}
            >
              Edit
            </Button>
            <CancelTicketButton
              formalId={ticket.formal.id}
              confirmText={`Cancel ${queue ? " Request" : " Ticket"}`}
            />
          </>
        )}
        {!queue && <QrCodeModal ticket={ticket} />}
      </HStack>
      {/* {queue && (
        <Progress
          colorScheme="brand"
          borderRadius={3}
          size="sm"
          isIndeterminate
          mt={3}
        />
      )} */}
    </Card>
  );
}
