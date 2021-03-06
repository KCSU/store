import { Heading, SimpleGrid } from "@chakra-ui/react";
import { RolesList } from "../../components/admin/RolesList";
import { UserRolesList } from "../../components/admin/UserRolesList";
import { Card } from "../../components/utility/Card";

export function AdminPermissionsView() {
  return (
    <>
      <Heading size="xl" mb={5}>
        Permissions
      </Heading>
      <SimpleGrid columns={{base: 1, xl: 2}} alignItems="start" gap={5}>
        <Card gap={3}>
          <RolesList/>
        </Card>
        <Card gap={3}>
          <Heading size="md" as="h3">User Roles</Heading>
          <UserRolesList/>
        </Card>
      </SimpleGrid>
    </>
  );
}
