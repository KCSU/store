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
import { QueueRequest } from "../../model/QueueRequest";
import { TicketRequest } from "../../model/TicketRequest";
import { TicketOptions } from "./TicketOptions";

interface TicketBuyFormProps {
  formal: Formal;
  hasShadow?: boolean;
  value: QueueRequest;
  onChange?: (req: QueueRequest) => void;
}

export function TicketBuyForm({
  formal,
  onChange = _ => {},
  value,
  hasShadow = true,
}: TicketBuyFormProps) {
  const ticket = value.ticket;
  const guestTickets = value.guestTickets;

  // CALLBACKS
  const setTicket = (t: TicketRequest) => {
    onChange({
      ...value,
      ticket: t
    });
  }
  const setGuestTickets = (gt: TicketRequest[]) => {
    onChange({
      ...value,
      guestTickets: gt
    })
  }
  const setTicketOption = (option: string) => {
    setTicket({
      ...value.ticket,
      option,
    });
  };
  const setGuestTicket = (index: number, option: string) => {
    const prev = value.guestTickets;
    setGuestTickets([
      ...prev.slice(0, index),
      { option },
      ...prev.slice(index + 1),
    ]);
  };
  const addGuestTicket = () => {
    const prev = value.guestTickets;
    setGuestTickets([
      ...prev,
      {
        option: "Normal",
      },
    ]);
  };
  const removeGuestTicket = (index: number) => {
    const prev = value.guestTickets;
    setGuestTickets([
      ...prev.slice(0, index),
      ...prev.slice(index + 1),
    ]);
  };
  
  return (
    <VStack spacing={2}>
      <Box>
        <Text as="b" fontSize="lg">
          Buying tickets for {formal.name}:
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
