import { Alert, AlertIcon, Text, VStack } from "@chakra-ui/react";
import dayjs from "dayjs";

interface QueueWarningProps {
  firstSaleStart: Date;
}

export function QueueWarning({ firstSaleStart }: QueueWarningProps) {
  if (dayjs(firstSaleStart).diff(dayjs(), "hour", true) >= 24) {
    return null;
  }
  return (
    <Alert
      status="warning"
      borderRadius="md"
      gridColumn="1 / -1"
      alignItems="flex-start"
      p={4}
    >
      <AlertIcon />
      <VStack align="start" spacing={1}>
        <Text fontSize="sm" fontWeight="bold" mb={1}>
          WARNING: There is less than 24 hours until your configured sale start
          time.
        </Text>
        <Text fontSize="sm">
          You may want to ensure buyers have enough time to join the pre-sale
          queue before launch.
        </Text>
        <Text fontSize="sm">
          Everyone in the queue is entered into a random draw when the sale
          starts, for a fair outcome.
        </Text>
        <Text fontSize="sm">
          Once the sale starts, tickets are first-come first-served, which is
          potentially undesirable.
        </Text>
      </VStack>
    </Alert>
  );
}
