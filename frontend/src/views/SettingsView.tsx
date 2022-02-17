import { Heading } from "@chakra-ui/layout";
import {
  Button,
  FormControl,
  FormLabel,
  Switch,
  useColorMode,
} from "@chakra-ui/react";
import { Card } from "../components/utility/Card";
import { useLogout } from "../hooks/useLogout";

export function SettingsView() {
  const { colorMode, toggleColorMode } = useColorMode();
  const mutation = useLogout();
  return (
    <Card alignItems="flex-start">
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
      <Button
        mt={4}
        colorScheme="brand"
        isLoading={mutation.isLoading}
        onClick={() => mutation.mutate()}
      >
        Logout
      </Button>
    </Card>
  );
}
