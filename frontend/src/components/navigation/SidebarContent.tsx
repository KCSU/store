import { FaHome, FaTicketAlt } from "react-icons/fa";
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
];

export const SidebarContent = () => {
  return (
    <>
      {routes.map(({ to, title, icon, end }) => (
        <SidebarItem to={to} icon={icon} end={end}>
          {title}
        </SidebarItem>
      ))}
    </>
  );
};
