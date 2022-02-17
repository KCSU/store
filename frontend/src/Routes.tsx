import { Box, BoxProps } from "@chakra-ui/react";
import { AnimatePresence } from "framer-motion";
import { Route, Routes as ReactRoutes, useLocation } from "react-router";
import { generateMotion } from "./components/utility/generateMotion";
import { RequireAdmin } from "./components/utility/RequireAdmin";
import { RequireAuth } from "./components/utility/RequireAuth";
import { EditFormal } from "./views/admin/EditFormal";
import { FormalList } from "./views/admin/FormalList";
import { EditTickets } from "./views/EditTickets";
import { FormalInfo } from "./views/FormalInfo";
import { Home } from "./views/Home";
import { Login } from "./views/Login";
import { Settings } from "./views/Settings";
import { TicketsView } from "./views/TicketsView";

const MotionBox = generateMotion<BoxProps, 'div'>(Box);

export function Routes() {
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
          <Route path="/login" element={<Login/>} />
          <Route element={<RequireAuth/>}>
            <Route path="/" element={<Home />} />
            <Route path="/settings" element={<Settings />} />
            <Route path="/formals/:formalId" element={<FormalInfo />} />
            <Route path="/tickets" element={<TicketsView />} />
            <Route path="/tickets/:id" element={<EditTickets />}/>
            {/* ADMIN ROUTES */}
            <Route path="/admin/formals" element={
              <RequireAdmin resource="formals" action="read">
                <FormalList />
              </RequireAdmin>
            }/>
            <Route path="/admin/formals/:id" element={
              <RequireAdmin resource="formals" action="read">
                <EditFormal />
              </RequireAdmin>
            }/>
          </Route>
        </ReactRoutes>
      </MotionBox>
    </AnimatePresence>
  );
};
