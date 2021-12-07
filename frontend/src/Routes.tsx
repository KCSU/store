import { Box } from "@chakra-ui/layout";
import { AnimatePresence, motion } from "framer-motion";
import { Route, Routes as ReactRoutes, useLocation } from "react-router";
import { Home } from "./views/Home";
import { Settings } from "./views/Settings";

const MotionBox = motion(Box);

export function Routes() {
  const location = useLocation();
  return (
    <AnimatePresence exitBeforeEnter initial={false}>
      <MotionBox
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
          <Route path="/" element={<Home />} />
          <Route path="/settings" element={<Settings />} />
        </ReactRoutes>
      </MotionBox>
    </AnimatePresence>
  );
};
