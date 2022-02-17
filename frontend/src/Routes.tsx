import { Box, BoxProps } from "@chakra-ui/react";
import { AnimatePresence } from "framer-motion";
import { Route, Routes as ReactRoutes, useLocation } from "react-router";
import { generateMotion } from "./components/utility/generateMotion";
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

const MotionBox = generateMotion<BoxProps, 'div'>(Box);

export function Routes() {
  const adminRoutes = [
    {
      path: "/admin/formals",
      element: <AdminFormalListView/>,
      resource: "formals",
      action: "read"
    },
    {
      path: "/admin/formals/:id",
      element: <AdminEditFormalView/>,
      resource: "formals",
      action: "read"
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
                <Route path={path} element={
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
