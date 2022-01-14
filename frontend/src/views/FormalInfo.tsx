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
import { Link, Navigate, useNavigate, useParams } from "react-router-dom";
import { TicketBuyForm } from "../components/display/TicketBuyForm";
import { Card } from "../components/utility/Card";
import { formatMoney } from "../helpers/formatMoney";
import { getBuyText } from "../helpers/getBuyText";
import { useBuyTicket } from "../hooks/useBuyTicket";
import { useDateTime } from "../hooks/useDateTime";
import { useFormals } from "../hooks/useFormals";
import { useQueueRequest } from "../hooks/useQueueRequest";
import { Formal } from "../model/Formal";

interface TicketStatsProps {
  prefix?: string;
  price: number;
  tickets: number;
  ticketsRemaining: number;
}

const TicketStats: React.FC<TicketStatsProps> = (props) => {
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

interface FormalInfoViewProps {
  formal: Formal;
}

function FormalInfoView({formal}: FormalInfoViewProps) {
  // Formal Data
  const datetime = useDateTime(formal.dateTime);
  const prefix = formal.guestLimit > 0 ? "King's " : "";
  const mutation = useBuyTicket();
  const navigate = useNavigate();
  // State management
  const [queueRequest, dispatchQR] = useQueueRequest(formal.id);

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
        <Heading as="h3" size="lg" mb={1}>
          {formal.name}
        </Heading>
        <Text as="b" mb={4}>
          {datetime}
        </Text>
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
            <TicketBuyForm
              formal={formal}
              hasShadow={false}
              value={queueRequest}
              onChange={dispatchQR}
            />
            <Button
              colorScheme="brand"
              onClick={async () => {
                await mutation.mutateAsync(queueRequest);
                navigate("/tickets");
              }}
              isLoading={mutation.isLoading}
            >
              {getBuyText(formal)}
            </Button>
          </VStack>
        </VStack>
      </Card>
    </Container>
  );
}

// TODO: Date and time!
export function FormalInfo() {
  // Get the formal
  const { formalId } = useParams();
  const formalIdNum = parseInt(formalId ?? "0");
  const { data: formals, isLoading, isError } = useFormals();
  const formal = formals?.find((f) => f.id === formalIdNum);
  if (isError) {
    // TODO: return an error!
    return <Navigate to="/" />;
  }
  if (isLoading && !formal) {
    // TODO: return something better!
    return <Box></Box>;
  }
  if (!formal) {
    // Hmmm...
    return <Box></Box>
  }
  return <FormalInfoView formal={formal} />;
}
