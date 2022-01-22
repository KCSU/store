import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
  Badge,
  Box,
  Button,
  Heading,
  HStack,
  Progress,
  Table,
  Tbody,
  Td,
  Text,
  Tfoot,
  Th,
  Thead,
  Tooltip,
  Tr,
  useDisclosure,
} from "@chakra-ui/react";
import { useRef } from "react";
import { FaEdit, FaTrashAlt } from "react-icons/fa";
import { formatMoney } from "../../helpers/formatMoney";
import { useCancelTickets } from "../../hooks/useCancelTickets";
import { useDateTime } from "../../hooks/useDateTime";
import { FormalTicket } from "../../model/Ticket";
import { Card } from "../utility/Card";

interface CancelTicketButtonProps {
  formalId: number;
  isQueue: boolean;
}

function CancelTicketButton({ isQueue, formalId }: CancelTicketButtonProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = useRef(null);
  const mutation = useCancelTickets();

  return (
    <>
      <Button
        size="sm"
        colorScheme="red"
        leftIcon={<FaTrashAlt />}
        onClick={onOpen}
      >
        Cancel {isQueue ? " Request" : " Ticket"}
      </Button>

      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Cancel Ticket
            </AlertDialogHeader>

            <AlertDialogBody>
              Are you sure you want to cancel your ticket?
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button
                colorScheme="red"
                onClick={async () => {
                  await mutation.mutateAsync(formalId);
                  onClose();
                }}
                isLoading={mutation.isLoading}
              >
                Cancel {isQueue ? " Request" : " Ticket"}
              </Button>
              <Button
                ref={cancelRef}
                onClick={onClose}
                ml={3}
                isDisabled={mutation.isLoading}
              >
                Go Back
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  );
}

interface TicketOverviewProps {
  ticket: FormalTicket;
  queue?: boolean;
}

// TODO: make this responsive!
export function TicketOverview({ ticket, queue = false }: TicketOverviewProps) {
  const price =
    ticket.formal.price + ticket.formal.guestPrice * ticket.guestTickets.length;
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
        {queue && <Badge colorScheme="brand">In Queue</Badge>}
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
        <Button size="sm" variant="outline" leftIcon={<FaEdit />}>
          Edit
        </Button>
        <CancelTicketButton formalId={ticket.formal.id} isQueue={queue} />
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
