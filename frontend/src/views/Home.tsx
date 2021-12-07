import {
  Heading,
  SimpleGrid,
  SkeletonCircle,
  SkeletonText,
} from "@chakra-ui/react";
import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useState } from "react";
import {
  FormalOverview,
  FormalProps,
} from "../components/display/FormalOverview";
import { Card, CardProps } from "../components/utility/Card";
import { Formal } from "../model/Formal";

function createFormal(data: Partial<Formal>): Formal {
  const template: Formal = {
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
    title: "Example Formal",
    guestLimit: 2,
    tickets: 100,
    ticketsRemaining: 50,
    saleStart: new Date("2021/12/25"),
    saleEnd: new Date("2022/01/01"),
  },
  {
    title: "Example Superformal",
    guestLimit: 1,
    saleEnd: new Date("2022/01/01"),
  },
  {
    title: "One More Formal",
  },
].map(createFormal);

const MotionCard = motion<CardProps>(Card);
const MotionOverview = motion<FormalProps>(FormalOverview);

export function Home() {
  const [loading, setLoading] = useState(true);
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
      <SimpleGrid
        templateColumns="repeat(auto-fill, minmax(300px, 1fr))"
        spacing="40px"
      >
        {/* TODO: fix weird visual behaviour */}
        <AnimatePresence exitBeforeEnter initial={false}>
          {loading ? (
            <MotionCard
              exit={{ scale: 0.5, opacity: 0 }}
              transition={{ ease: "easeIn" }}
            >
              <SkeletonCircle size="10" />
              <SkeletonText mt="4" noOfLines={4} spacing="4" />
            </MotionCard>
          ) : (
            formals.map((f, i) => (
              <MotionOverview
                // TODO: use actual DB ID as key
                key={`formal.${i}`}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{
                  delay: 0.1 * i
                }}
                formal={f}
              />
            ))
          )}
        </AnimatePresence>
      </SimpleGrid>
    </>
  );
}
