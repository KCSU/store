import {
  As,
  useColorModeValue,
  Button,
  Flex,
  Icon,
  Text,
} from "@chakra-ui/react";
import { To, useResolvedPath, useMatch, Link } from "react-router-dom";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { IconBox } from "../utility/IconBox";

interface SidebarItemProps {
  icon: As<any>;
  to: To;
  onClick?: () => void;
  end?: boolean;
}

export const SidebarItem: React.FC<SidebarItemProps> = ({
  children,
  icon,
  to,
  onClick,
  end,
}) => {
  const activeBg = useColorModeValue("white", "gray.700");
  const activeColor = useColorModeValue("gray.700", "white");
  const inactiveColor = useColorModeValue("gray.400", "gray.400");
  const sidebarActiveShadow = "0px 7px 11px rgba(0, 0, 0, 0.04)";
  const trans = "0.2s linear";
  const resolved = useResolvedPath(to);
  const active = useMatch({ path: resolved.pathname, end: end ?? false });
  return (
    <Button
      as={Link}
      to={to}
      onClick={onClick}
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
      ringColor="brand.500"
      _focus={{}}
      _focusVisible={{
        boxShadow: active ? "0px 7px 11px rgba(0, 0, 0, 0.2)" : "none",
        ring: 3,
      }}
    >
      <Flex>
        <IconBox
          bg={active ? "brand.500" : activeBg}
          color={active ? "white" : "brand.400"}
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

type AdminSidebarItemProps = SidebarItemProps & {
  resource: string;
  action: string;
}

export const AdminSidebarItem: React.FC<AdminSidebarItemProps> = ({resource, action, ...props}) => {
  const hasAccess = useHasPermission(resource, action);
  // return <SidebarItem {...props}/>;
  return hasAccess ? <SidebarItem {...props}/> : null;
}