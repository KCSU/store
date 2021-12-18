import {
  Heading,
  SimpleGrid,
  SimpleGridProps,
  SkeletonCircle,
  SkeletonText,
} from "@chakra-ui/react";
import { AnimatePresence, motion } from "framer-motion";
import {
  FormalOverview,
  FormalProps,
} from "../components/display/FormalOverview";
import { generateMotion } from "../components/utility/generateMotion";
import { Card, CardProps } from "../components/utility/Card";
import { useFormals } from "../hooks/useFormals";

const MotionCard = generateMotion<CardProps, "div">(Card);
const MotionOverview = motion<FormalProps>(FormalOverview);
const MotionSimpleGrid = motion<SimpleGridProps>(SimpleGrid);

export function Home() {
  const { data: formals, isLoading, isError } = useFormals();
  const gridVariant = {
    hidden: {},
    show: {
      transition: {
        staggerChildren: 0.1,
      },
    },
  };
  const itemVariant = {
    hidden: {
      opacity: 0,
      y: 20,
    },
    show: {
      opacity: 1,
      y: 0,
    },
  };
  return (
    <>
      <Heading size="xl" mb={5}>
        Upcoming Formals
      </Heading>
      <AnimatePresence exitBeforeEnter initial={false}>
        {isLoading ? (
          <SimpleGrid
            key="loadingGrid"
            templateColumns="repeat(auto-fill, minmax(300px, 1fr))"
            spacing="40px"
          >
            <MotionCard
              exit={{ scale: 0.5, opacity: 0 }}
              transition={{ duration: 0.25 }}
            >
              <SkeletonCircle size="10" />
              <SkeletonText mt="4" noOfLines={4} spacing="4" />
            </MotionCard>
          </SimpleGrid>
        ) : (
          <MotionSimpleGrid
            key="loadedGrid"
            variants={gridVariant}
            initial="hidden"
            animate="show"
            templateColumns="repeat(auto-fill, minmax(300px, 1fr))"
            spacing="40px"
          >
            {formals?.map((f, i) => (
              <MotionOverview
                // TODO: use actual DB ID as key
                key={`formal.${i}`}
                variants={itemVariant}
                formal={f}
              />
            ))}
          </MotionSimpleGrid>
        )}
      </AnimatePresence>
    </>
  );
}
