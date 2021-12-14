// import './App.css'

import {
  Box,
  Container,
  Flex,
  Heading,
  IconButton,
  useBreakpointValue,
  useColorModeValue,
} from "@chakra-ui/react";
import { useState } from "react";
import { FaBars } from "react-icons/fa";
import { Sidebar } from "./components/navigation/Sidebar";
import { Routes } from "./Routes";

const smVariant = { navigation: "drawer", navigationButton: true };
const mdVariant = { navigation: "sidebar", navigationButton: false };

function App() {
  const [isSidebarOpen, setSidebarOpen] = useState(false);
  const variants = useBreakpointValue({ base: smVariant, md: mdVariant });

  // TODO: refactoring with useCallback
  const toggleSidebar = () => setSidebarOpen(!isSidebarOpen);
  const bg = useColorModeValue("gray.50", "gray.800");

  return (
    <Flex height="100vh" bg={bg} pl={4}>
      <Sidebar
        variant={variants?.navigation as "drawer" | "sidebar"}
        isOpen={isSidebarOpen}
        onClose={toggleSidebar}
      />
      <Box flex="1" overflowY="auto" height="100%">
      <Container
        // mb={4}
        maxW="container.xl"
        py={6}
      >
        {variants?.navigationButton && (
          <Flex justifyContent="space-between" alignItems="center" mb={5}>
            <Heading size="lg" as="h1">
              KiFoMaSy
            </Heading>
            <IconButton
              size="sm"
              aria-label="open sidebar"
              onClick={toggleSidebar}
            >
              <FaBars />
            </IconButton>
          </Flex>
        )}
        <Routes />
      </Container>
      </Box>
    </Flex>
  );
}

export default App;
