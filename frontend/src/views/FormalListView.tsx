import {
  Box,
  Heading,
  SimpleGrid,
  SimpleGridProps,
  SkeletonCircle,
  SkeletonText,
} from "@chakra-ui/react";
import { AnimatePresence, motion } from "framer-motion";
import {
  FormalInfoCard,
  FormalProps,
} from "../components/formals/FormalInfoCard";
import { motionComponent } from "../components/utility/generateMotion";
import { Card, CardProps } from "../components/utility/Card";
import { useFormals } from "../hooks/queries/useFormals";

const MotionCard = motionComponent<CardProps, "div">(Card);
const MotionFormalInfoCard = motion<FormalProps>(FormalInfoCard);
const MotionSimpleGrid = motion<SimpleGridProps>(SimpleGrid);

function LoadingState() {
  return (
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
  );
}

function ErrorState() {
  // TODO: implement this
  return <Box></Box>
}

function FormalGrid() {
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
  if (isError) {
    return <ErrorState />;
  }
  if (isLoading && !formals) {
    return <LoadingState />;
  }
  return (
    <MotionSimpleGrid
      key="loadedGrid"
      variants={gridVariant}
      initial="hidden"
      animate="show"
      templateColumns="repeat(auto-fill, minmax(300px, 1fr))"
      spacing="40px"
    >
      {formals?.map((f, i) => (
        <MotionFormalInfoCard layout="position"
          // TODO: use actual DB ID as key
          key={`formal.${f.id}`}
          variants={itemVariant}
          formal={f}
        />
      ))}
    </MotionSimpleGrid>
  );
}


export function FormalListView() {
  return (
    <>
      <Heading size="xl" mb={5}>
        Upcoming Formals
      </Heading>
      <AnimatePresence exitBeforeEnter initial={false}>
        {FormalGrid()}
      </AnimatePresence>
    </>
  );
}
