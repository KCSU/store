import { Container, Heading } from "@chakra-ui/react";
import { useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { GroupDetailsForm } from "../../components/admin/GroupDetailsForm";
import { BackButton } from "../../components/utility/BackButton";
import { Card } from "../../components/utility/Card";
import { useCreateGroup } from "../../hooks/admin/useCreateGroup";
import { Group } from "../../model/Group";

export function AdminCreateGroupView() {
  const mutation = useCreateGroup();
  const navigate = useNavigate();
  const defaultGroup: Group = {
    id: 0,
    name: '',
    type: 'inst',
    lookup: ''
  }
  return (
    <Container maxW="container.md" p={0}>
      <BackButton to="/admin/groups">Back Home</BackButton>
      <Card mb={5}>
        <Heading as="h3" size="lg" mb={4}>
          Create a Group
        </Heading>
        <GroupDetailsForm group={defaultGroup} onSubmit={async (values) => {
          await mutation.mutateAsync(values);
          navigate('/admin/groups');
        }}>
          Create Group
        </GroupDetailsForm>
      </Card>
    </Container>
  );
}
