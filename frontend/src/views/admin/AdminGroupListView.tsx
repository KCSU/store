import { Box, Button, Code, Heading, LinkBox, LinkOverlay, SimpleGrid, Text } from "@chakra-ui/react";
import { FaPlus } from "react-icons/fa";
import { Link } from "react-router-dom";
import { Card } from "../../components/utility/Card";
import { useGroups } from "../../hooks/admin/useGroups";
import { Group, groupType } from "../../model/Group";

interface GroupProps {
  group: Group;
}

function AdminGroupCard({ group }: GroupProps) {
  return (
    <LinkBox
      as={Card}
      borderRadius={3}
      p={4}
      transition="box-shadow 0.2s"
      _hover={{ shadow: "lg" }}
      _focusWithin={{ shadow: "lg" }}
    >
      <LinkOverlay as={Link} to={`/admin/groups/${group.id}`}>
        <Heading as="h5" size="sm" mb={2}>
          {group.name}
        </Heading>
      </LinkOverlay>
      <Text fontSize="sm">
        {groupType(group)}
        {group.lookup && ': '}
        {group.lookup && <Code>
          {group.lookup}
        </Code>}
        {/* TODO: Member count? */}
      </Text>
    </LinkBox>
  );
}

export function AdminGroupListView() {
  const { data, isLoading, isError } = useGroups();
  if (!data) {
    return <></>;
  }

  return (
    <>
      <Heading size="xl" mb={5}>
        Manage Groups
      </Heading>
      <Button
        colorScheme="brand"
        mb={4}
        leftIcon={<FaPlus />}
        as={Link}
        to="/admin/groups/create"
      >
        Create Group
      </Button>
      <SimpleGrid
        templateColumns="repeat(auto-fill, minmax(250px, 1fr))"
        spacing="20px"
      >
        {data.map((g) => (
          <AdminGroupCard group={g} key={g.id} />
        ))}
      </SimpleGrid>
    </>
  );
}
