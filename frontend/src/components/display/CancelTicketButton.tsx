import {
  useDisclosure,
  Button,
  AlertDialog,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogBody,
  AlertDialogFooter,
} from "@chakra-ui/react";
import { useRef } from "react";
import { FaTrashAlt } from "react-icons/fa";
import { useCancelTickets } from "../../hooks/useCancelTickets";

export interface CancelTicketButtonProps {
  formalId: number;
  size?: string;
  confirmText: string | React.ReactNode;
  onSuccess?: () => void;
}

export function CancelTicketButton({
  formalId,
  size="sm",
  confirmText,
  onSuccess = () => {}
}: CancelTicketButtonProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = useRef(null);
  const mutation = useCancelTickets();

  return (
    <>
      <Button
        size={size}
        colorScheme="red"
        leftIcon={<FaTrashAlt />}
        onClick={onOpen}
      >
        {confirmText}
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
                  await mutation.mutateAsync(formalId);
                  onClose();
                  onSuccess();
                }}
                isLoading={mutation.isLoading}
              >
                {confirmText}
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
