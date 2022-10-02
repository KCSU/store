import {
  Badge,
  Box,
  Button,
  Container,
  Flex,
  Heading,
  Icon,
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
import { FaEdit, FaUsers } from "react-icons/fa";
import { FormalGuestList } from "../components/admin/FormalGuestList";

interface FormalTicketStatsProps {
  prefix?: string;
  price: number;
  tickets: number;
  ticketsRemaining: number;
  hint?: string;
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
      {props.hint && (
        <Text fontStyle="italic" fontSize="sm" color="gray.300">
          {props.hint}
        </Text>
      )}
    </WrapItem>
  );
};

interface FormalCardProps {
  formal: Formal;
}

function FormalCard({ formal }: FormalCardProps) {
  // Formal Data
  const datetime = useDateTime(formal.dateTime);
  const firstSaleStart = useDateTime(formal.firstSaleStart);
  const secondSaleStart = useDateTime(formal.secondSaleStart);
  const prefix = formal.guestLimit > 0 ? "King's " : "";
  const mutation = useBuyTicket();
  const navigate = useNavigate();
  // State management
  const [queueRequest, dispatchQR] = useQueueRequest(formal.id);
  const {
    hasTicket,
    canBuy,
    isSaleEnded,
    isFirstSaleStarted,
    isSecondSaleStarted,
  } = useTicketPermissions(formal);
  const ft = useMemo(
    () => formal.myTickets?.find((t) => !t.isGuest)?.id ?? "",
    [formal]
  );
  return (
    // TODO: guest list, responsive meal option
    <Container maxW="container.md" p={0}>
      <BackButton>Back Home</BackButton>
      <Card mb={5}>
        <Flex justifyContent="space-between" alignItems="center">
          <Heading as="h3" size="lg" mb={1}>
            {formal.name}
          </Heading>
          {formal.hasGuestList && <FormalGuestList formal={formal} />}
        </Flex>
        <Text fontWeight="bold">{datetime}</Text>
        {!isFirstSaleStarted && !isSaleEnded && (
          <Text fontSize="sm">Tickets on sale from {firstSaleStart}</Text>
        )}
        {isFirstSaleStarted && !isSecondSaleStarted && !isSaleEnded && (
          <Text fontSize="sm">More tickets {secondSaleStart}</Text>
        )}
        <Text mt={1} mb={4}>
          Available to {formal.groups?.map((g) => g.name).join(", ")}
        </Text>
        <VStack alignItems="stretch">
          <Wrap justifyContent="space-between">
            {isSecondSaleStarted ? (
              <FormalTicketStats
                price={formal.price}
                tickets={formal.firstSaleTickets + formal.secondSaleTickets}
                ticketsRemaining={formal.ticketsRemaining}
                prefix={prefix}
              ></FormalTicketStats>
            ) : (
              <FormalTicketStats
                price={formal.price}
                tickets={formal.firstSaleTickets}
                ticketsRemaining={
                  formal.ticketsRemaining - formal.secondSaleTickets
                }
                prefix={prefix}
                hint={`${formal.secondSaleTickets} unreleased tickets`}
              ></FormalTicketStats>
            )}

            {formal.guestLimit > 0 ? (
              <FormalTicketStats
                price={formal.guestPrice}
                tickets={
                  isSecondSaleStarted
                    ? formal.firstSaleGuestTickets +
                      formal.secondSaleGuestTickets
                    : formal.firstSaleGuestTickets
                }
                ticketsRemaining={
                  isSecondSaleStarted
                    ? formal.guestTicketsRemaining
                    : formal.guestTicketsRemaining -
                      formal.secondSaleGuestTickets
                }
                prefix="Guest "
                hint={
                  isSecondSaleStarted
                    ? undefined
                    : `${formal.secondSaleGuestTickets} unreleased tickets`
                }
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
            <Text whiteSpace="pre-line">
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
                  <Badge
                    ml={2}
                    colorScheme={t.isQueue ? "brand" : "green"}
                    verticalAlign="text-top"
                  >
                    {t.isQueue ? "In Queue" : "Confirmed"}
                  </Badge>
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
