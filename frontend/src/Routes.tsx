import { Box, BoxProps } from "@chakra-ui/react";
import { AnimatePresence } from "framer-motion";
import { Route, Routes as ReactRoutes, useLocation } from "react-router";
import { motionComponent } from "./components/utility/generateMotion";
import { RequireAdmin } from "./components/utility/RequireAdmin";
import { RequireAuth } from "./components/utility/RequireAuth";
import { AdminEditFormalView } from "./views/admin/AdminEditFormalView";
import { AdminFormalListView } from "./views/admin/AdminFormalListView";
import { EditFormalTicketsView } from "./views/EditFormalTicketsView";
import { FormalView } from "./views/FormalView";
import { FormalListView } from "./views/FormalListView";
import { LoginView } from "./views/LoginView";
import { SettingsView } from "./views/SettingsView";
import { TicketsView } from "./views/TicketsView";
import { AdminCreateFormalView } from "./views/admin/AdminCreateFormalView";
import { AdminGroupListView } from "./views/admin/AdminGroupListView";
import { AdminEditGroupView } from "./views/admin/AdminEditGroupView";
import { AdminCreateGroupView } from "./views/admin/AdminCreateGroupView";
import { AdminPermissionsView } from "./views/admin/AdminPermissionsView";
import { AdminEditRoleView } from "./views/admin/AdminEditRoleView";
import { AdminBillsListView } from "./views/admin/AdminBillsListView";
import { AdminEditBillView } from "./views/admin/AdminEditBillView";
import { AdminCreateBillView } from "./views/admin/AdminCreateBillView";

const MotionBox = motionComponent<BoxProps, 'div'>(Box);

export function Routes() {
  const adminRoutes = [
    {
      path: "/admin/formals",
      element: <AdminFormalListView/>,
      resource: "formals",
      action: "read"
    },
    {
      path: "/admin/formals/create",
      element: <AdminCreateFormalView/>,
      resource: "formals",
      action: "write"
    },
    {
      path: "/admin/formals/:formalId",
      element: <AdminEditFormalView/>,
      resource: "formals",
      action: "read"
    },
    {
      path: "/admin/groups",
      element: <AdminGroupListView/>,
      resource: "groups",
      action: "read" // Should this be "write"?
    },
    {
      path: "/admin/groups/create",
      element: <AdminCreateGroupView/>,
      resource: "groups",
      action: "write"
    },
    {
      path: "/admin/groups/:groupId",
      element: <AdminEditGroupView/>,
      resource: "groups",
      action: "read"
    },
    {
      path: "/admin/roles",
      element: <AdminPermissionsView/>,
      resource: "roles",
      action: "read"
    },
    {
      path: "/admin/roles/:id",
      element: <AdminEditRoleView/>,
      resource: "roles",
      action: "read"
    },
    {
      path: "/admin/bills",
      element: <AdminBillsListView/>,
      resource: "billing",
      action: "read"
    },
    {
      path: "/admin/bills/:id",
      element: <AdminEditBillView/>,
      resource: "billing",
      action: "read"
    },
    {
      path: "/admin/bills/create",
      element: <AdminCreateBillView/>,
      resource: "billing",
      action: "write"
    }
  ];
  const location = useLocation();
  return (
    <AnimatePresence exitBeforeEnter initial={false}>
      <MotionBox mr={{md: 4}}
        key={location.pathname}
        initial={{
          opacity: 0,
          y: 10
        }}
        animate={{
          opacity: 1,
          y: 0
        }}
        exit={{
          opacity: 0,
          y: 10
        }}
        transition={{
          duration: 0.15,
          ease: 'linear'
        }}
      >
        <ReactRoutes location={location}>
          <Route path="/login" element={<LoginView/>} />
          <Route element={<RequireAuth/>}>
            <Route path="/" element={<FormalListView />} />
            <Route path="/settings" element={<SettingsView />} />
            <Route path="/formals/:formalId" element={<FormalView />} />
            <Route path="/tickets" element={<TicketsView />} />
            <Route path="/tickets/:id" element={<EditFormalTicketsView />}/>
            {
              adminRoutes.map(({path, element, resource, action}) => (
                <Route key={path} path={path} element={
                  <RequireAdmin resource={resource} action={action}>
                    {element}
                  </RequireAdmin>
                }/>
              ))
            }
          </Route>
        </ReactRoutes>
      </MotionBox>
    </AnimatePresence>
  );
};
