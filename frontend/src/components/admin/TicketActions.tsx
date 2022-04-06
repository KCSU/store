import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
  Button,
  Flex,
  Icon,
  IconButton,
  Link,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  useDisclosure,
} from "@chakra-ui/react";
import { useRef, useState } from "react";
import { FaTrashAlt, FaExternalLinkAlt, FaPen } from "react-icons/fa";
import { useCancelTicket } from "../../hooks/admin/useCancelTicket";
import { useEditTicket } from "../../hooks/admin/useEditTicket";
import { AdminTicket } from "../../model/Ticket";
import { TicketOptionsInput } from "../tickets/TicketOptionsInput";

interface TicketProps {
  ticket: AdminTicket;
}

function EditTicketButton({ ticket }: TicketProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const mutation = useEditTicket(ticket.id);
  const [option, setOption] = useState(ticket.option);
  return (
    <>
      <IconButton
        variant="ghost"
        // variant="outline"
        size="xs"
        onClick={onOpen}
        aria-label="Edit"
        icon={<Icon as={FaPen} />}
      ></IconButton>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Edit Ticket</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <TicketOptionsInput
              value={option}
              onChange={setOption}
              hasShadow={false}
            />
          </ModalBody>

          <ModalFooter>
            <Button
              colorScheme="brand"
              mr={3}
              isLoading={mutation.isLoading}
              onClick={async () => {
                await mutation.mutateAsync(option);
                onClose();
              }}
            >
              Save
            </Button>
            <Button variant="ghost" onClick={onClose}>
              Cancel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}

function CancelTicketButton({ ticket }: TicketProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const leastDestructiveRef = useRef(null);
  const mutation = useCancelTicket(ticket.id);
  return (
    <>
      <IconButton
        variant="ghost"
        size="xs"
        colorScheme="red"
        aria-label="Delete"
        onClick={onOpen}
      >
        <Icon as={FaTrashAlt} />
      </IconButton>
      <AlertDialog
        isOpen={isOpen}
        onClose={onClose}
        leastDestructiveRef={leastDestructiveRef}
      >
        <AlertDialogOverlay />
        <AlertDialogContent>
          <AlertDialogHeader fontSize="lg" fontWeight="bold">
            Cancel Ticket
          </AlertDialogHeader>
          <AlertDialogBody>
            Are you sure you want to cancel this ticket?
            {!ticket.isGuest && (
              <Text mt={2}>
                This will also cancel all associated guest tickets for the user{" "}
                {ticket.userName} (
                <Link color="teal.500" href={`mailto:${ticket.userEmail}`}>
                  {ticket.userEmail.split("@")[0]}
                  <Icon boxSize={3} ml={1} as={FaExternalLinkAlt} />
                </Link>
                ).
              </Text>
            )}
          </AlertDialogBody>
          <AlertDialogFooter>
            <Button ref={leastDestructiveRef} onClick={onClose}>
              Close
            </Button>
            <Button
              isLoading={mutation.isLoading}
              colorScheme="red"
              onClick={async () => {
                await mutation.mutateAsync();
                onClose();
              }}
              ml={3}
            >
              Cancel Ticket
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

interface TicketActionsProps {
  canDelete: boolean;
  canWrite: boolean;
  ticket: AdminTicket;
}

export function TicketActions({
  canDelete,
  canWrite,
  ticket,
}: TicketActionsProps) {
  return (
    <Flex align="center" gap={2}>
      {canWrite && <EditTicketButton ticket={ticket} />}
      {canDelete && <CancelTicketButton ticket={ticket} />}
    </Flex>
  );
}
