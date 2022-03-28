import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
  Button,
  Container,
  Flex,
  Heading,
  Icon,
  IconButton,
  useDisclosure,
} from "@chakra-ui/react";
import { useRef } from "react";
import { FaPen, FaTrashAlt } from "react-icons/fa";
import { useNavigate, useParams } from "react-router-dom";
import { PermissionsTable } from "../../components/admin/PermissionsTable";
import { BackButton } from "../../components/utility/BackButton";
import { Card } from "../../components/utility/Card";
import { useDeleteRole } from "../../hooks/admin/useDeleteRole";
import { useRoles } from "../../hooks/admin/useRoles";
import { Role } from "../../model/Role";

export function AdminEditRoleView() {
  const { id } = useParams();
  const roleId = parseInt(id ?? "0");
  const { data, isLoading, isError } = useRoles();
  // TODO: loading states
  if (isLoading || isError || !data) {
    return <></>;
  }
  const role = data.find((r) => r.id === roleId);
  if (!role) {
    return <></>;
  }
  return (
    <>
      <Container maxW="container.md" p={0}>
        <BackButton to="/admin/roles">Back Home</BackButton>
        <Card mb={5}>
          <Flex gap={3}>
            <Heading as="h3" size="lg" mb={4} flex="1">
              {role.name}
            </Heading>
            <IconButton
              size="sm"
              variant="outline"
              // colorScheme="brand"
              aria-label="Edit"
              icon={<Icon as={FaPen} />}
            ></IconButton>
            <DeleteRoleButton role={role} />
          </Flex>
          <PermissionsTable role={role} />
        </Card>
      </Container>
    </>
  );
}

interface RoleProps {
  role: Role;
}

function DeleteRoleButton({ role }: RoleProps) {
  const mutation = useDeleteRole(role.id);
  const { isOpen, onOpen, onClose } = useDisclosure();
  const navigate = useNavigate();
  const cancelRef = useRef(null);
  return (
    <>
      <IconButton
        onClick={onOpen}
        size="sm"
        colorScheme="red"
        aria-label="Delete"
        icon={<Icon as={FaTrashAlt} />}
      ></IconButton>
      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Role
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
                  await mutation.mutateAsync();
                  onClose();
                  navigate("/admin/roles");
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
