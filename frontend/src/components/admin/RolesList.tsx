import { Heading, LinkBox, LinkOverlay, SimpleGrid, Text, useColorModeValue } from "@chakra-ui/react";
import { Link } from "react-router-dom";
import { useRoles } from "../../hooks/admin/useRoles";
import { Card } from "../utility/Card";
export function RolesList() {
  const { data, isLoading, isError } = useRoles();
  const hoverBg = useColorModeValue('gray.100', 'gray.750');
  // TODO: loading states
  if (isLoading || isError || !data) {
    return <></>;
  }
  return (
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
  );
}
