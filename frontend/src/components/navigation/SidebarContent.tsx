import { VStack } from "@chakra-ui/layout";
import { FaCog, FaHome, FaTicketAlt, FaUser } from "react-icons/fa";
import { SidebarItem } from "./SidebarItem";

const routes = [
  {
    to: "/",
    title: "Home",
    icon: FaHome,
    end: true,
  },
  {
    to: "/tickets",
    title: "Tickets",
    icon: FaTicketAlt,
  },
  {
    to: "/profile",
    title: "My Profile",
    icon: FaUser
  },
  {
    to: "/settings",
    title: "Settings",
    icon: FaCog
  }
];

interface SidebarContentProps {
  onClose?: () => void;
}

export function SidebarContent({onClose}: SidebarContentProps) {
  return (
    <VStack spacing="12px">
      {routes.map(({ to, title, icon, end }) => (
        <SidebarItem to={to} icon={icon} onClick={onClose} end={end} key={to}>
          {title}
        </SidebarItem>
      ))}
    </VStack>
  );
};
