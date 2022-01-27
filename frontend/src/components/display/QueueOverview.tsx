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
  AlertDialog,
  AlertDialogBody,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
} from "@chakra-ui/react";
import { useRef, useState } from "react";
import { FaEdit, FaSave, FaTrashAlt } from "react-icons/fa";
import { formatMoney } from "../../helpers/formatMoney";
import { useCancelTicket } from "../../hooks/useCancelTicket";
import { useDateTime } from "../../hooks/useDateTime";
import { useEditTicket } from "../../hooks/useEditTicket";
import { QueueTicket } from "../../model/Queue";
import { Card } from "../utility/Card";
import { TicketOptions } from "./TicketOptions";

interface QueueOverviewProps {
  ticket: QueueTicket;
}

export function QueueOverview({ ticket }: QueueOverviewProps) {
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
        <CancelQueueButton ticketId={ticket.ticket.id} />
        {/* <CancelTicketButton formalId={ticket.formal.id} isQueue={queue} /> */}
      </HStack>
      <EditQueueTicket isOpen={isOpen} onClose={onClose} ticket={ticket} />
    </Card>
  );
}

interface EditQueueTicketProps {
  isOpen: boolean;
  ticket: QueueTicket;
  onClose: () => void;
}

function EditQueueTicket({ isOpen, onClose, ticket }: EditQueueTicketProps) {
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
          <TicketOptions
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

interface CancelQueueButtonProps {
  ticketId: number;
}

function CancelQueueButton({ ticketId }: CancelQueueButtonProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = useRef(null);
  const mutation = useCancelTicket();

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
                  await mutation.mutateAsync(ticketId);
                  onClose();
                }}
                isLoading={mutation.isLoading}
              >
                Cancel Request
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
