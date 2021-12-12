import {
  Box,
  Button,
  CloseButton,
  Heading,
  HStack,
  Icon,
  Stat,
  StatGroup,
  StatHelpText,
  StatLabel,
  StatNumber,
  Text,
  VStack,
} from "@chakra-ui/react";
import { useState } from "react";
import { FaPlus } from "react-icons/fa";
import { formatMoney } from "../../helpers/formatMoney";
import { Formal } from "../../model/Formal";
import { TicketRequest } from "../../model/TicketRequest";
import { TicketOptions } from "./TicketOptions";

interface TicketBuyFormProps {
  formal: Formal;
  hasShadow?: boolean;
}

export function TicketBuyForm({ formal, hasShadow = true }: TicketBuyFormProps) {
  // TODO: Change this
  const [ticket, setTicket] = useState<TicketRequest>({
    option: "Normal",
  });
  const [guestTickets, setGuestTickets] = useState<TicketRequest[]>([]);
  const setTicketOption = (value: string) => {
    setTicket((prev) => ({
      ...prev,
      option: value,
    }));
  };
  const setGuestTicket = (index: number, value: string) => {
    setGuestTickets((prev) => [
      ...prev.slice(0, index),
      {
        option: value,
      },
      ...prev.slice(index + 1),
    ]);
  };
  const addGuestTicket = () => {
    setGuestTickets((prev) => [
      ...prev,
      {
        option: "Normal",
      },
    ]);
  };
  const removeGuestTicket = (index: number) => {
    setGuestTickets((prev) => [
      ...prev.slice(0, index),
      ...prev.slice(index + 1),
    ]);
  };
  return (
    <VStack spacing={2}>
      <Box>
        <Text as="b" fontSize="lg">
          Buying tickets for {formal.title}:
        </Text>
      </Box>
      <TicketOptions
        options={formal.options}
        hasShadow={hasShadow}
        value={ticket.option}
        onChange={setTicketOption}
      >
        <Heading as="h5" size="sm" mb={2}>
          King's Ticket: {formatMoney(formal.price)}
        </Heading>
      </TicketOptions>
      {guestTickets.map((t, i) => (
        <TicketOptions
          key={`guestTickets.${i}`}
          options={formal.options}
          hasShadow={hasShadow}
          value={t.option}
          onChange={(v) => setGuestTicket(i, v)}
        >
          <HStack justify="space-between" mb={2}>
            <Heading as="h5" size="sm">
              Guest Ticket: {formatMoney(formal.guestPrice)}
            </Heading>
            <CloseButton
              size="sm"
              onClick={() => removeGuestTicket(i)}
            ></CloseButton>
          </HStack>
        </TicketOptions>
      ))}
      {formal.guestLimit > 0 && (
        <Button
          disabled={guestTickets.length >= formal.guestLimit}
          onClick={addGuestTicket}
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
