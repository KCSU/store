import {
  Button,
  AlertDialog,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogBody,
  AlertDialogFooter,
} from "@chakra-ui/react";
import { useRef } from "react";
import { useCancelTicket } from "../../hooks/mutations/useCancelTicket";

export interface CancelGuestTicketDialogProps {
  isOpen: boolean;
  ticketId: string;
  confirmText?: string;
  onClose: () => void;
}

export function CancelGuestTicketDialog({isOpen, onClose, ticketId, confirmText}: CancelGuestTicketDialogProps) {
  const mutation = useCancelTicket();
  const cancelRef = useRef(null);
  return (
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
            Are you sure you want to cancel this ticket?
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
              {confirmText ?? 'Cancel Request'}
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
  );
}
