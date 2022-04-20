import {
  useColorModeValue,
  FormControl,
  HStack,
  FormLabel,
  Box,
} from "@chakra-ui/react";
import { isDisabled } from "@chakra-ui/utils";
import { CreatableSelect } from "chakra-react-select";
import { Card } from "../utility/Card";
import { TicketOptionsSelect } from "./TicketOptionsSelect";

interface TicketOptionsInputProps {
  value?: string;
  onChange?: (value: string) => void;
  hasShadow?: boolean;
  isDisabled?: boolean;
  footer?: React.ReactNode;
}

export const TicketOptionsInput: React.FC<TicketOptionsInputProps> = ({
  value,
  onChange,
  children,
  footer,
  isDisabled,
  hasShadow = true,
}) => {
  const bg = useColorModeValue("white", "gray.700");
  return (
    <Card
      bg={bg}
      borderRadius={5}
      p={3}
      {...(!hasShadow && {
        boxShadow: "none",
      })}
    >
      {children}
      <FormControl isDisabled={isDisabled} as={HStack} spacing={4} alignItems="center">
        <FormLabel m={0}>Meal Option:</FormLabel>
        <Box flex="1">
          <TicketOptionsSelect size='sm' value={value} onChange={onChange} />
        </Box>
      </FormControl>
      {footer}
    </Card>
  );
};
