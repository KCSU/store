import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Box,
  Container,
  Flex,
  Heading,
  Icon,
  IconButton,
  Link,
  Select,
  Text,
  Tooltip,
} from "@chakra-ui/react";
import dayjs from "dayjs";
import _ from "lodash";
import { useCallback, useMemo, useState } from "react";
import {
  FaAngleDoubleLeft,
  FaAngleLeft,
  FaAngleRight,
  FaExternalLinkAlt,
} from "react-icons/fa";
import { BackButton } from "../components/utility/BackButton";
import { Card } from "../components/utility/Card";
import { useAccessLogs } from "../hooks/queries/useAccessLogs";

export function AccessLogView() {
  const [page, setPage] = useState(1);
  const [size, setSize] = useState(20);
  const { isLoading, data, isError, isFetching } = useAccessLogs(page, size);
  const { logs, canNextPage, canPreviousPage } = useMemo(() => {
    return {
      logs: data?.slice(0, size) ?? [],
      canNextPage: (data?.length ?? 0) > size,
      canPreviousPage: page > 1,
    };
  }, [data, size]);
  const nextPage = useCallback(() => {
    setPage(page + 1);
  }, [page]);
  const previousPage = useCallback(() => {
    setPage(page - 1);
  }, [page]);
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
            <AccordionItem key={log.id}>
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
                  <Text as="span" fontWeight="bold">
                    Message:{" "}
                  </Text>
                  <Text as="span" fontFamily="mono" fontSize="sm">
                    {log.message}
                  </Text>
                </Text>
                <Text>
                  <Text as="span" fontWeight="bold">
                    User:{" "}
                  </Text>
                  <Link
                    href={`mailto:${log.email}`}
                    color="teal.500"
                    fontSize="sm"
                  >
                    {log.email}{" "}
                    <Icon as={FaExternalLinkAlt} boxSize={3} ml={0.5} />
                  </Link>
                </Text>
                {Object.entries(log.metadata).map(([key, value]) => (
                  <Text key={key}>
                    <Text as="span" fontWeight="bold">
                      {_.startCase(key)}:{" "}
                    </Text>
                    <Text as="span" fontFamily="mono" fontSize="sm">
                      {value}
                    </Text>
                  </Text>
                ))}
              </AccordionPanel>
            </AccordionItem>
          ))}
        </Accordion>
        <Flex justifyContent="space-between" alignItems="center" mt={2}>
          <Flex>
            <Tooltip label="First Page">
              <IconButton
                mr={1}
                size="sm"
                onClick={() => setPage(1)}
                isDisabled={!canPreviousPage}
                aria-label="first page"
                icon={<Icon as={FaAngleDoubleLeft} h={3} w={3} />}
              />
            </Tooltip>
            <Tooltip label="Previous Page">
              <IconButton
                size="sm"
                aria-label="previous page"
                onClick={previousPage}
                isDisabled={!canPreviousPage}
                icon={<Icon as={FaAngleLeft} h={3} w={3} />}
                mr={4}
              />
            </Tooltip>
          </Flex>
          <Text flexShrink="0" mx={1}>
            Page{" "}
            <Text fontWeight="bold" as="span">
              {page}
            </Text>
          </Text>

          <Select
            maxW={32}
            size="sm"
            value={size}
            onChange={(e) => {
              setSize(Number(e.target.value));
            }}
          >
            {[10, 20, 30, 40, 50].map((size) => (
              <option key={size} value={size}>
                Show {size}
              </option>
            ))}
          </Select>

          <Tooltip label="Next Page">
            <IconButton
              size="sm"
              onClick={nextPage}
              aria-label="next page"
              isDisabled={!canNextPage}
              isLoading={isFetching}
              icon={<Icon as={FaAngleRight} h={3} w={3} />}
              ml={4}
            />
          </Tooltip>
        </Flex>
      </Card>
    </Container>
  );
}
