import {
  Badge,
  Button,
  Divider,
  Text,
  forwardRef,
  Heading,
  HStack,
  Stat,
  StatHelpText,
  StatLabel,
  StatNumber,
  Box,
  Flex,
  useColorModeValue,
  Center,
  useDisclosure,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Icon,
} from "@chakra-ui/react";
import { FaArrowRight } from "react-icons/fa";
import { Link, useNavigate } from "react-router-dom";
import { canBuyTicket } from "../../helpers/canBuyTicket";
import { formatMoney } from "../../helpers/formatMoney";
import { getBuyText } from "../../helpers/getBuyText";
import { useBuyTicket } from "../../hooks/useBuyTicket";
import { useDateTime } from "../../hooks/useDateTime";
import { useQueueRequest } from "../../hooks/useQueueRequest";
import { Formal } from "../../model/Formal";
import { Card } from "../utility/Card";
import { TicketBuyForm } from "./TicketBuyForm";

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
        <StatNumber>{formatMoney(props.price)}</StatNumber>
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
        <Center flex="1" bg={noGuestBg} borderRadius={5} p={4}>
          <Text as="i" color={noGuestFg} textAlign="center">
            Guest tickets unavailable
          </Text>
        </Center>
      )}
    </Flex>
  );
}

const BuyButton = forwardRef<FormalProps, "button">(
  ({ formal, ...props }, ref) => {
    const text = getBuyText(formal);
    // What if ticket already bought??
    const isDisabled = !canBuyTicket(formal);
    return (
      <Button
        ref={ref}
        size="sm"
        rightIcon={<Icon as={FaArrowRight} />}
        colorScheme="brand"
        isDisabled={isDisabled}
        {...props}
      >
        {text}
      </Button>
    );
  }
);

interface BuyTicketModalProps {
  isOpen: boolean;
  onClose: () => void;
  formal: Formal;
}

function BuyTicketModal({ isOpen, onClose, formal }: BuyTicketModalProps) {
  const navigate = useNavigate();
  const modalBg = useColorModeValue("gray.50", "gray.800");
  const [queueRequest, dispatchQR] = useQueueRequest(formal.id);
  const mutation = useBuyTicket();
  const isDisabled = !canBuyTicket(formal);
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent bg={modalBg}>
        <ModalHeader>Ticket Purchase</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <TicketBuyForm
            value={queueRequest}
            formal={formal}
            onChange={dispatchQR}
          />
        </ModalBody>

        <ModalFooter>
          <Button
            isLoading={mutation.isLoading}
            isDisabled={isDisabled}
            colorScheme="brand"
            mr={3}
            onClick={async () => {
              await mutation.mutateAsync(queueRequest);
              onClose();
              // TODO: fix this
              setTimeout(() => navigate("/tickets"), 300);
            }}
          >
            {getBuyText(formal)}
          </Button>
          <Button
            variant="ghost"
            onClick={onClose}
            isDisabled={mutation.isLoading}
          >
            Cancel
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}

export const FormalOverview = forwardRef<FormalProps, "div">(
  ({ formal }, ref) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const datetime = useDateTime(formal.dateTime);
    return (
      <Card ref={ref}>
        <HStack mb="2">
          <Heading size="md">{formal.name}</Heading>
          <FormalStatusTag formal={formal} />
        </HStack>
        <Text as="b">{datetime}</Text>
        <Divider my={2} />
        <FormalStats formal={formal} />
        {/* <Divider my={2} /> */}
        <HStack mt={4}>
          <Button size="sm" as={Link} to={`/formals/${formal.id}`}>
            More Info
          </Button>
          <BuyButton formal={formal} onClick={onOpen}></BuyButton>
        </HStack>
        <BuyTicketModal formal={formal} onClose={onClose} isOpen={isOpen} />
      </Card>
    );
  }
);
