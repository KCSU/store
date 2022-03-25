import {
  Button,
  Heading,
  Icon,
  Input,
  LinkBox,
  LinkOverlay,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  SimpleGrid,
  Text,
  useColorModeValue,
  useDisclosure,
} from "@chakra-ui/react";
import { useState } from "react";
import { FaPlus } from "react-icons/fa";
import { Link } from "react-router-dom";
import { useCreateRole } from "../../hooks/admin/useCreateRole";
import { useRoles } from "../../hooks/admin/useRoles";
import { Card } from "../utility/Card";

interface CreateRoleModalProps {
  isOpen: boolean;
  onClose: () => void;
}

function CreateRoleModal({ isOpen, onClose }: CreateRoleModalProps) {
  const [name, setName] = useState('');
  const mutation = useCreateRole();
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Create Role</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <Input placeholder="Name" value={name} onChange={e => setName(e.target.value)} />
        </ModalBody>

        <ModalFooter>
          <Button colorScheme="brand" mr={3} isLoading={mutation.isLoading} onClick={async () => {
            await mutation.mutateAsync(name);
            setName('');
            onClose();
          }}>
            Create
          </Button>
          <Button variant="ghost" onClick={onClose}>
            Cancel
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}

export function RolesList() {
  const { data, isLoading, isError } = useRoles();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const hoverBg = useColorModeValue("gray.100", "gray.750");
  // TODO: loading states
  if (isLoading || isError || !data) {
    return <></>;
  }
  return (
    <>
      <Heading size="md" as="h3">
        Manage Roles
      </Heading>
      <Button
        onClick={onOpen}
        alignSelf="start"
        size="sm"
        colorScheme="brand"
        leftIcon={<Icon as={FaPlus} />}
      >
        Create Role
      </Button>
      <SimpleGrid
        templateColumns="repeat(auto-fill, minmax(200px, 1fr))"
        spacing={3}
      >
        {data.map((role) => (
          <LinkBox
            key={role.id}
            borderWidth={1}
            borderRadius="md"
            p={3}
            as={Card}
            transition="background-color 200ms"
            _hover={{ bgColor: hoverBg }}
            _focusWithin={{ bgColor: hoverBg }}
          >
            <LinkOverlay as={Link} to={`/admin/roles/${role.id}`}>
              <Heading size="sm" as="h4">
                {role.name}
              </Heading>
            </LinkOverlay>
            <Text fontSize="sm" mt={1}>
              {role.permissions?.length ?? 0} permission
              {role.permissions?.length !== 1 && "s"}
            </Text>
          </LinkBox>
        ))}
      </SimpleGrid>
      <CreateRoleModal isOpen={isOpen} onClose={onClose}/>
    </>
  );
}
