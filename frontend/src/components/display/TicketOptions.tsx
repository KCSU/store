import {
  useColorModeValue,
  FormControl,
  HStack,
  FormLabel,
  Box,
} from "@chakra-ui/react";
import { CreatableSelect } from "chakra-react-select";
import { Card } from "../utility/Card";

interface TicketOptionsProps {
  options: string[];
  value?: string;
  onChange?: (value: string) => void;
}

export const TicketOptions: React.FC<TicketOptionsProps> = ({
  options,
  value,
  onChange,
  children,
}) => {
  const handleChange = (option: Record<string, string>) => {
    onChange?.(option?.value ?? '')
  };
  const bg = useColorModeValue("white", "gray.700");
  return (
    <Card bg={bg} borderRadius={5} p={3}>
      {children}
      <FormControl as={HStack} spacing={4} alignItems="center">
        <FormLabel m={0}>Meal Option:</FormLabel>
        <Box flex="1">
          <CreatableSelect
            // TODO: fix long answers
            isClearable
            selectedOptionColor="purple"
            value={{
              label: value,
              value
            }}
            onChange={handleChange}
            size="sm"
            options={options.map((opt) => ({
              label: opt,
              value: opt,
            }))}
          ></CreatableSelect>
        </Box>
      </FormControl>
    </Card>
  );
};
