import {
  Box,
  Drawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerHeader,
  DrawerOverlay,
} from "@chakra-ui/react";
import React from "react";
import { SidebarContent } from "./SidebarContent";

interface SidebarProps {
  onClose: () => void;
  isOpen: boolean;
  variant: "drawer" | "sidebar";
}

// TODO: Close button & Hamburger button
export const Sidebar: React.FC<SidebarProps> = ({
  isOpen,
  variant,
  onClose,
}) => {
  return variant === "sidebar" ? (
    <Box left={0} p={5} w="300px" top={0} h="100%">
      <SidebarContent /*onClick={onClose}*/ />
    </Box>
  ) : (
    <Drawer isOpen={isOpen} placement="left" onClose={onClose}>
      <DrawerOverlay>
        <DrawerContent>
          <DrawerCloseButton />
          <DrawerHeader>Chakra-UI</DrawerHeader>
          <DrawerBody>
            <SidebarContent /*onClick={onClose}*/ />
          </DrawerBody>
        </DrawerContent>
      </DrawerOverlay>
    </Drawer>
  );
};
