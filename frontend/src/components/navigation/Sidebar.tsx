import {
  Box,
  Drawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerHeader,
  DrawerOverlay,
  Heading,
  useColorModeValue,
} from "@chakra-ui/react";
import { useLocation } from "react-router-dom";
import { SidebarContent } from "./SidebarContent";

interface SidebarProps {
  onClose: () => void;
  isOpen: boolean;
  variant: "drawer" | "sidebar";
}

// TODO: Close button & Hamburger button
export function Sidebar({ isOpen, variant, onClose }: SidebarProps) {
  const bg = useColorModeValue("gray.50", "gray.800");
  return variant === "sidebar" ? (
    <Box left={0} p={6} w="275px" top={0} h="100%">
      <Heading mb="12px" size="xl" as="h1">
        KiFoMaSy
      </Heading>
      <SidebarContent /*onClick={onClose}*/ />
    </Box>
  ) : (
    <Drawer isOpen={isOpen} placement="left" onClose={onClose}>
      <DrawerOverlay>
        <DrawerContent bg={bg}>
          <DrawerCloseButton />
          <DrawerHeader>KiFoMaSy</DrawerHeader>
          <DrawerBody>
            <SidebarContent /*onClick={onClose}*/ />
          </DrawerBody>
        </DrawerContent>
      </DrawerOverlay>
    </Drawer>
  );
}
