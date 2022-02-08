import {
  Badge,
  Box,
  Button,
  Heading,
  HStack,
  IconButton,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Tooltip,
  useColorModeValue,
  useDisclosure,
  VStack,
  Center,
} from "@chakra-ui/react";
import { useState } from "react";
import { FaPlus, FaSave, FaTrashAlt, FaUndo } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { useCanBuyTicket, useCanEditTicket } from "../../hooks/useCanBuyTicket";
import { formatMoney } from "../../helpers/formatMoney";
import { getBuyText } from "../../helpers/getBuyText";
import { useAddTicket } from "../../hooks/useAddTicket";
import { useEditTicket } from "../../hooks/useEditTicket";
import { Formal } from "../../model/Formal";
import { FormalTicket, Ticket } from "../../model/Ticket";
import { CancelGuestDialog } from "./CancelGuestDialog";
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
  const { isOpen, onOpen, onClose } = useDisclosure();
  const canEdit = useCanEditTicket(formal);
  return (
    <VStack spacing={3}>
      <SingleTicketForm formal={formal} ticket={ticket} hasShadow={hasShadow} isDisabled={!canEdit} />
      {guestTickets.map((t, i) => (
        <SingleTicketForm
          isDisabled={!canEdit}
          key={`guestTickets.${i}`}
          formal={formal}
          ticket={t}
          hasShadow={hasShadow}
        />
      ))}
      <HStack spacing={4}>
        <CancelTicketButton
          isDisabled={!canEdit}
          size="md"
          formalId={formal.id}
          confirmText="Cancel Tickets"
          title="Cancel Tickets"
          body={`Are you sure you want to cancel your tickets for "${formal.name}"?`}
          onSuccess={() => navigate("/tickets")}
        />
        <Button
          colorScheme="brand"
          leftIcon={<FaPlus />}
          isDisabled={guestTickets.length >= formal.guestLimit || !canEdit}
          onClick={onOpen}
        >
          Add Guest Ticket
        </Button>
      </HStack>
      <PriceStat formal={formal} guestTickets={guestTickets} />
      <AddGuestModal isOpen={isOpen} onClose={onClose} formal={formal} />
    </VStack>
  );
}

interface SingleTicketFormProps {
  formal: Formal;
  ticket: Ticket;
  hasShadow?: boolean;
  isDisabled?: boolean;
}

function SingleTicketForm({
  formal,
  ticket,
  hasShadow,
  isDisabled = false
}: SingleTicketFormProps) {
  const mutation = useEditTicket(ticket.id);
  const [option, setOption] = useState(ticket.option);
  const { isOpen, onOpen, onClose } = useDisclosure();
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
        isDisabled={option === ticket.option || isDisabled}
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
      // TODO: state, isDisabled
      isDisabled={isDisabled}
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
                isDisabled={isDisabled}
                icon={<FaTrashAlt />}
                aria-label="Cancel Ticket"
                size="sm"
                onClick={onOpen}
                //   onClick={() => onChange({ type: "removeGuestTicket", index: i })}
              ></IconButton>
            </Tooltip>
            <CancelGuestDialog
              confirmText="Cancel Ticket"
              isOpen={isOpen}
              onClose={onClose}
              ticketId={ticket.id}
            />
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

interface AddGuestModalProps {
  isOpen: boolean;
  onClose: () => void;
  formal: Formal;
}

function AddGuestModal({ isOpen, onClose, formal }: AddGuestModalProps) {
  const mutation = useAddTicket(formal.id);
  const [option, setOption] = useState("Normal");
  const modalBg = useColorModeValue("gray.50", "gray.800");
  const isDisabled = !useCanBuyTicket(formal);
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent bg={modalBg}>
        <ModalHeader>Add Guest Ticket</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <TicketOptions
            hasShadow={true}
            value={option}
            onChange={setOption}
          ></TicketOptions>
          <Center as="b" fontSize="lg" mt={4}>
            Ticket Price: {formatMoney(formal.guestPrice)}
          </Center>
        </ModalBody>

        <ModalFooter>
          <Button
            isLoading={mutation.isLoading}
            colorScheme="brand"
            mr={3}
            isDisabled={isDisabled}
            onClick={async () => {
              await mutation.mutateAsync({ option });
              setOption("Normal");
              onClose();
            }}
          >
            {getBuyText(formal)}
          </Button>
          <Button
            variant="ghost"
            onClick={onClose}
            isDisabled={mutation.isLoading}
          >
            Cancel
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}
