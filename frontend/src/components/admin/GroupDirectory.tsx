import {
  Heading,
  Icon,
  Input,
  InputGroup,
  InputLeftAddon,
  VStack,
  Box,
  HStack,
  Button,
  useColorModeValue,
} from "@chakra-ui/react";
import { useMemo, useState } from "react";
import { FaArrowRight, FaSearch } from "react-icons/fa";
import { useAddGroupUser } from "../../hooks/admin/useAddGroupUser";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { useRemoveGroupUser } from "../../hooks/admin/useRemoveGroupUser";
import { Group, GroupUser } from "../../model/Group";
import { GroupUserList } from "./GroupUserList";

interface GroupProps {
  group: Group;
}

function AddUserBox({group}: GroupProps) {
  const [email, setEmail] = useState("");
  const inputBg = useColorModeValue("gray.50", "gray.600");
  const mutation = useAddGroupUser(group.id);
  return <Box
    align="center"
    borderWidth={1}
    alignSelf="center"
    borderRadius="md"
    mb={4}
    p={2}
  >
    <Heading size="sm" mb={2}>
      Add a User
    </Heading>
    <HStack wrap="wrap" spacing={2}>
      <Input
        size="sm"
        bg={inputBg}
        maxW="200px"
        type="email"
        isDisabled={mutation.isLoading}
        placeholder="Email address"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <Button
        colorScheme="brand"
        isLoading={mutation.isLoading}
        size="sm"
        rightIcon={<Icon as={FaArrowRight} />}
        onClick={async () => {
          await mutation.mutateAsync(email);
          setEmail('');
        }}
      >
        Add
      </Button>
    </HStack>
  </Box>;
}

export function GroupDirectory({ group }: GroupProps) {
  const [query, setQuery] = useState("");
  const mutation = useRemoveGroupUser(group.id);
  const canWrite = useHasPermission("groups", "write");
  const { manualUsers, lookupUsers } = useMemo(() => {
    let manualUsers: GroupUser[] = [];
    let lookupUsers: GroupUser[] = [];
    group.users
      ?.sort((a, b) => {
        if (a.userEmail < b.userEmail) {
          return -1;
        }
        if (a.userEmail > b.userEmail) {
          return 1;
        }
        return 0;
      })
      .forEach((u) => {
        if (u.isManual) {
          manualUsers.push(u);
        } else {
          lookupUsers.push(u);
        }
      });
    return { manualUsers, lookupUsers };
  }, [group]);
  return (
    <VStack align="stretch">
      {canWrite && <AddUserBox group={group} />}
      <InputGroup size="sm" maxW="500px" mb={2}>
        <InputLeftAddon>
          <Icon as={FaSearch} />
        </InputLeftAddon>
        <Input
          id="query"
          autoComplete="off"
          // maxW="300px"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search email addresses..."
        />
      </InputGroup>

      {manualUsers.length > 0 && (
        <>
          <Heading size="sm">Manual Users</Heading>
          <GroupUserList users={manualUsers} query={query}
          onDelete={u => mutation.mutate(u.userEmail)}/>
        </>
      )}
      {group.type !== "manual" && (
        <>
          <Heading size="sm">Lookup Users</Heading>
          <GroupUserList users={lookupUsers} query={query} />
        </>
      )}
    </VStack>
  );
}
