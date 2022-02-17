import {
  Badge,
  Box,
  Button,
  Heading,
  HStack,
  LinkBox,
  LinkOverlay,
  SimpleGrid,
  Text,
} from "@chakra-ui/react";
import { FaPlus } from "react-icons/fa";
import { Link } from "react-router-dom";
import { Card } from "../../components/utility/Card";
import { useAllFormals } from "../../hooks/admin/useAllFormals";
import { useDateTime } from "../../hooks/useDateTime";
import { Formal } from "../../model/Formal";

interface FormalProps {
  formal: Formal;
}

function FormalStatusTag({ formal }: FormalProps) {
  if (formal.saleEnd < new Date()) {
    return <Badge>Closed</Badge>;
  } else if (formal.saleStart > new Date()) {
    return <Badge colorScheme="teal">Queueing</Badge>;
  } else if (
    formal.guestTicketsRemaining === 0 &&
    formal.ticketsRemaining === 0
  ) {
    return <Badge colorScheme="red">Sold Out</Badge>;
  } else {
    return <Badge colorScheme="purple">Sales Open</Badge>;
  }
}

function FormalAdminCard({ formal: f }: FormalProps) {
  const dateTime = useDateTime(f.dateTime);
  return (
    <LinkBox
      as={Card}
      borderRadius={3}
      p={4}
      transition="box-shadow 0.2s"
      _hover={{ shadow: "lg" }}
      _focusWithin={{ shadow: "lg" }}
    >
      <HStack mb={2}>
        <LinkOverlay as={Link} to={`/admin/formals/${f.id}`}>
          <Heading as="h5" size="sm">
            {f.name}
          </Heading>
        </LinkOverlay>
        <FormalStatusTag formal={f} />
      </HStack>
      <Text fontSize="sm">{dateTime}</Text>
    </LinkBox>
  );
}

export function FormalList() {
  const { data, isLoading, isError } = useAllFormals();
  if (!data) {
    return <></>;
  }

  return (
    <>
      <Heading size="xl" mb={5}>
        Manage Formals
      </Heading>
      <Button colorScheme="brand" mb={4} leftIcon={<FaPlus />} as={Link} to="/admin/formals/create">
        Create Formal
      </Button>
      <SimpleGrid
        templateColumns="repeat(auto-fill, minmax(300px, 1fr))"
        spacing="20px"
      >
        {data.map((f) => (
          <FormalAdminCard formal={f} key={f.id} />
        ))}
      </SimpleGrid>
    </>
  );
}
