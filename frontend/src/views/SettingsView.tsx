import {
  Button,
  Heading,
  Link,
  FormControl,
  FormLabel,
  Switch,
  useColorMode,
  Icon,
  Flex,
  Text,
} from "@chakra-ui/react";
import { Card } from "../components/utility/Card";
import { useLogout } from "../hooks/mutations/useLogout";
import { Link as RouterLink } from "react-router-dom";
import { FaArrowRight, FaExternalLinkAlt, FaMoon, FaSun } from "react-icons/fa";
import { useContext } from "react";
import { UserContext } from "../model/User";

export function SettingsView() {
  const { colorMode, toggleColorMode } = useColorMode();
  const modeColor = colorMode === "light" ? "blue" : "yellow";
  const mutation = useLogout();
  const user = useContext(UserContext);
  return (
    <Card alignItems="flex-start">
      <Heading as="h3" size="lg" mb={4}>
        Settings
      </Heading>
      <Flex gap={3} mb={4} wrap="wrap">
        <Button
          colorScheme={modeColor}
          onClick={toggleColorMode}
          leftIcon={<Icon as={
            colorMode === "light" ? FaMoon : FaSun
          } />}
        >
          Use {colorMode === "light" ? "Dark Mode" : "Light Mode"}
        </Button>
        <Button
          as={RouterLink}
          to="/settings/access"
          variant="outline"
          rightIcon={<Icon as={FaArrowRight} />}
        >
          View Admin Access Log
        </Button>
      </Flex>
      <Heading as="h4" size="md" mb={2}>
        Logged in as:
      </Heading>
      <Text>
        {user?.name} {' '}
        (<Link href={`mailto:${user?.email}`} color="teal.500">
          {user?.email} <Icon as={FaExternalLinkAlt} boxSize={3} ml={0.5}/>
        </Link>)
      </Text>
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
