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

export const SidebarContent = () => {
  return (
    <VStack spacing="12px">
      {routes.map(({ to, title, icon, end }) => (
        <SidebarItem to={to} icon={icon} end={end} key={to}>
          {title}
        </SidebarItem>
      ))}
    </VStack>
  );
};
