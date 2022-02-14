import { Heading, VStack } from "@chakra-ui/layout";
import {
  FaCalendarDay,
  FaCog,
  FaHome,
  FaTicketAlt,
  FaUser,
} from "react-icons/fa";
import { AdminSidebarItem, SidebarItem } from "./SidebarItem";

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
  // {
  //   to: "/profile",
  //   title: "My Profile",
  //   icon: FaUser
  // },
  {
    to: "/settings",
    title: "Settings",
    icon: FaCog,
  },
];

const adminRoutes = [
  {
    to: "/admin/formals",
    title: "Manage Formals",
    resource: "formals",
    action: "read",
    icon: FaCalendarDay,
  },
];

interface SidebarContentProps {
  onClose?: () => void;
}

export function SidebarContent({ onClose }: SidebarContentProps) {
  const adminItems = adminRoutes.map(
    ({ to, title, icon, resource, action }) => (
      <AdminSidebarItem
        key={to}
        to={to}
        icon={icon}
        onClick={onClose}
        // end={false}
        resource={resource}
        action={action}
      >
        {title}
      </AdminSidebarItem>
    )
  );
  const showAdmin = adminItems.some((item) => item !== null);
  return (
    <VStack spacing="12px">
      {routes.map(({ to, title, icon, end }) => (
        <SidebarItem to={to} icon={icon} onClick={onClose} end={end} key={to}>
          {title}
        </SidebarItem>
      ))}
      {showAdmin && (
        <Heading as="h3" size="sm" alignSelf="flex-start" pt={4}>
          Admin
        </Heading>
      )}
      {adminItems}
    </VStack>
  );
}
