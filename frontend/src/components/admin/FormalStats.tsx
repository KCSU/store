import {
  Alert,
  AlertDialog,
  AlertDialogBody,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
  AlertIcon,
  Button,
  Heading,
  Icon,
  useDisclosure,
} from "@chakra-ui/react";
import { useRef } from "react";
import { FaTrashAlt } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { useDeleteFormal } from "../../hooks/admin/useDeleteFormal";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { Formal } from "../../model/Formal";

interface FormalProps {
  formal: Formal;
}

export function FormalStats({ formal }: FormalProps) {
  const canDelete = useHasPermission("formals", "delete");
  return <>
    {canDelete && <FormalActions formal={formal}/>}
  </>
}

function FormalActions({formal}: FormalProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = useRef(null);
  const mutation = useDeleteFormal(formal.id);
  const navigate = useNavigate();
  return (
    <>
      <Alert status="warning" mb={4} variant="left-accent">
        <AlertIcon />
        The following actions are potentially destructive! Only proceed if you
        know what you're doing.
      </Alert>
      <Button
        colorScheme="red"
        leftIcon={<Icon as={FaTrashAlt} />}
        onClick={onOpen}
      >
        Delete Formal
      </Button>
      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Formal
            </AlertDialogHeader>
            <AlertDialogBody>
              Are you sure you want to delete this formal? This action cannot be
              undone.
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={onClose}>
                Cancel
              </Button>
              <Button
                colorScheme="red"
                onClick={async () => {
                  await mutation.mutateAsync();
                  onClose();
                  navigate("/admin/formals");
                }}
                ml={3}
                isLoading={mutation.isLoading}
              >
                Delete
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  );
}
