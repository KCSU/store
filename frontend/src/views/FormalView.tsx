import {
  Badge,
  Box,
  Button,
  Container,
  Heading,
  Text,
  VStack,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";
import { Link, Navigate, useNavigate, useParams } from "react-router-dom";
import { BuyTicketForm } from "../components/tickets/BuyTicketForm";
import { BackButton } from "../components/utility/BackButton";
import { Card } from "../components/utility/Card";
import { useTicketPermissions } from "../hooks/state/useTicketPermissions";
import { formatMoney } from "../helpers/formatMoney";
import { getBuyText } from "../helpers/getBuyText";
import { useBuyTicket } from "../hooks/mutations/useBuyTicket";
import { useDateTime } from "../hooks/state/useDateTime";
import { useFormals } from "../hooks/queries/useFormals";
import { useQueueRequest } from "../hooks/state/useQueueRequest";
import { Formal } from "../model/Formal";
import { useMemo } from "react";
import { FaEdit } from "react-icons/fa";

interface FormalTicketStatsProps {
  prefix?: string;
  price: number;
  tickets: number;
  ticketsRemaining: number;
}

const FormalTicketStats: React.FC<FormalTicketStatsProps> = (props) => {
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

interface FormalCardProps {
  formal: Formal;
}

function FormalCard({ formal }: FormalCardProps) {
  // Formal Data
  const datetime = useDateTime(formal.dateTime);
  const prefix = formal.guestLimit > 0 ? "King's " : "";
  const mutation = useBuyTicket();
  const navigate = useNavigate();
  // State management
  const [queueRequest, dispatchQR] = useQueueRequest(formal.id);
  const { hasTicket, canBuy } = useTicketPermissions(formal);
  const ft = useMemo(() => formal.myTickets?.find(t => !t.isGuest)?.id ?? "", [formal]);
  return (
    // TODO: guest list, responsive meal option
    <Container maxW="container.md" p={0}>
      <BackButton>Back Home</BackButton>
      <Card mb={5}>
        <Heading as="h3" size="lg" mb={1}>
          {formal.name}
        </Heading>
        <Text fontWeight="bold">
          {datetime}
        </Text>
        <Text mt={1} mb={4}>
          Available to {' '}
          {formal.groups?.map(g => g.name).join(", ")}
        </Text>
        <VStack alignItems="stretch">
          <Wrap justifyContent="space-between">
            <FormalTicketStats
              price={formal.price}
              tickets={formal.tickets}
              ticketsRemaining={formal.ticketsRemaining}
              prefix={prefix}
            ></FormalTicketStats>
            {formal.guestLimit > 0 ? (
              <FormalTicketStats
                price={formal.guestPrice}
                tickets={formal.guestTickets}
                ticketsRemaining={formal.guestTicketsRemaining}
                prefix="Guest "
              >
                <br />
                <Text as="i">
                  (up to {formal.guestLimit} per King's member)
                </Text>
              </FormalTicketStats>
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
          {canBuy && (
            <VStack align="stretch" borderWidth="1px" borderRadius="md" p={3}>
              <BuyTicketForm
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
          )}
          {hasTicket && (
            <VStack align="stretch" borderWidth="1px" borderRadius="md" p={3}>
              <Heading as="h5" size="sm">
                My Tickets
              </Heading>
              {formal.myTickets?.map((t) => (
                <Text key={t.id}>
                  {t.isGuest ? "Guest " : "King's "} Ticket ({t.option}) &mdash;{" "}
                  <Text as="b">
                    {formatMoney(t.isGuest ? formal.guestPrice : formal.price)}
                  </Text>
                  {t.isQueue && (
                    <Badge ml={2} colorScheme="brand" verticalAlign="text-top">
                      In Queue
                    </Badge>
                  )}
                </Text>
              ))}
              <Button
                size="sm"
                alignSelf="start"
                // variant="outline"
                // leftIcon={<FaEdit />}}
                as={Link}
                to={`/tickets/${ft}`}
              >
                View Details
              </Button>
            </VStack>
          )}
        </VStack>
      </Card>
    </Container>
  );
}

// TODO: Date and time!
export function FormalView() {
  // Get the formal
  const { formalId } = useParams();
  const { data: formals, isLoading, isError } = useFormals();
  const formal = formals?.find((f) => f.id === formalId);
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
    return <Box></Box>;
  }
  return <FormalCard formal={formal} />;
}
