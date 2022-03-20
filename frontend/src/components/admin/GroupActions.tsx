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
  HStack,
  Icon,
  useDisclosure,
} from "@chakra-ui/react";
import { useRef } from "react";
import { FaSync, FaTrashAlt } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { useDeleteGroup } from "../../hooks/admin/useDeleteGroup";
import { Group } from "../../model/Group";

interface GroupProps {
  group: Group;
}

export function GroupActions({ group }: GroupProps) {
  const navigate = useNavigate();
  const deleteMutation = useDeleteGroup(group.id);
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = useRef(null);
  return (
    <>
      <Alert status="warning" mb={4} variant="left-accent">
        <AlertIcon />
        These actions are potentially destructive! Only proceed if you know what
        you're doing.
      </Alert>
      <HStack justify="space-evenly" wrap="wrap" rowGap={2}>
        {group.type !== "manual" && (
          <Button
            colorScheme="brand"
            variant="outline"
            leftIcon={<Icon as={FaSync} />}
          >
            Sync with Lookup Directory
          </Button>
        )}
        <Button
          colorScheme="red"
          leftIcon={<Icon as={FaTrashAlt} />}
          onClick={onOpen}
        >
          Delete Group
        </Button>

        <AlertDialog
          isOpen={isOpen}
          leastDestructiveRef={cancelRef}
          onClose={onClose}
        >
          <AlertDialogOverlay>
            <AlertDialogContent>
              <AlertDialogHeader fontSize="lg" fontWeight="bold">
                Delete Group
              </AlertDialogHeader>

              <AlertDialogBody>
                Are you sure? This can't be undone.
              </AlertDialogBody>

              <AlertDialogFooter>
                <Button ref={cancelRef} onClick={onClose}>
                  Cancel
                </Button>
                <Button
                  colorScheme="red"
                  onClick={async () => {
                    await deleteMutation.mutateAsync();
                    onClose();
                    navigate('/admin/groups');
                  }}
                  ml={3}
                  isLoading={deleteMutation.isLoading}
                >
                  Delete
                </Button>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialogOverlay>
        </AlertDialog>
      </HStack>
    </>
  );
}
