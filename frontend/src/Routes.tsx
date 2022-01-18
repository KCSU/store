import { Box, BoxProps } from "@chakra-ui/react";
import { AnimatePresence, motion } from "framer-motion";
import { Route, Routes as ReactRoutes, useLocation } from "react-router";
import { generateMotion } from "./components/utility/generateMotion";
import { RequireAuth } from "./components/utility/RequireAuth";
import { FormalInfo } from "./views/FormalInfo";
import { Home } from "./views/Home";
import { Login } from "./views/Login";
import { Settings } from "./views/Settings";
import { Tickets } from "./views/Tickets";

const MotionBox = generateMotion<BoxProps, 'div'>(Box);

export function Routes() {
  const location = useLocation();
  return (
    <AnimatePresence exitBeforeEnter initial={false}>
      <MotionBox mr={4}
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
            <Route path="/tickets" element={<Tickets />} />
          </Route>
        </ReactRoutes>
      </MotionBox>
    </AnimatePresence>
  );
};
