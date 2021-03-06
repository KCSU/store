import { Heading, VStack } from "@chakra-ui/layout";
import { useContext } from "react";
import {
  FaCalendarDay,
  FaCog,
  FaHome,
  FaReceipt,
  FaShieldAlt,
  FaTicketAlt,
  FaUsers,
} from "react-icons/fa";
import { UserContext } from "../../model/User";
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
  {
    to: "/admin/groups",
    title: "Manage Groups",
    resource: "groups",
    action: "read", // write?
    icon: FaUsers,
  },
  {
    to: "/admin/roles",
    title: "Permissions",
    resource: "roles",
    action: "read",
    icon: FaShieldAlt,
  },
  {
    to: "/admin/bills",
    title: "Billing",
    resource: "billing",
    action: "read",
    icon: FaReceipt,
  }
];

interface SidebarContentProps {
  onClose?: () => void;
}

export function SidebarContent({ onClose }: SidebarContentProps) {
  const user = useContext(UserContext);
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
  const showAdmin = (user?.permissions.length ?? 0) > 0;
  return (
    <VStack spacing="12px">
      {routes.map(({ to, title, icon, end }) => (
        <SidebarItem to={to} icon={icon} onClick={onClose} end={end} key={to}>
          {title}
        </SidebarItem>
      ))}
      {showAdmin && (
        <>
          <Heading as="h3" size="sm" alignSelf="flex-start" pt={4}>
            Admin
          </Heading>
          {adminItems}
        </>
      )}
    </VStack>
  );
}
