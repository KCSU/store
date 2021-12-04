import {
  As,
  useColorModeValue,
  Button,
  Flex,
  Icon,
  Text,
} from "@chakra-ui/react";
import { To, useResolvedPath, useMatch, Link } from "react-router-dom";
import { IconBox } from "../utility/IconBox";

interface SidebarItemProps {
  icon: As<any>;
  to: To;
  end?: boolean;
}

export const SidebarItem: React.FC<SidebarItemProps> = ({
  children,
  icon,
  to,
  end,
}) => {
  const activeBg = useColorModeValue("white", "gray.700");
  const activeColor = useColorModeValue("gray.700", "white");
  const inactiveColor = useColorModeValue("gray.400", "gray.400");
  const purple = useColorModeValue("purple.400", "purple.600")
  const sidebarActiveShadow = "0px 7px 11px rgba(0, 0, 0, 0.04)";
  const trans = "0.2s linear";
  const resolved = useResolvedPath(to);
  const active = useMatch({ path: resolved.pathname, end: end ?? false });
  return (
    <Button
      as={Link}
      to={to}
      boxSize="initial"
      justifyContent="flex-start"
      alignItems="center"
      _hover={{}}
      boxShadow={active ? sidebarActiveShadow : "none"}
      bg={active ? activeBg : "transparent"}
      transition={trans}
      mx="auto"
      ps="16px"
      py="12px"
      borderRadius="15px"
      // _hover="none"
      w="100%"
      _active={{
        bg: "inherit",
        transform: "none",
        borderColor: "transparent",
      }}
      ringColor={purple}
      _focus={{}}
      _focusVisible={{
        boxShadow: active ? "0px 7px 11px rgba(0, 0, 0, 0.2)" : "none",
        ring: 3,
      }}
    >
      <Flex>
        <IconBox
          bg={purple}
          color="white"
          h="30px"
          w="30px"
          me="12px"
          transition={trans}
        >
          <Icon as={icon}></Icon>
        </IconBox>
        <Text
          color={active ? activeColor : inactiveColor}
          my="auto"
          fontSize="sm"
        >
          {children}
        </Text>
      </Flex>
    </Button>
  );
};
