import { Heading } from "@chakra-ui/layout";
import { FormControl, FormLabel, Switch, useColorMode } from "@chakra-ui/react";
import { Card } from "../components/utility/Card";

export function Settings() {
  const { colorMode, toggleColorMode } = useColorMode();
  return (
      <Card>
        <Heading as="h3" size="lg" mb={5}>
          Settings
        </Heading>
        <FormControl display="flex" alignItems="center">
          <FormLabel htmlFor="dark-mode" mb="0">
            Dark Mode:
          </FormLabel>
          <Switch
            id="dark-mode"
            colorScheme="brand"
            isChecked={colorMode === "dark"}
            onChange={toggleColorMode}
          />
        </FormControl>
      </Card>
  );
};
