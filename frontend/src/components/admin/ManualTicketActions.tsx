import {
  Button,
  Flex,
  Icon,
  IconButton,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  useDisclosure,
} from "@chakra-ui/react";
import { FaTrashAlt } from "react-icons/fa";
import { useCancelManualTicket } from "../../hooks/admin/useCancelManualTicket";
import { ManualTicket } from "../../model/ManualTicket";

interface ManualTicketProps {
  ticket: ManualTicket;
}

function CancelManualTicketButton({ ticket }: ManualTicketProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const mutation = useCancelManualTicket(ticket.id);
  return (
    <>
      <IconButton
        variant="ghost"
        colorScheme="red"
        size="xs"
        onClick={onOpen}
        aria-label="Delete"
        icon={<Icon as={FaTrashAlt} />}
      ></IconButton>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Cancel Ticket</ModalHeader>
          <ModalCloseButton />
          <ModalBody>Are you sure you want to cancel this ticket?</ModalBody>
          <ModalFooter>
            <Button variant="ghost" onClick={onClose}>
              Close
            </Button>
            <Button
              colorScheme="red"
              isLoading={mutation.isLoading}
              onClick={async () => {
                await mutation.mutateAsync();
                onClose();
              }}
            >
              Cancel Ticket
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}

interface ManualTicketActionsProps {
  ticket: ManualTicket;
  canDelete: boolean;
  canWrite: boolean;
}

export function ManualTicketActions({
  canDelete,
  canWrite,
  ticket,
}: ManualTicketActionsProps) {
  return (
    <Flex align="center" gap={2}>
      {/* {canWrite && <EditManualTicketButton ticket={ticket} />} */}
      {canDelete && <CancelManualTicketButton ticket={ticket} />}
    </Flex>
  );
}
