import {
  Box,
  Button,
  Checkbox,
  CheckboxGroup,
  HStack,
  Icon,
  VStack,
} from "@chakra-ui/react";
import { useState } from "react";
import { FaSave } from "react-icons/fa";
import { useEditFormalGroups } from "../../hooks/admin/useEditFormalGroups";
import { useGroups } from "../../hooks/admin/useGroups";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { Formal } from "../../model/Formal";

interface FormalProps {
  formal: Formal;
}

export function EditFormalGroupsForm({ formal }: FormalProps) {
  const { data: available, isError, isLoading } = useGroups();
  const mutation = useEditFormalGroups(formal.id);
  const canWrite = useHasPermission("formals", "write");
  const [groups, setGroups] = useState(
    formal.groups?.map((g) => g.id.toString()) ?? []
  );

  // TODO: handle with loading states
  if (isError || isLoading || !available) {
    return <Box></Box>;
  }
  return (
    <VStack align="start" gap={4}>
      <HStack gap={10} flexWrap="wrap">
        <CheckboxGroup
          isDisabled={!canWrite}
          colorScheme="brand"
          value={groups}
          onChange={(v) => setGroups(v as string[])}
        >
          {available.map((g) => (
            <Checkbox key={g.id} value={g.id.toString()}>
              {g.name}
            </Checkbox>
          ))}
        </CheckboxGroup>
      </HStack>
      {canWrite && (
        <Button
          isLoading={mutation.isLoading}
          onClick={async () => {
            await mutation.mutateAsync(groups.map((g) => parseInt(g)));
          }}
          colorScheme="brand"
          leftIcon={<Icon as={FaSave} />}
        >
          Save Changes
        </Button>
      )}
    </VStack>
  );
}
