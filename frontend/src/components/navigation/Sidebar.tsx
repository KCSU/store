import { Box, Drawer, DrawerBody, DrawerCloseButton, DrawerContent, DrawerHeader, DrawerOverlay } from "@chakra-ui/react"

interface SidebarProps {
    onClose: () => void
    isOpen: boolean
    variant: 'drawer' | 'sidebar'
}

const SidebarContent = () => {
    return <div>Sidebar Content</div>
}

// TODO: Close button & Hamburger button
export const Sidebar = ({ isOpen, variant, onClose }: SidebarProps) => {
    return variant === 'sidebar' ? (
      <Box
        left={0}
        p={5}
        w="300px"
        top={0}
        h="100%"
        bg="#dfdfdf"
      >
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
    )
  }