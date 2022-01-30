import { Button, Heading, HStack, IconButton, Stat, StatHelpText, StatLabel, StatNumber, VStack } from "@chakra-ui/react";
import { FaSave, FaTrashAlt, FaUndo } from "react-icons/fa";
import { formatMoney } from "../../helpers/formatMoney";
import { FormalTicket } from "../../model/Ticket";
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
  return (
    <VStack spacing={3}>
      <TicketOptions
        hasShadow={hasShadow}
        value={ticket.option}
        footer={<Footer />}
      >
        <Heading as="h5" size="sm" mb={2}>
          King's Ticket: {formatMoney(formal.price)}
        </Heading>
      </TicketOptions>
      {guestTickets.map((t, i) => (
        <TicketOptions
          key={`guestTickets.${i}`}
          hasShadow={hasShadow}
          // TODO: state
          value={t.option}
          footer={<Footer />}
        >
          <HStack justify="space-between" mb={2}>
            <Heading as="h5" size="sm">
              Guest Ticket: {formatMoney(formal.guestPrice)}
            </Heading>
            <IconButton
              colorScheme="red"
              icon={<FaTrashAlt />}
              aria-label="Cancel Ticket"
              size="sm"
              //   onClick={() => onChange({ type: "removeGuestTicket", index: i })}
            ></IconButton>
          </HStack>
        </TicketOptions>
      ))}
    <PriceStat formal={formal} guestTickets={guestTickets}/>
    </VStack>
  );
}

function Footer() {
  return <HStack mt={3} justify="flex-end">
    <Button size="sm" variant="ghost" leftIcon={<FaUndo />}>
        Reset
      </Button>
    <Button size="sm" colorScheme="brand" leftIcon={<FaSave />}>
      Save Changes
    </Button>
  </HStack>
}