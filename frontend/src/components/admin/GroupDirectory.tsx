import {
  Heading,
  Icon,
  Input,
  InputGroup,
  InputLeftAddon,
  InputRightAddon,
  VStack,
  Box,
  HStack,
  Button,
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

function AddUserBox({ group }: GroupProps) {
  const [crsid, setCrsid] = useState("");
  const mutation = useAddGroupUser(group.id);
  const cleanCrsid = crsid.trim().toLowerCase();
  const isValid = /^[a-z0-9]+$/.test(cleanCrsid);
  return (
    <Box
      as="form"
      mb={2}
      onSubmit={(e: React.FormEvent) => {
        e.preventDefault();
        if (isValid) {
          mutation.mutate(`${cleanCrsid}@cam.ac.uk`, {
            onSuccess: () => setCrsid(""),
          });
        }
      }}
    >
      <Heading size="sm" mb={2}>
        Add a Manual User
      </Heading>
      <HStack maxW="500px" spacing={2}>
        <InputGroup size="sm">
          <Input
            placeholder="CRSid"
            autoComplete="off"
            isDisabled={mutation.isLoading}
            value={crsid}
            onChange={(e) => setCrsid(e.target.value)}
          />
          <InputRightAddon>@cam.ac.uk</InputRightAddon>
        </InputGroup>
        <Button
          type="submit"
          colorScheme="brand"
          size="sm"
          flexShrink={0}
          isLoading={mutation.isLoading}
          isDisabled={!isValid}
          rightIcon={<Icon as={FaArrowRight} />}
        >
          Add
        </Button>
      </HStack>
    </Box>
  );
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
