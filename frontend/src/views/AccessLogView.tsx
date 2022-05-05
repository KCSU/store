import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Box,
  Code,
  Container,
  Heading,
  Icon,
  Link,
  Text,
} from "@chakra-ui/react";
import dayjs from "dayjs";
import _ from "lodash";
import { useCallback, useMemo, useState } from "react";
import { FaExternalLinkAlt } from "react-icons/fa";
import { BackButton } from "../components/utility/BackButton";
import { Card } from "../components/utility/Card";
import { useAccessLogs } from "../hooks/queries/useAccessLogs";

export function AccessLogView() {
  const [page, setPage] = useState(1);
  const [size, setSize] = useState(20);
  const { isLoading, data, isError } = useAccessLogs(page, size);
  const { logs, canNextPage } = useMemo(() => {
    return {
      logs: data?.slice(0, size) ?? [],
      canNextPage: (data?.length ?? 0) > size,
    };
  }, [data, size]);
  if (isError) {
    return <></>;
  }
  return (
    <Container maxW="container.md" p={0}>
      <BackButton to="/settings">Back to Settings</BackButton>
      <Card mb={5}>
        <Heading as="h3" size="lg" mb={4}>
          Admin Access Logs
        </Heading>
        <Accordion allowMultiple>
          {logs.map((log) => (
            <AccordionItem id={log.id}>
              <AccordionButton fontFamily="mono" fontSize="sm">
                <Box flex="1" textAlign="left" isTruncated>
                  [
                  <Text as="span" color="green.500">
                    {dayjs(log.createdAt).format("YYYY-MM-DD HH:mm")}
                  </Text>
                  ]{" "}
                  <Text as="span" color="gray">
                    {log.message}
                  </Text>
                </Box>
                <AccordionIcon />
              </AccordionButton>
              <AccordionPanel pb={4}>
                <Text>
                  <Text as="span" fontWeight="bold">Message: </Text>
                  <Text as="span" fontFamily="mono" fontSize="sm">{log.message}</Text>
                </Text>
                <Text>
                  <Text as="span" fontWeight="bold">User: </Text>
                  <Link href={`mailto:${log.email}`} color="teal.500" fontSize="sm">
                    {log.email} <Icon as={FaExternalLinkAlt} boxSize={3} ml={0.5} />
                  </Link>
                </Text>
                {Object.entries(log.metadata).map(([key, value]) => (
                  <Text id={key}>
                    <Text as="span" fontWeight="bold">{_.startCase(key)}: </Text>
                    <Text as="span" fontFamily="mono" fontSize="sm">{value}</Text>
                  </Text>
                ))}
              </AccordionPanel>
            </AccordionItem>
          ))}
        </Accordion>
      </Card>
    </Container>
  );
}
