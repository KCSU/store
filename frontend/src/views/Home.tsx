import {
  Heading,
  SimpleGrid,
  SimpleGridProps,
  SkeletonCircle,
  SkeletonText,
} from "@chakra-ui/react";
import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useState } from "react";
import {
  FormalOverview,
  FormalProps,
} from "../components/display/FormalOverview";
import { generateMotion } from "../components/helpers/motion";
import { Card, CardProps } from "../components/utility/Card";
import { Formal } from "../model/Formal";

function createFormal(data: Partial<Formal>): Formal {
  const template: Formal = {
    id: 0,
    title: "",
    menu: "",
    price: 0,
    guestPrice: 0,
    options: [],
    guestLimit: 0,
    guestTickets: 0,
    guestTicketsRemaining: 0,
    tickets: 0,
    ticketsRemaining: 0,
    saleStart: new Date("2020/01/01"),
    saleEnd: new Date(),
  };
  return Object.assign(template, data);
}

const formals: Formal[] = [
  {
    id: 1,
    title: "Example Formal",
    guestLimit: 2,
    tickets: 100,
    ticketsRemaining: 50,
    saleStart: new Date("2021/12/25"),
    saleEnd: new Date("2022/01/01"),
  },
  {
    id: 2,
    title: "Example Superformal",
    guestLimit: 1,
    saleEnd: new Date("2022/01/01"),
  },
  {
    id: 3,
    title: "One More Formal",
  },
].map(createFormal);

const MotionCard = generateMotion<CardProps, 'div'>(Card);
const MotionOverview = motion<FormalProps>(FormalOverview);
const MotionSimpleGrid = motion<SimpleGridProps>(SimpleGrid);

export function Home() {
  const [loading, setLoading] = useState(true);
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
  useEffect(() => {
    let timer = setTimeout(() => {
      setLoading(false);
    }, 1000);
    return () => {
      clearTimeout(timer);
    };
  }, []);
  return (
    <>
      <Heading mb={5}>Upcoming Formals</Heading>
      <AnimatePresence exitBeforeEnter initial={false}>
        {loading ? (
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
            {formals.map((f, i) => (
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
