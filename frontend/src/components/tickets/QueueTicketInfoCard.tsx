import {
  HStack,
  Heading,
  Badge,
  Text,
  Button,
  useDisclosure,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
} from "@chakra-ui/react";
import { useState } from "react";
import { FaEdit, FaSave, FaTrashAlt } from "react-icons/fa";
import { formatMoney } from "../../helpers/formatMoney";
import { useDateTime } from "../../hooks/state/useDateTime";
import { useEditTicket } from "../../hooks/mutations/useEditTicket";
import { QueueTicket } from "../../model/Queue";
import { Card } from "../utility/Card";
import { CancelGuestTicketDialog } from "./CancelGuestTicketDialog";
import { TicketOptionsInput } from "./TicketOptionsInput";

interface QueueTicketInfoCardProps {
  ticket: QueueTicket;
}

export function QueueTicketInfoCard({ ticket }: QueueTicketInfoCardProps) {
  const datetime = useDateTime(ticket.formal.dateTime);
  const { isOpen, onOpen, onClose } = useDisclosure();
  // const modalBg = useColorModeValue("gray.50", "gray.800");
  return (
    <Card p={3} borderRadius="md">
      <HStack>
        <Heading size="md" as="h4">
          {ticket.formal.name}
        </Heading>
        <Badge colorScheme="teal">Guest</Badge>
        <Badge colorScheme="brand">In Queue</Badge>
      </HStack>
      <Text as="b">{datetime}</Text>
      <Text fontSize="sm">Meal Option: {ticket.ticket.option}</Text>
      <HStack justifyContent="flex-end">
        <Text as="b">{formatMoney(ticket.formal.price)}</Text>
        <Button
          size="xs"
          variant="outline"
          leftIcon={<FaEdit />}
          onClick={onOpen}
        >
          Edit
        </Button>
        <CancelGuestTicketButton ticketId={ticket.ticket.id} />
        {/* <CancelTicketButton formalId={ticket.formal.id} isQueue={queue} /> */}
      </HStack>
      <EditQueueTicketModal isOpen={isOpen} onClose={onClose} ticket={ticket} />
    </Card>
  );
}

interface EditQueueTicketModalProps {
  isOpen: boolean;
  ticket: QueueTicket;
  onClose: () => void;
}

function EditQueueTicketModal({ isOpen, onClose, ticket }: EditQueueTicketModalProps) {
  const [option, setOption] = useState(ticket.ticket.option);
  const mutation = useEditTicket(ticket.ticket.id);

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Edit Guest Ticket</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <Heading as="h4" size="sm">
            Editing ticket for "{ticket.formal.name}":
          </Heading>
          <TicketOptionsInput
            hasShadow={false}
            value={option}
            onChange={setOption}
          />
        </ModalBody>
        <ModalFooter>
          <Button
            isLoading={mutation.isLoading}
            leftIcon={<FaSave />}
            colorScheme="brand"
            size="sm"
            mr={3}
            onClick={async () => {
              await mutation.mutateAsync(option);
              onClose();
            }}
          >
            Save Changes
          </Button>
          <Button
            variant="ghost"
            size="sm"
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

interface CancelGuestTicketButtonProps {
  ticketId: number;
}

function CancelGuestTicketButton({ ticketId }: CancelGuestTicketButtonProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <Button
        size="xs"
        colorScheme="red"
        leftIcon={<FaTrashAlt />}
        onClick={onOpen}
      >
        Cancel
      </Button>
      <CancelGuestTicketDialog isOpen={isOpen} onClose={onClose} ticketId={ticketId}/>
    </>
  );
}

