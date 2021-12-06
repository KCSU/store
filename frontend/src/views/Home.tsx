import {
  Heading,
  SimpleGrid,
  SkeletonCircle,
  SkeletonText,
} from "@chakra-ui/react";
import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useState } from "react";
import { Card, CardProps } from "../components/utility/Card";

const formals = [
  {
    title: "Example Formal",
  },
  // {
  //   title: "Example Superformal",
  // },
  // {
  //   title: "One More Formal",
  // },
];

const MotionCard = motion<CardProps>(Card);

export const Home: React.FC = () => {
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
        <AnimatePresence exitBeforeEnter>
          {loading ? (
            <MotionCard exit={{ scale: 0.5, opacity: 0 }} transition={{ease: 'easeIn'}}>
              <SkeletonCircle size="10" />
              <SkeletonText mt="4" noOfLines={4} spacing="4" />
            </MotionCard>
          ) : (
            formals.map((f, i) => (
              <MotionCard key={`formal.${i}`}
                initial={{ opacity: 0, translateY: '20px'}}
                animate={{ opacity: 1, translateY: 0 }}
              >
                <Heading size="md">{f.title}</Heading>
              </MotionCard>
            ))
          )}
        </AnimatePresence>
      </SimpleGrid>
    </>
  );
};
