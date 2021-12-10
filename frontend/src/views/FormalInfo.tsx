import {
  Box,
  Button,
  Container,
  Heading,
  Icon,
  Text,
  VStack,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";
import { FaArrowLeft } from "react-icons/fa";
import { Link, Navigate, useParams } from "react-router-dom";
import { TicketBuyForm } from "../components/display/TicketBuyForm";
import { Card } from "../components/utility/Card";
import { formatMoney } from "../helpers/formatMoney";
import { getBuyText } from "../helpers/getBuyText";
import { useFormal } from "../hooks/useFormal";

interface TicketStatsProps {
  prefix?: string;
  price: number;
  tickets: number;
  ticketsRemaining: number;
}

export const TicketStats: React.FC<TicketStatsProps> = (props) => {
  return (
    <WrapItem
      display="block"
      flex="1"
      flexBasis="250px"
      borderWidth="1px"
      borderRadius="md"
      p={3}
    >
      <Heading as="h6" size="sm">
        {props.prefix} Ticket Price:
      </Heading>
      <Text mb={3}>{formatMoney(props.price)}</Text>
      <Heading as="h6" size="sm">
        {props.prefix} Tickets Remaining:
      </Heading>
      <Text>
        {props.ticketsRemaining} out of {props.tickets}
        &nbsp;
        {props.children}
      </Text>
    </WrapItem>
  );
};

export function FormalInfo() {
  const { formalId } = useParams();
  const formal = useFormal(parseInt(formalId ?? "0"));
  if (!formal) {
    // TODO: return an error!
    return <Navigate to="/"/>;
  }
  const prefix = formal.guestLimit > 0 ? "King's " : "";
  return (
    // TODO: guest list, responsive meal option
    <Container maxW="container.md" p={0}>
      <Button
        as={Link}
        // size="sm"
        to="/"
        variant="ghost"
        mb={4}
        leftIcon={<Icon as={FaArrowLeft} />}
      >
        Back Home
      </Button>
      <Card mb={5}>
        <Heading as="h3" size="lg" mb={5}>
          {formal.title}
        </Heading>
        <VStack alignItems="stretch">
          <Wrap justifyContent="space-between">
            <TicketStats
              price={formal.price}
              tickets={formal.tickets}
              ticketsRemaining={formal.ticketsRemaining}
              prefix={prefix}
            ></TicketStats>
            {formal.guestLimit > 0 ? (
              <TicketStats
                price={formal.guestPrice}
                tickets={formal.guestTickets}
                ticketsRemaining={formal.guestTicketsRemaining}
                prefix="Guest "
              >
                <br />
                <Text as="i">
                  (up to {formal.guestLimit} per King's member)
                </Text>
              </TicketStats>
            ) : (
              <WrapItem
                display="block"
                flex="1"
                flexBasis="250px"
                borderWidth="1px"
                borderRadius="md"
                p={3}
              >
                <Text as="i">Guest tickets not available.</Text>
              </WrapItem>
            )}
          </Wrap>
          <Box borderWidth="1px" borderRadius="md" p={3}>
            <Heading as="h5" size="sm">
              Menu
            </Heading>
            <Text>
              {/* TODO: rich text */}
              {formal.menu}
            </Text>
          </Box>
          <VStack align="stretch" borderWidth="1px" borderRadius="md" p={3}>
            <TicketBuyForm formal={formal} hasShadow={false} />
            <Button colorScheme="purple">{getBuyText(formal)}</Button>
          </VStack>
        </VStack>
      </Card>
    </Container>
  );
}
