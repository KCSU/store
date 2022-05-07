import {
  Alert,
  AlertIcon,
  background,
  Button,
  Flex,
  HStack,
  Icon,
  IconButton,
  Input,
  InputGroup,
  InputLeftAddon,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  Select,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tooltip,
  Tr,
  useColorModeValue,
  useDisclosure,
} from "@chakra-ui/react";
import { useEffect, useMemo, useState } from "react";
import {
  FaAngleDoubleLeft,
  FaAngleDoubleRight,
  FaAngleLeft,
  FaAngleRight,
  FaSearch,
  FaUsers,
} from "react-icons/fa";
import { Column, useGlobalFilter, usePagination, useTable } from "react-table";
import { useFormalGuestList } from "../../hooks/queries/useFormalGuestList";
import { Formal } from "../../model/Formal";
import { FormalGuest } from "../../model/FormalGuest";

interface FormalProps {
  formal: Formal;
}

function FormalGuestsTable({ formal }: FormalProps) {
  const { data, isLoading, isError } = useFormalGuestList(formal.id);
  const guests = useMemo(() => data ?? [], [data]);
  const [query, setQuery] = useState("");
  const columns = useMemo<Column<FormalGuest>[]>(() => {
    return [
      {
        Header: "Name",
        accessor: "name",
      },
      // {
      //   Header: "CRSID",
      //   accessor: "crsid",
      // }
      {
        Header: "Guests",
        accessor: "guests",
        isNumeric: true,
      },
    ];
  }, []);
  const {
    getTableProps,
    getTableBodyProps,
    prepareRow,
    headerGroups,
    gotoPage,
    canPreviousPage,
    canNextPage,
    page,
    previousPage,
    nextPage,
    setGlobalFilter,
    pageOptions,
    state: { pageIndex, pageSize },
    pageCount,
  } = useTable({ columns, data: guests }, useGlobalFilter, usePagination);
  useEffect(() => setGlobalFilter(query), [query]);
  const background = useColorModeValue("white", "gray.750");
  if (isError) {
    return (
      <Alert status="error" mb={3}>
        <AlertIcon />
        An error occurred while loading the guest list.
      </Alert>
    );
  }
  if (isLoading && !data) {
    return <></>;
  }
  return (
    <>
      <InputGroup size="sm" mb={4}>
        <InputLeftAddon>
          <Icon as={FaSearch} />
        </InputLeftAddon>
        <Input
          id="query"
          autoComplete="off"
          // maxW="300px"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search guests..."
        />
      </InputGroup>
      <Table my={2} variant="striped" size="sm" {...getTableProps()}>
        <Thead>
          {headerGroups.map((headerGroup) => (
            <Tr {...headerGroup.getHeaderGroupProps()}>
              {headerGroup.headers.map((column) => (
                <Th
                  {...column.getHeaderProps()}
                  isNumeric={(column as any).isNumeric}
                >
                  {column.render("Header")}
                </Th>
              ))}
            </Tr>
          ))}
        </Thead>
        <Tbody {...getTableBodyProps()} bg={background}>
          {page.map((row) => {
            prepareRow(row);
            return (
              <Tr {...row.getRowProps()}>
                {row.cells.map((cell) => (
                  <Td
                    {...cell.getCellProps()}
                    isNumeric={(cell.column as any).isNumeric}
                  >
                    {cell.render("Cell")}
                  </Td>
                ))}
              </Tr>
            );
          })}
        </Tbody>
      </Table>
      <Flex justifyContent="space-between" alignItems="center" mt={2} mb={3}>
        <Flex>
          <Tooltip label="First Page">
            <IconButton
              mr={1}
              size="sm"
              onClick={() => gotoPage(0)}
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
            {pageIndex + 1}
          </Text>{" "}
          of{" "}
          <Text fontWeight="bold" as="span">
            {pageOptions.length}
          </Text>
        </Text>
        <Flex>
          <Tooltip label="Next Page">
            <IconButton
              size="sm"
              onClick={nextPage}
              aria-label="next page"
              isDisabled={!canNextPage}
              icon={<Icon as={FaAngleRight} h={3} w={3} />}
              ml={4}
            />
          </Tooltip>
          <Tooltip label="Last Page">
            <IconButton
              size="sm"
              aria-label="last page"
              onClick={() => gotoPage(pageCount - 1)}
              isDisabled={!canNextPage}
              icon={<Icon as={FaAngleDoubleRight} h={3} w={3} />}
              ml={1}
            />
          </Tooltip>
        </Flex>
      </Flex>
    </>
  );
}

export function FormalGuestList({ formal }: FormalProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  return (
    <>
      <Button leftIcon={<Icon as={FaUsers} />} size="sm" onClick={onOpen}>
        View Guest List
      </Button>
      <Modal isOpen={isOpen} size="lg" onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>{formal.name}: Guest List</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <FormalGuestsTable formal={formal} />
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
}
