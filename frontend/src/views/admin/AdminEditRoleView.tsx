import { Container, Heading } from "@chakra-ui/react";
import { useParams } from "react-router-dom";
import { PermissionsTable } from "../../components/admin/PermissionsTable";
import { BackButton } from "../../components/utility/BackButton";
import { Card } from "../../components/utility/Card";
import { useRoles } from "../../hooks/admin/useRoles";

export function AdminEditRoleView() {
  const {id} = useParams();
  const roleId = parseInt(id ?? "0");
  const { data, isLoading, isError } = useRoles();
  // TODO: loading states
  if (isLoading || isError || !data) {
    return <></>;
  }
  const role = data.find(r => r.id === roleId);
  if (!role) {
    return <></>;
  }
  return <>
    <Container maxW="container.md" p={0}>
      <BackButton to="/admin/roles">Back Home</BackButton>
      <Card mb={5}>
        <Heading as="h3" size="lg" mb={4}>
          {role.name}
        </Heading>
        <PermissionsTable role={role} />
      </Card>
    </Container>
  </>
}
