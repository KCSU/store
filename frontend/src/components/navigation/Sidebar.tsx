import {
  Box,
  BoxProps,
  Drawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerHeader,
  DrawerOverlay,
  Heading,
  useColorModeValue,
} from "@chakra-ui/react";
import { motionComponent } from "../utility/generateMotion";
import { SidebarContent } from "./SidebarContent";

interface SidebarProps {
  onClose: () => void;
  isOpen: boolean;
  variant: "drawer" | "sidebar";
}

const MotionBox = motionComponent<BoxProps, "div">(Box);

// TODO: Close button & Hamburger button
// TODO: move into routes for better visibility control
export function Sidebar({ isOpen, variant, onClose }: SidebarProps) {
  const bg = useColorModeValue("gray.50", "gray.800");
  return variant === "sidebar" ? (
    <MotionBox
      left={0}
      p={6}
      w="275px"
      top={0}
      h="100%"
      animate={{
        x: 0
      }}
      exit={{
        x: "-100%"
      }}
      transition={{
        duration: 0.15,
        // ease: 'linear'
      }}
    >
      <Heading mb="12px" size="xl" as="h1">
        KiFoMaSy
      </Heading>
      <SidebarContent onClose={onClose} />
    </MotionBox>
  ) : (
    <Drawer isOpen={isOpen} placement="left" onClose={onClose}>
      <DrawerOverlay>
        <DrawerContent bg={bg}>
          <DrawerCloseButton />
          <DrawerHeader>KiFoMaSy</DrawerHeader>
          <DrawerBody>
            <SidebarContent onClose={onClose} />
          </DrawerBody>
        </DrawerContent>
      </DrawerOverlay>
    </Drawer>
  );
}
