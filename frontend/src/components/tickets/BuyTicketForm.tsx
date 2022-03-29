import {
  Box,
  Button,
  CloseButton,
  Heading,
  HStack,
  Icon,
  Text,
  VStack,
} from "@chakra-ui/react";
import { FaPlus } from "react-icons/fa";
import { formatMoney } from "../../helpers/formatMoney";
import { QueueRequestAction } from "../../hooks/state/useQueueRequest";
import { Formal } from "../../model/Formal";
import { QueueRequest } from "../../model/QueueRequest";
import { PriceStat } from "../formals/PriceStat";
import { TicketOptionsInput } from "./TicketOptionsInput";

interface BuyTicketFormProps {
  formal: Formal;
  hasShadow?: boolean;
  value: QueueRequest;
  onChange?: React.Dispatch<QueueRequestAction>;
}

export function BuyTicketForm({
  formal,
  onChange = () => {},
  value,
  hasShadow = true,
}: BuyTicketFormProps) {
  const ticket = value.ticket;
  const guestTickets = value.guestTickets;
  
  return (
    <VStack spacing={2}>
      <Box>
        <Text as="b" fontSize="lg">
          Buying tickets for {formal.name}:
        </Text>
      </Box>
      <TicketOptionsInput
        hasShadow={hasShadow}
        value={ticket.option}
        onChange={value => onChange({type: 'option', value})}
      >
        <Heading as="h5" size="sm" mb={2}>
          King's Ticket: {formatMoney(formal.price)}
        </Heading>
      </TicketOptionsInput>
      {guestTickets.map((t, i) => (
        <TicketOptionsInput
          key={`guestTickets.${i}`}
          hasShadow={hasShadow}
          value={t.option}
          onChange={value => onChange({type: 'guestTicket', index: i, value})}
        >
          <HStack justify="space-between" mb={2}>
            <Heading as="h5" size="sm">
              Guest Ticket: {formatMoney(formal.guestPrice)}
            </Heading>
            <CloseButton
              size="sm"
              onClick={() => onChange({type: 'removeGuestTicket', index: i})}
            ></CloseButton>
          </HStack>
        </TicketOptionsInput>
      ))}
      {formal.guestLimit > 0 && (
        <Button
          isDisabled={guestTickets.length >= formal.guestLimit}
          onClick={() => onChange({type: 'addGuestTicket'})}
          size="sm"
          variant="outline"
          leftIcon={<Icon as={FaPlus} />}
          colorScheme="brand"
        >
          Add Guest Ticket
        </Button>
      )}
      <PriceStat formal={formal} guestTickets={guestTickets}/>
    </VStack>
  );
}
