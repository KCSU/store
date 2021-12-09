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
import { Link } from "react-router-dom";
import { formatMoney } from "../../helpers/formatMoney";
import { Formal } from "../../model/Formal";
import { Card } from "../utility/Card";
import { TicketBuyForm } from "./TicketBuyForm";
import { TicketOptions } from "./TicketOptions";

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

function getBuyText(formal: Formal): string {
  if (formal.saleStart > new Date()) {
    return "Join Queue";
  } else if (
    formal.guestTicketsRemaining === 0 &&
    formal.ticketsRemaining === 0
  ) {
    return "Join Waiting List";
  }
  return "Buy Tickets";
}

const BuyButton = forwardRef<FormalProps, "button">(
  ({ formal, ...props }, ref) => {
    const text = getBuyText(formal);
    const disabled = formal.saleEnd < new Date();
    return (
      <Button
        ref={ref}
        size="sm"
        rightIcon={<Icon as={FaArrowRight} />}
        colorScheme="purple"
        disabled={disabled}
        {...props}
      >
        {text}
      </Button>
    );
  }
);

export const FormalOverview = forwardRef<FormalProps, "div">(
  ({ formal }, ref) => {
    const modalBg = useColorModeValue("gray.50", "gray.800");
    const { isOpen, onOpen, onClose } = useDisclosure();
    return (
      <Card ref={ref}>
        <HStack>
          <Heading size="md">{formal.title}</Heading>
          <FormalStatusTag formal={formal} />
        </HStack>
        <Divider my={2} />
        <FormalStats formal={formal} />
        {/* <Divider my={2} /> */}
        <HStack mt={4}>
          <Button size="sm" as={Link} to={`/formals/${formal.id}`}>
            More Info
          </Button>
          <BuyButton formal={formal} onClick={onOpen}></BuyButton>
        </HStack>
        <Modal isOpen={isOpen} onClose={onClose}>
          <ModalOverlay />
          <ModalContent bg={modalBg}>
            <ModalHeader>Ticket Purchase</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
              <TicketBuyForm formal={formal} />
            </ModalBody>

            <ModalFooter>
              <Button colorScheme="purple" mr={3} onClick={onClose}>
                {getBuyText(formal)}
              </Button>
              <Button variant="ghost" onClick={onClose}>
                Cancel
              </Button>
            </ModalFooter>
          </ModalContent>
        </Modal>
      </Card>
    );
  }
);
