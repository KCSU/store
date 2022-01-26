import {
  Box,
  Button,
  CloseButton,
  Heading,
  HStack,
  Icon,
  Stat,
  StatHelpText,
  StatLabel,
  StatNumber,
  Text,
  VStack,
} from "@chakra-ui/react";
import { FaPlus } from "react-icons/fa";
import { formatMoney } from "../../helpers/formatMoney";
import { QueueRequestAction } from "../../hooks/useQueueRequest";
import { Formal } from "../../model/Formal";
import { QueueRequest } from "../../model/QueueRequest";
import { TicketOptions } from "./TicketOptions";

interface TicketBuyFormProps {
  formal: Formal;
  hasShadow?: boolean;
  value: QueueRequest;
  onChange?: React.Dispatch<QueueRequestAction>;
}

export function TicketBuyForm({
  formal,
  onChange = () => {},
  value,
  hasShadow = true,
}: TicketBuyFormProps) {
  const ticket = value.ticket;
  const guestTickets = value.guestTickets;
  
  return (
    <VStack spacing={2}>
      <Box>
        <Text as="b" fontSize="lg">
          Buying tickets for {formal.name}:
        </Text>
      </Box>
      <TicketOptions
        hasShadow={hasShadow}
        value={ticket.option}
        onChange={value => onChange({type: 'option', value})}
      >
        <Heading as="h5" size="sm" mb={2}>
          King's Ticket: {formatMoney(formal.price)}
        </Heading>
      </TicketOptions>
      {guestTickets.map((t, i) => (
        <TicketOptions
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
        </TicketOptions>
      ))}
      {formal.guestLimit > 0 && (
        <Button
          disabled={guestTickets.length >= formal.guestLimit}
          onClick={() => onChange({type: 'addGuestTicket'})}
          size="sm"
          variant="outline"
          leftIcon={<Icon as={FaPlus} />}
          colorScheme="brand"
        >
          Add Guest Ticket
        </Button>
      )}
      <Stat textAlign="center">
        <StatLabel>Overall ticket price:</StatLabel>
        <StatNumber>
          {formatMoney(formal.price + formal.guestPrice * guestTickets.length)}
        </StatNumber>
        <StatHelpText>Will be added to college bill</StatHelpText>
      </Stat>
    </VStack>
  );
}
