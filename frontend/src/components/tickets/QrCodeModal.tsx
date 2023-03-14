import {
  Button,
  Heading,
  VStack,
  Icon,
  Image,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  useDisclosure,
  Divider,
  Link,
  useColorModeValue,
  Box,
} from "@chakra-ui/react";
import { FaExternalLinkAlt, FaQrcode } from "react-icons/fa";
import { FormalTicket, Ticket } from "../../model/Ticket";

interface QrCodeModalProps {
  ticket: FormalTicket | Ticket;
}

function getQrCodeUrl(id: string) {
  return `https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=${id}`;
}

const QrCode: React.FC<{ id: string }> = ({ id, children }) => {
  const linkColor = useColorModeValue("teal.600", "teal.300");
  return (
    <>
      <Heading size="sm" as="h3">
        {children}
      </Heading>
      <Image src={getQrCodeUrl(id)} />
      <Box pb={4}>
      <Link
        fontSize="sm"
        color={linkColor}
        href={getQrCodeUrl(id)}
        download={`${id}.png`}
        isExternal
      >
        Download <Icon boxSize={3} as={FaExternalLinkAlt} />
      </Link>
      </Box>
    </>
  );
};

export function QrCodeModal({ ticket }: QrCodeModalProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  return (
    <>
      <Button
        size="sm"
        onClick={onOpen}
        leftIcon={<Icon as={FaQrcode} />}
        // to={`/tickets/${ticket.formal.id}/qr`}
      >
        View QR
      </Button>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Ticket QR Codes</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <VStack>
              {"ticket" in ticket ? <>
              <QrCode id={ticket.ticket.id}>King's Ticket</QrCode>
                {ticket.guestTickets.map((guestTicket, i) => (
                  <QrCode key={guestTicket.id} id={guestTicket.id}>
                    Guest Ticket {i + 1} of {ticket.guestTickets.length}
                  </QrCode>)
                )}
              </> : <QrCode id={ticket.id}>Your Ticket</QrCode>}
            </VStack>
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
}
