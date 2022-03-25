import {
  Box,
  Divider,
  Heading,
  VStack,
} from "@chakra-ui/react";
import { useRoles } from "../../hooks/admin/useRoles";
import { Role } from "../../model/Role";
import { PermissionsTable } from "./PermissionsTable";
export function RolesList() {
  const { data, isLoading, isError } = useRoles();
  // TODO: loading states
  if (isLoading || isError || !data) {
    return <></>;
  }
  return (
    <>
      {data.map((role) => (
        <VStack align="stretch" key={role.id} borderWidth={1} borderRadius="md" p={2}>
          <Heading size="sm" as="h4" m={1}>
            {role.name}
          </Heading>
          <Divider my={2}/>
          <PermissionsTable role={role} />
        </VStack>
      ))}
    </>
  );
}
