import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Box,
  Text,
} from "@chakra-ui/react";
import { useRoles } from "../../hooks/admin/useRoles";

export function RolesList() {
  const { data, isLoading, isError } = useRoles();
  // TODO: loading states
  if (isLoading || isError || !data) {
    return <></>;
  }
  return (
    <Accordion defaultIndex={[]} allowMultiple>
      {data.map((role) => (
        <AccordionItem key={role.id}>
          <AccordionButton>
            <Box flex="1" textAlign="left">
              {role.name}
            </Box>
            <AccordionIcon />
          </AccordionButton>
          <AccordionPanel>
            {role.permissions?.map(p => (
              <Text key={p.id}>{p.resource}: {p.action}</Text>
            ))}
          </AccordionPanel>
        </AccordionItem>
      ))}
    </Accordion>
  );
}
