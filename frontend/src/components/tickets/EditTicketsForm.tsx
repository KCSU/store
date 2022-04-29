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
import { useCanEditTicket, useTicketPermissions } from "../../hooks/state/useTicketPermissions";
import { formatMoney } from "../../helpers/formatMoney";
import { getBuyText } from "../../helpers/getBuyText";
import { useAddTicket } from "../../hooks/mutations/useAddTicket";
import { useEditTicket } from "../../hooks/mutations/useEditTicket";
import { Formal } from "../../model/Formal";
import { FormalTicket, Ticket } from "../../model/Ticket";
import { CancelGuestTicketDialog } from "./CancelGuestTicketDialog";
import { CancelTicketButton } from "./CancelTicketButton";
import { PriceStat } from "../formals/PriceStat";
import { TicketOptionsInput } from "./TicketOptionsInput";

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
      <EditSingleTicketForm formal={formal} ticket={ticket} hasShadow={hasShadow} isDisabled={!canEdit} />
      {guestTickets.map((t, i) => (
        <EditSingleTicketForm
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
      <AddGuestTicketModal isOpen={isOpen} onClose={onClose} formal={formal} />
    </VStack>
  );
}

interface EditSingleTicketFormProps {
  formal: Formal;
  ticket: Ticket;
  hasShadow?: boolean;
  isDisabled?: boolean;
}

function EditSingleTicketForm({
  formal,
  ticket,
  hasShadow,
  isDisabled = false
}: EditSingleTicketFormProps) {
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
    <TicketOptionsInput
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
            <CancelGuestTicketDialog
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
    </TicketOptionsInput>
  );
}

interface AddGuestTicketModalProps {
  isOpen: boolean;
  onClose: () => void;
  formal: Formal;
}

function AddGuestTicketModal({ isOpen, onClose, formal }: AddGuestTicketModalProps) {
  const mutation = useAddTicket(formal.id);
  const [option, setOption] = useState("Normal");
  const modalBg = useColorModeValue("gray.50", "gray.800");
  const isDisabled = !useTicketPermissions(formal).canBuy;
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent bg={modalBg}>
        <ModalHeader>Add Guest Ticket</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <TicketOptionsInput
            hasShadow={true}
            value={option}
            onChange={setOption}
          ></TicketOptionsInput>
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
