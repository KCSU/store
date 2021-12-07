import {
  Badge,
  Button,
  Divider,
  Text,
  forwardRef,
  Heading,
  HStack,
  Stat,
  StatGroup,
  StatHelpText,
  StatLabel,
  StatNumber,
  Box,
  SimpleGrid,
  VStack,
  Flex,
  useColorModeValue,
  Center,
} from "@chakra-ui/react";
import { Formal } from "../../model/Formal";
import { Card } from "../utility/Card";

export interface FormalProps {
  formal: Formal;
}

function FormalStatusTag({ formal }: FormalProps) {
  if (formal.saleEnd < new Date()) {
    return <Badge>Closed</Badge>;
  } else if (formal.saleStart > new Date()) {
    return <Badge colorScheme="teal">Queue Now</Badge>;
  } else if (
    formal.guestTicketsRemaining === 0 &&
    formal.ticketsRemaining === 0
  ) {
    return <Badge colorScheme="red">Sold Out</Badge>;
  } else {
    return <Badge colorScheme="purple">Buy Now</Badge>;
  }
}

interface TicketStatsProps {
  prefix?: string;
  price: number;
  tickets: number;
  ticketsRemaining: number;
}

function TicketStats(props: TicketStatsProps) {
  return (
    <Box flex="1">
      <Stat>
        <StatLabel>
          {props.prefix}
          Ticket Price
        </StatLabel>
        <StatNumber>&#163;{props.price}</StatNumber>
      </Stat>
      <Stat>
        <StatLabel>
          {props.prefix}
          Tickets Left
        </StatLabel>
        <StatNumber>{props.ticketsRemaining}</StatNumber>
        <StatHelpText>out of {props.tickets}</StatHelpText>
      </Stat>
    </Box>
  );
}

function FormalStats({ formal }: FormalProps) {
  const prefix = formal.guestLimit > 0 ? "King's " : "";
  const noGuestBg = useColorModeValue("gray.100", "gray.600");
  const noGuestFg = useColorModeValue("gray.600", "gray.300");
  return (
    <Flex justifyContent="space-evenly">
      <TicketStats
        prefix={prefix}
        price={formal.price}
        tickets={formal.tickets}
        ticketsRemaining={formal.ticketsRemaining}
      />
      {formal.guestLimit > 0 ? (
        <>
        <Divider orientation="vertical" mx={2} />
        <TicketStats
          prefix="Guest "
          price={formal.guestPrice}
          tickets={formal.guestTickets}
          ticketsRemaining={formal.guestTicketsRemaining}
        />
        </>
      ) : (
        <Center flex="1" bg={noGuestBg} borderRadius={5}>
          <Text as="i" color={noGuestFg} textAlign="center">
            Guest tickets unavailable
          </Text>
        </Center>
      )}
    </Flex>
  );
}

export const FormalOverview = forwardRef<FormalProps, "div">(
  ({ formal }, ref) => {
    return (
      <Card ref={ref}>
        <HStack>
          <Heading size="md">{formal.title}</Heading>
          <FormalStatusTag formal={formal} />
        </HStack>
        <Divider my={2} />
        <FormalStats formal={formal} />
        {/* <Divider my={2} /> */}
        <HStack>
          <Button>More Info</Button>
        </HStack>
      </Card>
    );
  }
);
