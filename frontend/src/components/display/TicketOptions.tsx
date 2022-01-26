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
  value?: string;
  onChange?: (value: string) => void;
  hasShadow?: boolean;
}

export const TicketOptions: React.FC<TicketOptionsProps> = ({
  value,
  onChange,
  children,
  hasShadow = true,
}) => {
  const options = ["Normal", "Vegetarian", "Vegan", "Pescetarian"];
  const handleChange = (option: Record<string, string>) => {
    onChange?.(option?.value ?? "");
  };
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
      <FormControl as={HStack} spacing={4} alignItems="center">
        <FormLabel m={0}>Meal Option:</FormLabel>
        <Box flex="1">
          <CreatableSelect
            // TODO: fix long answers
            isClearable
            selectedOptionColor="brand"
            value={{
              label: value,
              value,
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
