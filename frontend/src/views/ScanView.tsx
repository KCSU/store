import {
  Alert,
  AlertIcon,
  Badge,
  Box,
  Heading,
  Skeleton,
  Spinner,
  Text,
  useColorModeValue,
  VStack,
} from "@chakra-ui/react";
import dayjs from "dayjs";
import { useState } from "react";
import { QrReader } from "react-qr-reader";
import { Card } from "../components/utility/Card";
import { useScanTicket } from "../hooks/queries/useScanTicket";

function validateUUIDv4(uuid: string) {
  const uuidRegex =
    /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
  return uuidRegex.test(uuid);
}

function ScanTicket({ id }: { id: string }) {
  const { isLoading, isError, data } = useScanTicket(id);
  if (isLoading) {
    return <Spinner />;
  }
  if (isError) {
    return (
      <Alert status="error" maxW="500px">
        <AlertIcon />
        Invalid ticket.
      </Alert>
    );
  }
  const colours: Record<string, string | undefined> = {
    "Normal": "",
    "Pescetarian": "blue",
    "Vegetarian": "green",
    "Vegan": "teal",
  }
  let colour = colours[data?.option ?? ""] ?? "pink";
  if (colour !== "") {
    colour = `${colour}.500`;
  }
  return (
    <Card maxW="500px" border="2px" borderColor={colour || "transparent"}>
      {data?.isScanned && (
        <Alert status="warning" mb={2}>
          <AlertIcon />
          Ticket already scanned.
        </Alert>
      )}
      <Heading size="md" display="flex" alignItems="center">
        {data?.formalName}
        {data?.isGuest ? (
          <Badge fontSize="sm" colorScheme="teal" ml={2}>
            Guest
          </Badge>
        ) : (
          <Badge fontSize="sm" colorScheme="brand" ml={2}>
            King's Member
          </Badge>
        )}
      </Heading>
      <Text fontSize="md">{dayjs(data?.formalDate).format("ll, HH:mm")}</Text>
      <Text>
        <Text as="strong">Name: </Text>
        {data?.userName}
      </Text>
      <Text>
        <Text as="strong">Meal option: </Text>
        {data?.option}
      </Text>
    </Card>
  );
}

export function ScanView() {
  const [data, setData] = useState<string | null>(null);
  const isValid = validateUUIDv4(data ?? "");
  return (
    <VStack>
      <Heading size="xl" as="h1" mb={5}>
        KiFoMaSy
      </Heading>
      <Heading size="lg" as="h2" mb={5}>
        Scan Tickets
      </Heading>
      <QrReader
        onResult={(result, error) => {
          if (!!result) {
            setData(result?.toString());
          }

          if (!!error) {
            console.info(error);
          }
        }}
        constraints={{
          aspectRatio: 1,
        }}
        containerStyle={{
          maxWidth: "300px",
          width: "100%"
        }}
        videoStyle={{
          borderRadius: "8px",
        }}
        // style={{ width: '100%' }}
      />
      <Box mt={4} />
      {data === null ? (
        <Alert status="info" maxW="500px">
          <AlertIcon />
          No QR code detected
        </Alert>
      ) : isValid ? (
        <ScanTicket id={data} />
      ) : (
        <Alert status="error" maxW="500px">
          <AlertIcon />
          Invalid ticket.
        </Alert>
      )}
    </VStack>
  );
}
