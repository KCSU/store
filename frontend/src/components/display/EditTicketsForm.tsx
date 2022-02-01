import {
  Badge,
  Box,
  Button,
  Flex,
  Heading,
  HStack,
  IconButton,
  Tooltip,
  VStack,
} from "@chakra-ui/react";
import { useCallback, useState } from "react";
import { FaPlus, FaSave, FaTrashAlt, FaUndo } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { formatMoney } from "../../helpers/formatMoney";
import { useEditTicket } from "../../hooks/useEditTicket";
import { Formal } from "../../model/Formal";
import { FormalTicket, Ticket } from "../../model/Ticket";
import { CancelTicketButton } from "./CancelTicketButton";
import { PriceStat } from "./PriceStat";
import { TicketOptions } from "./TicketOptions";

export interface EditTicketsFormProps {
  ticket: FormalTicket;
  hasShadow?: boolean;
}

export function EditTicketsForm({
  ticket: { formal, ticket, guestTickets },
  hasShadow,
}: EditTicketsFormProps) {
  const navigate = useNavigate();
  return (
    <VStack spacing={3}>
      <SingleTicketForm formal={formal} ticket={ticket} hasShadow={hasShadow} />
      {guestTickets.map((t, i) => (
        <SingleTicketForm
          key={`guestTickets.${i}`}
          formal={formal}
          ticket={t}
          hasShadow={hasShadow}
        />
      ))}
      <HStack spacing={4}>
        <CancelTicketButton
          size="md"
          formalId={formal.id}
          confirmText="Cancel Tickets"
          onSuccess={() => navigate("/tickets")}
        />
        <Button
          colorScheme="brand"
          leftIcon={<FaPlus />}
          isDisabled={guestTickets.length >= formal.guestLimit}
        >
          Add Guest Ticket
        </Button>
      </HStack>
      <PriceStat formal={formal} guestTickets={guestTickets} />
    </VStack>
  );
}

interface SingleTicketFormProps {
  formal: Formal;
  ticket: Ticket;
  hasShadow?: boolean;
}

function SingleTicketForm({
  formal,
  ticket,
  hasShadow,
}: SingleTicketFormProps) {
  const mutation = useEditTicket(ticket.id);
  const [option, setOption] = useState(ticket.option);
  const footer = (
    <HStack mt={3} justify="flex-end">
      <Button
        size="sm"
        variant="ghost"
        leftIcon={<FaUndo />}
        isDisabled={option === ticket.option}
        onClick={() => setOption(ticket.option)}
      >
        Reset
      </Button>
      <Button
        size="sm"
        colorScheme="brand"
        leftIcon={<FaSave />}
        isDisabled={option === ticket.option}
        isLoading={mutation.isLoading}
        onClick={() => mutation.mutate(option)}
      >
        Save Changes
      </Button>
    </HStack>
  );
  return (
    <TicketOptions
      hasShadow={hasShadow}
      // TODO: state
      value={option}
      onChange={setOption}
      footer={footer}
    >
      <HStack mb={2} align="start">
        {ticket.isGuest ? (
          <>
            <Heading as="h4" size="sm">
              Guest Ticket: {formatMoney(formal.guestPrice)}
            </Heading>
            {ticket.isQueue && <Badge colorScheme="brand">In Queue</Badge>}
            <Box flex="1"></Box>
            <Tooltip label="Cancel Ticket">
              <IconButton
                justifySelf="flex-end"
                colorScheme="red"
                icon={<FaTrashAlt />}
                aria-label="Cancel Ticket"
                size="sm"
                //   onClick={() => onChange({ type: "removeGuestTicket", index: i })}
              ></IconButton>
            </Tooltip>
          </>
        ) : (
          <>
            <Heading as="h4" size="sm" mb={2}>
              King's Ticket: {formatMoney(formal.price)}
            </Heading>
            {ticket.isQueue && <Badge colorScheme="brand">In Queue</Badge>}
          </>
        )}
      </HStack>
    </TicketOptions>
  );
}
