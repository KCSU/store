import {
  Box,
  Button,
  Checkbox,
  CheckboxGroup,
  HStack,
  Icon,
  VStack,
} from "@chakra-ui/react";
import { useContext, useState } from "react";
import { FaSave } from "react-icons/fa";
import { useEditFormalGroups } from "../../hooks/admin/useEditFormalGroups";
import { useGroups } from "../../hooks/admin/useGroups";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { Formal, FormalContext } from "../../model/Formal";

export function EditFormalGroupsForm() {
  const formal = useContext(FormalContext);
  const { data: available, isError, isLoading } = useGroups();
  const mutation = useEditFormalGroups(formal.id);
  const canWrite = useHasPermission("formals", "write");
  const [groups, setGroups] = useState(
    formal.groups?.map((g) => g.id) ?? []
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
            <Checkbox key={g.id} value={g.id}>
              {g.name}
            </Checkbox>
          ))}
        </CheckboxGroup>
      </HStack>
      {canWrite && (
        <Button
          isLoading={mutation.isLoading}
          onClick={async () => {
            await mutation.mutateAsync(groups);
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
